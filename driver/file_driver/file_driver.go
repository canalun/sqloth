package file_driver

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/canalun/sqloth/domain/model"
)

type FileDriver struct {
	FilePath string
}

func NewFileDriver(filePath string) FileDriver {
	return FileDriver{
		FilePath: filePath,
	}
}

// GetSchema() parses the given schema file and return Schema model.
// The returned model has the relation of tables and columns, and foreign key constraints.
func (fd FileDriver) GetSchema() model.Schema {
	f, err := os.Open(fd.FilePath)
	if err != nil {
		fmt.Println("error, cannot open the file")
		return model.Schema{}
	}
	defer f.Close()

	schema := model.Schema{}

	// prepare map slice to avoid double loop in parsing constraint
	columnIndexWithinTable := []map[model.ColumnName]int{}
	columnIndex := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lt := checkLineType(line)

		switch lt {
		case "tableLine":
			table := parseTableLine(line)
			schema.AddTable(table)

			columnIndexWithinTable = append(columnIndexWithinTable, map[model.ColumnName]int{})
			columnIndex = 0

		case "columnLine":
			tableName := schema.LastTable().Name
			column, err := parseColumnLine(line, tableName)
			if err != nil {
				fmt.Print(err)
				return model.Schema{}
			}
			schema.LastTable().AddColumns(column)

			columnIndexWithinTable[len(columnIndexWithinTable)-1][column.Name] = columnIndex
			columnIndex += 1

		case "constraintLine":
			parsedConstraints := parseConstraintLine(line)
			for _, c := range parsedConstraints {
				if index, ok := columnIndexWithinTable[len(columnIndexWithinTable)-1][c.BoundedColumnName]; ok {
					schema.LastTable().Columns[index].Constraints = append(schema.LastTable().Columns[index].Constraints, c.Constraint)
				} else {
					fmt.Println("error when parsing constraint")
					return model.Schema{}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error, something happens when/while reading the schema")
		return model.Schema{}
	}

	return schema
}

// the below regular expressions assumes;
// 		tableLines is an array as ["CREATE TABLE `name` ..."]
// 		columnLines is an array as ["  `name` type ..."]
// 		constraintLines is an array like ["[CONSTRAINT [symbol]] FOREIGN KEY (`column_name`) REFERENCES `table_name` (`column_name`) ..."]
var regexForTableLine = regexp.MustCompile("(?m)^ *CREATE TABLE .*")
var regexForColumnLine = regexp.MustCompile("(?m)^ *`.*` *[^ ]+.*,")
var regexForConstraintLine = regexp.MustCompile("(?m)^.* FOREIGN KEY .*")

func checkLineType(line string) string {
	if tableLines := regexForTableLine.FindStringSubmatch(line); len(tableLines) > 0 {
		return "tableLine"
	} else if columnLines := regexForColumnLine.FindStringSubmatch(line); len(columnLines) > 0 {
		return "columnLine"
	} else if constraintLines := regexForConstraintLine.FindStringSubmatch(line); len(constraintLines) > 0 {
		return "constraintLine"
	} else {
		return ""
	}
}

////// parsers ///////////////////////////

func parseTableLine(tableLine string) model.Table {
	_tableName := trimSymbols(strings.Fields(tableLine)[2])
	tableName := model.TableName(_tableName)

	table := model.NewTable(tableName)
	return table
}

var regexForAutoIncrement = regexp.MustCompile("AUTO_INCREMENT")
var regexForUnsigned = regexp.MustCompile("(UNSIGNED)|(Unsigned)|(unsigned)")

// parseColumnLine() parses not only column name, but also data type and type attribute.
func parseColumnLine(columnLine string, tableName model.TableName) (model.Column, error) {
	columnName := trimSymbols(strings.Fields(columnLine)[0])
	_columnFullName := string(tableName) + "." + columnName
	columnFullName := model.ColumnFullName(_columnFullName)

	columnType, err := model.NewColumnType(strings.Fields(columnLine)[1])
	if err != nil {
		fmt.Println("unexpected data type")
		return model.Column{}, err
	}

	column := model.NewColumn(columnFullName, columnType)
	if len(regexForAutoIncrement.FindStringSubmatch(columnLine)) > 0 {
		column.SetAutoIncrement()
	}
	if len(regexForUnsigned.FindStringSubmatch(columnLine)) > 0 {
		column.SetUnsigned(true)
	} else {
		column.SetUnsigned(false)
	}
	return column, nil
}

// need to return slices in case constraints are written in multi-index style
type parsedConstraint struct {
	Constraint        model.Constraint
	BoundedColumnName model.ColumnName
}

func parseConstraintLine(constraintLine string) []parsedConstraint {
	// foreign key format is as below for mysql.
	// [CONSTRAINT [symbol]] FOREIGN KEY
	// 		[index_name] (index_col_name, ...)
	// 		REFERENCES tbl_name (index_col_name,...)
	// 		[ON DELETE reference_option]
	// 		[ON UPDATE reference_option]
	// https://dev.mysql.com/doc/refman/5.6/ja/create-table-foreign-keys.html

	words := strings.Fields(constraintLine)

	// cut optional words before (index_col_name, ...)
	for i, w := range words {
		fmt.Printf("%#v\n", w)
		if w[0:1] == "(" {
			words = words[i:]
			break
		}
	}

	boundedColumnNames := strings.Split(trimSymbols(words[0]), ",")
	referencedTableName := trimSymbols(words[2])
	referencedColumnNames := strings.Split(trimSymbols(words[3]), ",")

	re := []parsedConstraint{}
	for i, boundedColumnName := range boundedColumnNames {
		constraint := model.NewConstraint(model.TableName(referencedTableName), model.ColumnName(referencedColumnNames[i]))
		re = append(re, parsedConstraint{
			Constraint:        constraint,
			BoundedColumnName: model.ColumnName(boundedColumnName),
		})
	}

	return re
}

//////////////////////////////////////////

func trimSymbols(str string) string {
	trimmingTarget := "()`"
	return strings.Trim(str, trimmingTarget)
}

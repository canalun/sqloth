package file_driver

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/canalun/sqloth/domain/model"
)

var regexForTable = regexp.MustCompile("(?m)^ *CREATE TABLE .*")
var regexForColumn = regexp.MustCompile("(?m)^ *`.*` *[^ ]+.*,")
var regexForAutoIncrement = regexp.MustCompile("AUTO_INCREMENT")
var regexForUnsigned = regexp.MustCompile("UNSIGNED")
var regexForColumnConstraint = regexp.MustCompile("(?m)^ *CONSTRAINT .*")

type FileDriver struct {
	FilePath string
}

func NewFileDriver(filePath string) FileDriver {
	return FileDriver{
		FilePath: filePath,
	}
}

// TODO: fix to handle multi-index-key
func (fd FileDriver) GetSchema() model.Schema {
	f, err := os.Open(fd.FilePath)
	if err != nil {
		fmt.Println("error, cannot open the file")
		return model.Schema{}
	}
	defer f.Close()

	schema := model.Schema{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		tableLines := regexForTable.FindStringSubmatch(line)
		if len(tableLines) > 0 {
			// assuming tableLines are like ["CREATE TABLE `name` ..."]
			_tableName := trimSqlQuery(strings.Fields(tableLines[0])[2])
			tableName := model.TableName(_tableName)

			table := model.NewTable(tableName, []model.Column{})
			schema.AddTable(table)
		}

		columnLines := regexForColumn.FindStringSubmatch(line)
		if len(columnLines) > 0 {
			// assuming columnLines are like ["  `name` type ..."]
			_columnName := trimSqlQuery(strings.Fields(columnLines[0])[0])
			columnName := model.ColumnName(_columnName)

			_columnFullName := string(schema.LastTable().Name) + "." + string(columnName)
			columnFullName := model.ColumnFullName(_columnFullName)

			columnType, err := strToColumnType(strings.Fields(columnLines[0])[1])
			if err != nil {
				fmt.Println("unexpected data type")
				return model.Schema{}
			}

			column := model.NewColumn(columnFullName, columnType)
			if len(regexForAutoIncrement.FindStringSubmatch(columnLines[0])) > 0 {
				column.SetAutoIncrement()
			}
			if len(regexForUnsigned.FindStringSubmatch(columnLines[0])) > 0 {
				column.SetUnsigned()
			}
			schema.LastTable().AddColumns(column)
		}

		columnKeyLines := regexForColumnConstraint.FindStringSubmatch(line)
		if len(columnKeyLines) > 0 {
			// assuming columnLines are like [" CONSTRAINT `XX` FOREIGN KEY (`column_name`) REFERENCES `table_name` (`column_name`) ..."]
			_boundedColumnName := trimSqlQuery(strings.Fields(columnKeyLines[0])[4])
			boundedColumnName := model.ColumnName(_boundedColumnName)

			_referencedTableName := trimSqlQuery(strings.Fields(columnKeyLines[0])[6])
			referencedTableName := model.TableName(_referencedTableName)

			_referencedColumnName := trimSqlQuery(strings.Fields(columnKeyLines[0])[7])
			referencedColumnName := model.ColumnName(_referencedColumnName)

			constraint := model.NewConstraint(referencedTableName, referencedColumnName)
			for i, c := range schema.LastTable().Columns {
				if c.Name == boundedColumnName {
					schema.LastTable().Columns[i].Constraints = append(schema.LastTable().Columns[i].Constraints, constraint)
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

func strToColumnType(str string) (ct model.ColumnType, err error) {
	f := func(r rune) bool {
		return r == '(' || r == ')' || r == ','
	}
	l := strings.FieldsFunc(str, f)

	base, err := model.StrToColumnTypeBase(strings.ToLower(l[0]))
	if err != nil {
		return model.ColumnType{}, err
	}

	var param int
	if len(l) > 1 {
		param, err = strconv.Atoi(l[1])
		if err != nil {
			return model.ColumnType{}, err
		}
	}
	if base == model.Text {
		param = 100
	}
	// TODO: handle varbinary appropriately
	if base == model.Varbinary {
		param = 100
	}
	// TODO: handle mediumblob appropriately
	if base == model.Mediumblob {
		param = 100
	}

	return model.ColumnType{
		Base:  base,
		Param: model.ColumnTypeParam(param),
	}, nil
}

func trimSqlQuery(str string) string {
	trimmingTarget := "()`"
	return strings.Trim(str, trimmingTarget)
}

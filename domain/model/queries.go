package model

func GenerateQuery(rft map[TableName][]Record) []string {
	if len(rft) == 0 {
		return []string{}
	}
	re := []string{"SET FOREIGN KEY = 0;"}
	for tableName, records := range rft {
		q := "INSERT INTO " + string(tableName) + " VALUES "
		for _, record := range records {
			q += querizeRecord(record)
		}
		q = q[:len(q)-1] + ";"
		re = append(re, q)
	}
	re = append(re, "SET FOREIGN KEY = 1;")
	return re
}

func querizeRecord(record Record) string {
	re := "("
	for _, v := range record {
		re += "'" + string(v) + "',"
	}
	re = re[:len(re)-1] + "),"
	return re
}

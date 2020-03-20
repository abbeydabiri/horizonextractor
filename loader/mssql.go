package loader

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" //mssql database for import
	"github.com/jmoiron/sqlx"
)

//MSSQL database
var MSSQL *sqlx.DB

//Connected is used to confirm data connectivity
var Connected bool

//Connect to the MSSQL database
func Connect(host, port, user, pass, db string) {
	if host == "" || port == "" || user == "" || pass == "" || db == "" {
		sMsg := "mssql connection not set"
		fmt.Println(sMsg)
		log.Printf(sMsg)
		return
	}

	//SQL Connection for MSSQL
	// conn := "mssql://%s:%s@%s:%s/instance?database=%s"
	conn := "server=%s;user id=%s;password=%s;database=%s"
	conn = fmt.Sprintf(conn, host, user, pass, db)

	var err error
	MSSQL, err = sqlx.Open("mssql", conn)
	if err != nil {
		sMsg := "Error Connecting Database: " + err.Error()
		fmt.Println(sMsg)
		log.Printf(sMsg)
		return
	}

	err = MSSQL.Ping()
	sMsg := "Connecting To  Database.."
	fmt.Println(sMsg)
	log.Printf(sMsg)

	if err != nil {
		sMsg := "Error Connecting Database: " + err.Error()
		fmt.Println(sMsg)
		log.Printf(sMsg)
		return
	}

	Connected = true
	//SQL Connection for MSSQL
}

//Create  ...
func Create(table Tables) string {
	reflectType := reflect.TypeOf(table).Elem()
	tablename := strings.ToLower(reflectType.Name())
	sqlDrop := "drop table " + tablename

	sqlIndex := ""
	sqlCreate := "create table " + tablename + " ("
	sqlCreate, sqlIndex = createFields(reflectType, tablename, sqlCreate, sqlIndex)
	sqlCreate = strings.TrimSuffix(sqlCreate, ", ") + "); "

	MSSQL.Exec(sqlDrop)
	_, err := MSSQL.Exec(sqlCreate)
	if err != nil {
		log.Println("tablename: " + tablename + " - " + err.Error())
	} else {
		if sqlIndex != "" {
			MSSQL.Exec(sqlIndex)
		}
	}
	return fmt.Sprintf("Table %s created", reflectType.Name())
}

//CreateFields
func createFields(reflectType reflect.Type, tablename, sqlCreate, sqlIndex string) (string, string) {
	indexFmt := "\ncreate %s " + tablename + "_%s on " + tablename + " (%s);"
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tag := field.Tag.Get("sql")
		fieldName := strings.ToLower(field.Name)
		fieldType := sqlTypes[strings.ToLower(field.Type.Name())]

		if fieldType == "" {
			if field.Name == "Fields" {
				sqlCreate, sqlIndex = createFields(field.Type, tablename, sqlCreate, sqlIndex)
			}
			continue
		}

		if fieldName != "id" {
			defaultValue := ""
			switch fieldType {
			case "bool":
				defaultValue = "DEFAULT false"
			case "timestamp":
				defaultValue = "DEFAULT current_timestamp"
			case "text":
				defaultValue = "DEFAULT ''"
			case "float", "float64":
				defaultValue = "DEFAULT 0.0"
			case "int", "int8":
				defaultValue = "DEFAULT 0"
			}
			sqlCreate += fmt.Sprintf("%s %s %s", fieldName, fieldType, defaultValue)
		}

		switch tag {
		case "pk":
			if fieldName == "id" {
				sqlCreate += "id SERIAL PRIMARY KEY"
			}
		case "index", "unique index":
			sqlIndex += fmt.Sprintf(indexFmt, tag, fieldName, fieldName)
		}
		sqlCreate += ", "
	}
	// return sqlCreate, sqlIndex
	return sqlCreate, ""
}

//Insert  ...
func Insert(table Tables, tableMap map[string]interface{}) (string, []interface{}) {
	delete(tableMap, "ID")

	reflectType := reflect.TypeOf(table).Elem()
	tablename := strings.ToLower(reflectType.Name())

	var sqlParams []interface{}
	sqlFields, sqlValues := " (", " ("
	sqlInsert := "insert into " + tablename
	sqlFields, sqlValues, sqlParams = insertFields(sqlFields, sqlValues, reflectType, tableMap, sqlParams)

	sqlInsert += strings.TrimSuffix(sqlFields, ", ") + " ) "
	sqlInsert += " OUTPUT INSERTED.col1 VALUES "
	sqlInsert += strings.TrimSuffix(sqlValues, ", ") + " ) "
	return sqlInsert, sqlParams
}

//InsertFields
func insertFields(sqlFields, sqlValues string, reflectType reflect.Type,
	tableMap map[string]interface{}, sqlParams []interface{}) (string, string, []interface{}) {
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		if tableMap[field.Name] != nil || field.Name == "Fields" {
			fieldName := strings.ToLower(field.Name)
			fieldType := sqlTypes[strings.ToLower(field.Type.Name())]
			if fieldType == "" {
				if fieldName == "fields" {
					sqlFields, sqlValues, sqlParams = insertFields(sqlFields, sqlValues, field.Type, tableMap, sqlParams)
				}
				continue
			}

			switch strings.ToLower(field.Type.Name()) {
			case "int", "int64", "uint", "uint64":
				tableMapFieldType := reflect.TypeOf(tableMap[field.Name])
				switch tableMapFieldType.Kind() {
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
					sqlParams = append(sqlParams, fmt.Sprintf("%d", tableMap[field.Name]))
				default:
					sqlParams = append(sqlParams, fmt.Sprintf("%.f", tableMap[field.Name]))
				}
			default:
				sqlParams = append(sqlParams, tableMap[field.Name])
			}

			sqlFields += fieldName + ", "
			sqlValues += fmt.Sprintf("$%v, ", len(sqlParams))
		}
	}
	return sqlFields, sqlValues, sqlParams
}

//BulkInsert ...
func BulkInsert(table Tables, tableMapSlice []map[string]interface{}) string {
	reflectType := reflect.TypeOf(table).Elem()
	tablename := strings.ToLower(reflectType.Name())
	var sqlFields, sqlValues, sqlInsertBulk string

	for index, tableMap := range tableMapSlice {

		if tableMap["ID"] != nil {
			delete(tableMap, "ID")
		}

		sqlFields, sqlValues = bulkInsertFields(" (", " (", reflectType, tableMap)
		if sqlFields != "" && sqlValues != "" {
			if index%160 == 0 && index < len(tableMap)-1 {
				if sqlInsertBulk != "" {
					sqlInsertBulk = strings.TrimSuffix(sqlInsertBulk, "), ") + ");"
				}
				sqlInsertBulk += fmt.Sprintf("insert into %s %s) VALUES ",
					tablename, strings.TrimSuffix(sqlFields, ", "))
			}
			sqlInsertBulk += strings.TrimSuffix(sqlValues, ", ") + "), "
		}

	}
	return strings.TrimSuffix(sqlInsertBulk, "), ") + "); "
}

//bulkInsertFields
func bulkInsertFields(sqlFields, sqlValues string, reflectType reflect.Type,
	tableMap map[string]interface{}) (string, string) {
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		if tableMap[field.Name] != nil {
			fieldName := strings.ToLower(field.Name)
			fieldType := sqlTypes[strings.ToLower(field.Type.Name())]
			if fieldType == "" || fieldName == "id" {
				continue
			}

			switch strings.ToLower(field.Type.Name()) {
			case "int", "int64", "uint", "uint64":
				tableMapFieldType := reflect.TypeOf(tableMap[field.Name])
				switch tableMapFieldType.Kind() {
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
					sqlValues += fmt.Sprintf("%d, ", tableMap[field.Name])
				default:
					sqlValues += fmt.Sprintf("%.f, ", tableMap[field.Name])
				}
			case "string", "time":
				sqlValues += fmt.Sprintf("'%v', ", strings.Replace(strings.Replace(tableMap[field.Name].(string), "'", "", -1), `"`, ``, -1))
			default: //bool and float
				sqlValues += fmt.Sprintf("%v, ", tableMap[field.Name])
			}
			sqlFields += fieldName + ", "
		}
	}
	return sqlFields, sqlValues
}

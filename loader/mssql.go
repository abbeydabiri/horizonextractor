package loader

import (
	"encoding/json"
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

var dbname string

//Connect to the MSSQL database
func Connect(server, db, user, pass string) {
	if server == "" || db == "" {
		sMsg := "sql server and db not set"
		fmt.Println(sMsg)
		log.Printf(sMsg)
		return
	}

	dbname = db
	conn := "server=%s;database=%s"

	//SQL Connection for MSSQL
	if user == "" || pass == "" {
		conn += ";trusted_connection=yes"
		conn = fmt.Sprintf(conn, server, db)
	} else {
		conn += ";user id=%s;password=%s"
		conn = fmt.Sprintf(conn, server, db, user, pass)
	}

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
		sMsg := conn + " -|- Error Connecting Database: " + err.Error()
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
	// sqlDrop := "drop table " + tablename
	// MSSQL.Exec(sqlDrop)

	//check if table exists before attempting to create
	sqlCheck := "select count(table_name) from information_schema.tables where table_type = 'BASE TABLE' and table_catalog = '%v' and table_name = '%v'"
	sqlCheck = fmt.Sprintf(sqlCheck, dbname, tablename)

	totalExists := 0
	if MSSQL.Get(&totalExists, sqlCheck); totalExists == 0 {
		sqlIndex := ""
		sqlCreate := "create table " + tablename + " ("
		sqlCreate, sqlIndex = createFields(reflectType, tablename, sqlCreate, sqlIndex)
		sqlCreate = strings.TrimSuffix(sqlCreate, ", ") + "); "

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
	return ""
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
			case "datetime":
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
				sqlCreate += "id int IDENTITY(1,1) PRIMARY KEY"
			}
		case "index", "unique index":
			sqlIndex += fmt.Sprintf(indexFmt, tag, fieldName, fieldName)
		}
		sqlCreate += ", "
	}
	// return sqlCreate, sqlIndex
	return sqlCreate, ""
}

//ToMap ...
func ToMap(table Tables) (mapInterface map[string]interface{}) {
	jsonTable, _ := json.Marshal(table)
	json.Unmarshal(jsonTable, &mapInterface)
	return
}

//Insert  ...
func Insert(table Tables, tableMap map[string]interface{}) (string, []interface{}) {
	delete(tableMap, "ID")
	delete(tableMap, "Createdate")
	delete(tableMap, "Updatedate")

	reflectType := reflect.TypeOf(table).Elem()
	tablename := strings.ToLower(reflectType.Name())

	var sqlParams []interface{}
	sqlFields, sqlValues := " (", " ("
	sqlInsert := "insert into " + tablename
	sqlFields, sqlValues, sqlParams = insertFields(sqlFields, sqlValues, reflectType, tableMap, sqlParams)

	sqlInsert += strings.TrimSuffix(sqlFields, ", ") + " ) "
	sqlInsert += " VALUES "
	sqlInsert += strings.TrimSuffix(sqlValues, ", ") + " ); SELECT SCOPE_IDENTITY()"

	return sqlInsert, sqlParams
}

//InsertFields ...
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

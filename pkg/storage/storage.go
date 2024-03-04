package storage

import (
	"database/sql"
	"fmt"
	"os"
	refl "reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

/*
Короче, краткое описание того, что будет происходить в этом модуле.
Пусть у нас есть какие-то внешние типы-структуры. Тогда их названия - названия таблиц, названия их полей - колонки таблиц.
На запуске мы даем модулю потрогать все эти структуры и создаем таблички под них, если до этого таковых не было.
А потом, при необходимости обратиться к таблице, все будет происходить по вышеописанным правилам.
Главное это все научиться делать через reflect, и вот это уже нетривиальная задача.
*/

var db *sql.DB

// Auxilary error type.
type storageError struct{
	what string
}

func (err storageError) Error() string{
	return fmt.Sprintf("Storage error: %s", err.what)
}

// Get the table name from struct type name.
func getName(s any) string{
	return refl.TypeOf(s).Name()
}

// Get the bindings from struct fields.
func getBindings(s any) ([]string, []any, error){
	val := refl.ValueOf(s)
	
	columns := []string{}
	values := []any{}

	if refl.TypeOf(s).Kind() != refl.Struct{
		return columns, values, storageError{"the kind of input interface but be a struct and represent the specific database table columns."}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldKind := val.Field(i).Type().Kind()
		fieldValue := val.Field(i)

		if fieldKind == refl.Int64 || fieldKind == refl.String{
			columns = append(columns, field.Name)
			values = append(values, fieldValue.Interface())
		}else{
			return columns, values, storageError{"the interface struct fields must be either of type int64 or string."}
		}
	}

	return columns, values, nil
}

// Function that creates the table for specific struct type. 
// Should be only ran initially, once for each type that will interact with this package functions later.
func TouchTable(s any) error{
	table := getName(s) 
	
	columns, values, err := getBindings(s)
	if err != nil{
		return err
	}

	fields := []string{}
	for i := 0; i < len(columns); i++{
		kind := ""
		if refl.ValueOf(values[i]).Type().Kind() == refl.Int64{
			kind = "INTEGER"
		}else{
			kind = "TEXT"
		}

		if columns[i] == "id"{
			fields = append(fields, fmt.Sprintf("%v %v PRIMARY KEY NOT NULL", columns[i], kind))
		}else{
			fields = append(fields, fmt.Sprintf("%v %v", columns[i], kind))
		}
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", table, strings.Join(fields, ", "))
	_, err = db.Exec(query, values...)
	return err
}

// Function representing the sqlite INSERT query.
func Insert(s any) error{
	table := getName(s) 
	
	columns, values, err := getBindings(s)
	if err != nil{
		return err
	}
	questionMarks := strings.Join(strings.Split(strings.Repeat("?", len(columns)), ""), ", ")

	query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Join(columns, ", "), questionMarks)

	_, err = db.Query(query, values...)
	return err
}

// Function that executes initial connection to the database.
func Init() error{
	path := "./database/.db"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		
		err = file.Close()
		if err != nil {
			return err
		}
	}

	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	return nil
}

func Close() error{
	err := db.Close()
	if err != nil{
		return err
	}

	return nil
}
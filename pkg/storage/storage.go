package storage

import (
	"database/sql"
	"fmt"
	"os"
	refl "reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// This library is sqlite ORM


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

// Get list of pointers to struct fields. Recieves a pointer to a struct!
func getPointers(s any) ([]any, error) {
	ptr := refl.ValueOf(s)
	if ptr.Kind() != refl.Pointer{
		return nil, storageError{"the kind of interface must be a pointer to a struct."}
	}
	if ptr.Elem().Kind() != refl.Struct{
		return nil, storageError{"the kind of interface must be a pointer to a struct."}
	}
	
	n := ptr.Elem().NumField()
	pointers := make([]any, n)
    for i := 0; i < n; i++ {
        field := ptr.Elem().Field(i).Addr().Interface()
        pointers[i] = field
    }

	return pointers, nil
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

		if columns[i] == "Id"{
			fields = append(fields, fmt.Sprintf("%v %v PRIMARY KEY NOT NULL", columns[i], kind))
		}else{
			fields = append(fields, fmt.Sprintf("%v %v", columns[i], kind))
		}
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v)", table, strings.Join(fields, ", "))
	_, err = db.Exec(query, values...)
	return err
}

// Parametrized function, that performs the sqlite INSERT query.
func Insert[T any](s T) error{
	table := getName(s) 
	
	columns, values, err := getBindings(s)
	if err != nil{
		return err
	}
	questionMarks := strings.Join(strings.Split(strings.Repeat("?", len(columns)), ""), ", ")

	query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Join(columns, ", "), questionMarks)

	_, err = db.Exec(query, values...)
	return err
}

// Parametrized function, that performs the sqlite SELECT query and returns struct of type T.
func Select[T any](id int64) (T, error){
	var s T

	table := getName(s) 
	pointers, err := getPointers(&s)
	if err != nil{
		return s, err
	}

	query := fmt.Sprintf("SELECT * FROM %v WHERE Id = ?", table)
	rows, err := db.Query(query, id)
	if err != nil{
		return s, err
	}

	if !rows.Next(){
		return s, sql.ErrNoRows
	}

	err = rows.Scan(pointers...)
	return s, err
}

// Parametrized function, that performs the sqlite UPDATE query.
func Update[T any](id int64, s T) error{
	table := getName(s) 
	
	columns, values, err := getBindings(s)
	if err != nil{
		return err
	}
	
	update := ""
	for i, column := range columns{
		if i == 0{
			update += fmt.Sprintf("%v = ?", column)
		}else{
			update += fmt.Sprintf(", %v = ?", column)
		}
	}
	query := fmt.Sprintf("UPDATE %v SET %v WHERE Id=?", table, update)

	values = append(values, id)
	_, err = db.Exec(query, values...)
	return err
}

func Delete[T any](id int64) error{
	var s T
	table := getName(s) 
	
	query := fmt.Sprintf("DELETE FROM %v WHERE Id = ?", table)

	_, err := db.Exec(query, id)
	return err
}

// Parametrized function, that performs the sqlite SELECT query for the whole table.
func SelectAll[T any]() ([]T, error){
	var s T
	var ss []T

	table := getName(s) 

	query := fmt.Sprintf("SELECT * FROM %v", table)
	rows, err := db.Query(query)
	if err != nil{
		return ss, err
	}

	for rows.Next(){
		var s T
		pointers, err := getPointers(&s)
		if err != nil{
			return ss, err
		}

		err = rows.Scan(pointers...)
		if err != nil{
			return ss, err
		}

		ss = append(ss, s)
	}

	return ss, err
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

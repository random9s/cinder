package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/go-sql-driver/mysql"
)

//MySQL is a wrapper around a sql DB struct
type MySQL struct {
	*sql.DB
}

//NewMySQL wraps db connection
func NewMySQL(db *sql.DB) *MySQL {
	return &MySQL{db}
}

//MarshalRow updates or deletes a row in a mysql database
func (db *MySQL) MarshalRow(sql string, args ...interface{}) (interface{}, error) {
	return db.Exec(sql, args...)
}

//MarshalRows updates or deletes many rows in a mysql database
func (db *MySQL) MarshalRows(sql string, args ...interface{}) (interface{}, error) {
	return db.Exec(sql, args...)
}

//MarshalField updates or deletes a rows field in a mysql database
func (db *MySQL) MarshalField(sql string, args ...interface{}) (interface{}, error) {
	return db.Exec(sql, args...)
}

//MarshalFields updates or deletes many fields in a mysql database
func (db *MySQL) MarshalFields(sql string, args ...interface{}) (interface{}, error) {
	return db.Exec(sql, args...)
}

//UnmarshalRow retrieves a row from a mysql database
func (db *MySQL) UnmarshalRow(v interface{}, sql string, args ...interface{}) error {
	//v must be a pointer!
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return fmt.Errorf("could not translate object %v", reflect.TypeOf(v))
	}

	//Reports whether rv represents a value
	if ptr.IsValid() {
		rows, err := db.Query(sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			ms, err := structToMS(cols, v)
			if err != nil {
				return err
			}

			//Dump the rows we want to scan into the scanner
			err = rows.Scan(nullValues(ms)...)
			if err != nil {
				return err
			}

			err = assignValsToStruct(ms)
			if err != nil {
				return err
			}

			break
		}

		return rows.Err()
	}

	return fmt.Errorf("invalid type %v", reflect.TypeOf(v))
}

//UnmarshalRows retrieves many rows from a mysql database
func (db *MySQL) UnmarshalRows(v interface{}, sql string, args ...interface{}) error {
	//v must be a pointer!
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return fmt.Errorf("could not translate non pointer %v", reflect.TypeOf(v))
	}

	//v must be a pointer to a slice
	sliceType := reflect.TypeOf(v).Elem()
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("could not translate non slice %v", reflect.TypeOf(v))
	}

	//Return the underlying type of slice
	sliceType = sliceType.Elem()
	//If the item is a pointer to something, get the type it points to
	var appendPointerType = false
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
		appendPointerType = true
	}

	//Reports whether rv represents a value
	if ptr.IsValid() {
		rows, err := db.Query(sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		slice := reflect.ValueOf(v).Elem()
		if slice.Kind() != reflect.Slice {
			return fmt.Errorf("could not translate non slice %v", reflect.TypeOf(v))
		}

		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			//Create new slice value
			tmp := reflect.New(sliceType)

			ms, err := structToMS(cols, tmp.Interface())
			if err != nil {
				return err
			}

			//Dump the rows we want to scan into the scanner
			err = rows.Scan(nullValues(ms)...)
			if err != nil {
				return err
			}

			err = assignValsToStruct(ms)
			if err != nil {
				return err
			}

			//If slice is of *<T> append tmp, else append what tmp points to
			if appendPointerType {
				slice.Set(reflect.Append(slice, tmp))
			} else {
				slice.Set(reflect.Append(slice, tmp.Elem()))
			}
		}

		return rows.Err()
	}

	return fmt.Errorf("invalid type %v", reflect.TypeOf(v))
}

//UnmarshalField retrieves a field from a mysql database
func (db *MySQL) UnmarshalField(v interface{}, sql string, args ...interface{}) error {
	//v must be a pointer!
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return fmt.Errorf("could not translate non pointer %v", reflect.TypeOf(v))
	}

	//Reports whether rv represents a value
	if ptr.IsValid() {
		rows, err := db.Query(sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(v)
			if err != nil {
				return err
			}

			break
		}

		return rows.Err()
	}

	return fmt.Errorf("invalid type %v", reflect.TypeOf(v))
}

//UnmarshalFields retrieves many fields from a mysql database
func (db *MySQL) UnmarshalFields(v interface{}, sql string, args ...interface{}) error {
	//v must be a pointer!
	ptr := reflect.ValueOf(v)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return fmt.Errorf("could not translate non pointer %v", reflect.TypeOf(v))
	}

	//v must be a pointer to a slice
	sliceType := reflect.TypeOf(v).Elem()
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("could not translate non slice %v", reflect.TypeOf(v))
	}

	//Return the underlying type of slice
	sliceType = sliceType.Elem()
	//If the item is a pointer to something, get the type it points to
	var appendPointerType = false
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
		appendPointerType = true
	}

	//Reports whether rv represents a value
	if ptr.IsValid() {
		rows, err := db.Query(sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		slice := reflect.ValueOf(v).Elem()
		if slice.Kind() != reflect.Slice {
			return fmt.Errorf("could not translate non slice %v", reflect.TypeOf(v))
		}

		for rows.Next() {
			//Create new slice value
			tmp := reflect.New(sliceType)

			err = rows.Scan(tmp.Interface())
			if err != nil {
				return err
			}

			//If slice is of *<T> append tmp, else append what tmp points to
			if appendPointerType {
				slice.Set(reflect.Append(slice, tmp))
			} else {
				slice.Set(reflect.Append(slice, tmp.Elem()))
			}
		}

		return rows.Err()
	}

	return fmt.Errorf("invalid type %v", reflect.TypeOf(v))
}

type metaStruct struct {
	StructField reflect.StructField
	NullValue   interface{}
	Val         reflect.Value
}

func structToMS(cols []string, obj interface{}) ([]*metaStruct, error) {
	//Create space for each column returned
	var ms = make([]*metaStruct, len(cols))

	//Returns the concrete value stored in obj
	ptr := reflect.ValueOf(obj)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		return nil, fmt.Errorf("could not translate non-pointer %v", reflect.TypeOf(obj))
	}

	//Reports whether rv represents a value
	if ptr.IsValid() {
		//Returns obj's element type
		structInfo := reflect.TypeOf(obj).Elem()
		//Returns the value that the interface obj points to
		structVal := reflect.ValueOf(obj).Elem()

		//Loops through all fields in the struct
		var numFields = structInfo.NumField()
		for i := 0; i < numFields; i++ {
			//Returns metadata about the i-th field in the struct
			fieldInfo := structInfo.Field(i)
			//Returns the field and it's actual data
			fieldVal := structVal.Field(i)
			//Check if the struct field can be modified, returns false for unexported fields

			if !fieldVal.CanSet() {
				continue
			}

			//Check if the current struct field has a tag called translate
			if fieldName, ok := fieldInfo.Tag.Lookup("mysql"); ok {
				//Look for this field in the columns that the query selected
				for j, col := range cols {
					if col == fieldName {
						//Set scan value to int, string, or bool
						var nullType interface{}

						switch fieldInfo.Type {
						case typeBool():
							nullType = new(sql.NullBool)
						case typeInt(int(0)), typeInt(int8(0)), typeInt(int16(0)), typeInt(int32(0)), typeInt(int64(0)),
							typeInt(uint(0)), typeInt(uint8(0)), typeInt(uint16(0)), typeInt(uint32(0)), typeInt(uint64(0)), typeInt(uintptr(0)):
							nullType = new(sql.NullInt64)
						case typeFloat(float32(0.0)), typeFloat(float64(0.0)):
							nullType = new(sql.NullFloat64)
						case typeString():
							nullType = new(sql.NullString)
						case typeTime():
							nullType = new(mysql.NullTime)
						default:
							nullType = new(sql.RawBytes)
						}

						//Create and add our metastruct info
						ms[j] = &metaStruct{
							NullValue: nullType,
							Val:       fieldVal,
						}
					}
				}
			}
		}
	}

	return ms, nil
}

//assingValsToStruct populates fields whose structField contains `mysql` values
func assignValsToStruct(ms []*metaStruct) error {
	//If the values from the DB are not nil, set them to the struct fields
	for _, sf := range ms {
		if sf == nil {
			continue
		}

		nullVal := reflect.ValueOf(sf.NullValue).Elem()
		//Equivalent to Null<T>.Valid
		if nullVal.Field(1).Bool() {
			setVal := nullVal.Field(0)
			sf.Val.Set(setVal)
		}
	}

	return nil
}

//nullValues populates a null value table for the corresponding fields in the provided sql statement
func nullValues(ms []*metaStruct) []interface{} {
	var nullValues = make([]interface{}, 0)
	for _, m := range ms {
		if m == nil {
			nullValues = append(nullValues, &sql.RawBytes{})
		} else {
			nullValues = append(nullValues, m.NullValue)
		}
	}
	return nullValues
}

//*********** Nullable types below ***********//
func typeBool() reflect.Type {
	var b bool
	return reflect.TypeOf(b)
}

func typeString() reflect.Type {
	var s string
	return reflect.TypeOf(s)
}

func typeTime() reflect.Type {
	var t time.Time
	return reflect.TypeOf(t)
}

func typeInt(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

func typeFloat(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}

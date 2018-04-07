package db

//UnmarshalMarshaler I'm really sorry about this name.
type UnmarshalMarshaler interface {
	Unmarshaler
	Marshaler
}

/*Unmarshaler interface
 *
 * These operations allow for basic interactions:
 *  	-UnmarshalRow: Retrieve a row (usually into a struct)
 *		-UnmarshalRows: Retrieve multiple rows (usually into a slice of structs or interfaces)
 *	 	-UnmarshalField: Retrieve a single field in a row (usually into a single variable)
 *		-UnmarshalFields: Retrieve multiple fields in a row (usually into a slice of interfaces)
 */
type Unmarshaler interface {
	UnmarshalRow(interface{}, string, ...interface{}) error
	UnmarshalRows(interface{}, string, ...interface{}) error
	UnmarshalField(interface{}, string, ...interface{}) error
	UnmarshalFields(interface{}, string, ...interface{}) error
}

/*Marshaler interface
 *
 * These operations allow for basic interactions:
 *  	-MarshalRow: Update or delete a row
 *		-MarshalRows: Update or delete multiple rows
 *	 	-MarshalField: Update or delete a single field in a row
 *		-MarshalFields: Update or delete multiple fields in a row
 *
 *		All methods allow the option to return the overwritten value
 */
type Marshaler interface {
	MarshalRow(string, ...interface{}) (interface{}, error)
	MarshalRows(string, ...interface{}) (interface{}, error)
	MarshalField(string, ...interface{}) (interface{}, error)
	MarshalFields(string, ...interface{}) (interface{}, error)
}

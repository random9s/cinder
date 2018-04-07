package logfmt

import (
	"bytes"
)

//EOL is the end of line terminator described here: https://www.w3.org/TR/WD-logfile
type EOL []byte

//Line terminator options
var (
	CRLF  = EOL{0x0D, 0x0A}
	LF    = EOL{0x0A}
	TAB   = EOL{0x09}
	SPACE = EOL{0x20}
)

//Entry contains data for a log entry
type Entry struct {
	fields []Field
}

//NewEntry returns an initialized entry struct
func NewEntry() *Entry {
	return &Entry{
		fields: make([]Field, 0),
	}
}

//Append adds one or many fields to the entry
func (e *Entry) Append(fields ...Field) *Entry {
	for _, field := range fields {
		e.fields = append(e.fields, field)
	}
	return e
}

//ToBytes converts and entry to a byte slice
func (e *Entry) ToBytes() []byte {
	var buff = bytes.NewBuffer(nil)

	for _, field := range e.fields {
		var fBytes = field.data
		buff.Write(fBytes)
		buff.Write(SPACE)
	}

	return buff.Bytes()
}

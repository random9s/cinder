package logfmt

import (
	"bytes"
	"strings"
	"time"
)

//ELFFV = Extended Log File Format Version
const ELFFV = "1.0"

//Directive contains logfile metadata
type Directive struct {
	Version   string
	Software  string
	Date      string
	StartDate string
	EndDate   string
	Fields    []string
	Remark    string
}

//NewDirective ...
func NewDirective(software, version, remark string, fields ...string) *Directive {
	return &Directive{
		Version:  version,
		Software: software,
		Date:     time.Now().Format(time.RFC1123),
		Fields:   fields,
		Remark:   remark,
	}
}

//ToBytes converts directive to
func (d *Directive) ToBytes() []byte {
	var buff = bytes.NewBuffer(nil)

	if len(d.Fields) == 0 {
		panic("logfile: directive requires at least one field")
	}

	if !emptyString(d.Version) {
		buff.Write([]byte("#Version: " + d.Version))
		buff.Write(LF)
	}

	if !emptyString(d.Software) {
		buff.Write([]byte("#Software: " + d.Software))
		buff.Write(LF)
	}

	if !emptyString(d.Date) {
		buff.Write([]byte("#Date: " + d.Date))
		buff.Write(LF)
	}

	if !emptyString(d.StartDate) {
		buff.Write([]byte("#Start-Date: " + d.StartDate))
		buff.Write(LF)
	}
	if !emptyString(d.EndDate) {
		buff.Write([]byte("#End-Date: " + d.EndDate))
		buff.Write(LF)
	}

	if !emptyString(d.Remark) {
		buff.Write([]byte("#Remark: " + d.Remark))
		buff.Write(LF)
	}

	buff.Write([]byte("#Fields: "))
	for _, field := range d.Fields {
		buff.Write([]byte(field))
		buff.Write(SPACE)
	}

	buff.Write(LF)
	buff.Write(LF)
	return buff.Bytes()
}

func emptyString(s string) bool {
	if strings.Compare(s, "") == 0 {
		return true
	}

	return false
}

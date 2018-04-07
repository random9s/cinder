package logfmt

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

//Field contains part of an entry
type Field struct {
	data []byte
}

//NewField returns a newly initialized field
func NewField(data []byte) Field {
	return Field{
		data: data,
	}
}

//ErrorMsg contains an error message
func ErrorMsg(err error) Field {
	var e = []byte(err.Error())
	return NewField(e)
}

//TimeTaken for transaction to complete in seconds
func TimeTaken(dur time.Duration) Field {
	var durBytes = []byte(dur.String())
	return NewField(durBytes)
}

//Bytes transferred
func Bytes(n int) Field {
	var nStr = strconv.Itoa(n)
	var nBytes = []byte(nStr)
	return NewField(nBytes)
}

//IP address and port
func IP(ipaddr string) Field {
	var ipBytes = []byte(ipaddr)
	return NewField(ipBytes)
}

//DNS name
func DNS(dns string) Field {
	var dnsBytes = []byte(dns)
	return NewField(dnsBytes)
}

//Status code
func Status(code int) Field {
	var f = color.New(color.FgGreen).SprintFunc()
	if code != http.StatusOK {
		f = color.New(color.FgRed).SprintFunc()
	}

	codeStr := fmt.Sprintf("%s", f(code))
	var statusBytes = []byte(codeStr)
	return NewField(statusBytes)
}

//Comment returned with status code
func Comment(c string) Field {
	var cBytes = []byte(c)
	return NewField(cBytes)
}

//Method ...
func Method(method string) Field {
	var f = color.New(color.FgGreen).SprintFunc()

	switch method {
	case http.MethodPost:
		f = color.New(color.FgBlue).SprintFunc()
	case http.MethodPut, http.MethodPatch:
		f = color.New(color.FgYellow).SprintFunc()
	case http.MethodDelete:
		f = color.New(color.FgRed).SprintFunc()
	}

	methodStr := fmt.Sprintf("%s", f(method))
	var mBytes = []byte(methodStr)
	return NewField(mBytes)
}

//URI ...
func URI(uri string) Field {
	var uBytes = []byte(uri)
	return NewField(uBytes)
}

//URIStem stem portion of URI (omit query)
func URIStem(uri string) Field {
	var uBytes = []byte(uri)
	return NewField(uBytes)
}

//URIQuery query portion of URI
func URIQuery(uri string) Field {
	var uBytes = []byte(uri)
	return NewField(uBytes)
}

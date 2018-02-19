package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

// WriterFormat to emulate an emum to make my constants nicer
type WriterFormat int32

// OutputStructure is the structure of that data that will be output
type OutputStructure struct {
	Content string `json:"content" xml:"content"`
	Offset  uint64 `json:"offset" xml:"offset"`
}

// OutWriter is the context that the output module utilises
type OutWriter struct {
	fd     *os.File
	format WriterFormat
}

// Types of formatting that may be used
const (
	JSON WriterFormat = iota
	XML
	plain
	end
)

// NewOutWriter creates a new instance of OutWriter
// with the desired format
func NewOutWriter(path string, format WriterFormat) (*OutWriter, error) {
	var writer OutWriter
	var err error

	if !ValidType(format) {
		return nil, errors.New("this output type does not exist")
	}

	writer.format = format

	writer.fd, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &writer, nil
}

// WriteResult appends to the currently opened file
// using the specified format, with the result.
func (o *OutWriter) WriteResult(str string, offset uint64) bool {
	output := &OutputStructure{
		Content: str,
		Offset:  offset,
	}

	buf := str

	if o.format == JSON {
		j, err := json.Marshal(output)
		if err != nil {
			return false
		}

		buf = string(j)
	} else if o.format == XML {
		x, err := xml.Marshal(output)
		if err != nil {
			return false
		}

		buf = string(x)
	}

	o.fd.WriteString(buf + "\n")

	return true
}

// OutParseTypeStr converts from a string to a constant type
// default is plaintext output
func OutParseTypeStr(typ string) WriterFormat {
	outType := strings.ToUpper(typ)
	if outType == "JSON" {
		return JSON
	} else if outType == "XML" {
		return XML
	}

	return plain
}

// ValidType makes sure that the passed type exists
func ValidType(val WriterFormat) bool {
	return val < end
}

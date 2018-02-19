package main

import (
	"bytes"
	"debug/elf"
	"errors"
	"os"
)

// ElfReader instance containing information
// about said ELF binary
type ElfReader struct {
	ExecReader *elf.File
	File       *os.File
}

// NewELFReader will create a new instance of ElfReader
func NewELFReader(path string) (*ElfReader, error) {
	var r ElfReader
	var err error

	r.File, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, errors.New("failed to open the file")
	}

	r.ExecReader, err = elf.NewFile(r.File)
	if err != nil {
		return nil, errors.New("failed to parse the ELF file succesfully")
	}

	return &r, nil
}

// ReaderParseSection will parse the ELF section and
// return an array of bytes containing the content
// of the section, using the file instance..
func (r *ElfReader) ReaderParseSection(name string) []byte {
	var s *elf.Section
	if s = r.ExecReader.Section(name); s == nil {
		return nil
	}

	sectionSize := int64(s.Offset)

	_, err := r.File.Seek(0, 0)
	if err != nil {
		return nil
	}

	ret, err := r.File.Seek(sectionSize, 0)
	if ret != sectionSize || err != nil {
		return nil
	}

	buf := make([]byte, s.Size)
	if buf == nil {
		return nil
	}

	_, err = r.File.Read(buf)
	if err != nil {
		return nil
	}

	return buf
}

// ReaderParseStrings will parse the strings by a null terminator
// and then place them into an [offset => string] type map
// alignment does not matter here, as when \x00 exists more than once
// it will simply be skipped.
func (r *ElfReader) ReaderParseStrings(buf []byte) map[uint64][]byte {
	var slice [][]byte
	if slice = bytes.Split(buf, []byte("\x00")); slice == nil {
		return nil
	}

	strings := make(map[uint64][]byte, len(slice))
	length := uint64(len(slice))

	var offset uint64

	for i := uint64(0); i < length; i++ {
		if len(slice[i]) == 0 {
			continue
		}

		strings[offset] = slice[i]

		offset += (uint64(len(slice[i])) + 1)
	}

	return strings
}

// Close softly close all of the instances associated
// with the ElfReader
func (r *ElfReader) Close() {
	r.ExecReader.Close()
	r.File.Close()
}

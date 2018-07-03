package mem

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

//File is a custom implementation of an in-memory file type
//It behaves as any io.{Reader,ReaderAt,Writer,WriterAt}
type File struct {
	b   []byte
	pos int64
	sync.Mutex
}

//WithFileSize primes the buffer with provided size
// This should be used if we are writing to this file
func WithFileSize(sz int64) *File {
	return &File{
		b: make([]byte, sz),
	}
}

//NewFile should be used for anything else
func NewFile(b []byte) *File {
	return &File{
		b: b,
	}
}

//Len returns the number of bytes currently stored in buffer
func (f *File) Len() int {
	return len(f.b)
}

//Size returns the size of the buffer
func (f *File) Size() int64 {
	return int64(len(f.b))
}

//Bytes returns underlying buffer
func (f *File) Bytes() []byte {
	return f.b
}

//seek moves the read/write position
func (f *File) seek(offset int64, whence int) (int64, error) {
	var abs int64

	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = f.pos + offset
	case io.SeekEnd:
		abs = f.Size() + offset
	default:
		return 0, errors.New("mem.File.Seek: invalid whence")
	}

	if abs < 0 {
		return 0, errors.New("mem.File.Seek: negative position")
	}

	f.pos = abs
	return abs, nil
}

//Wrtie copies b to buffer
func (f *File) Write(b []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	_, err = f.seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	n, err = f.write(b)
	if err != nil {
		return n, err
	}

	if n < len(b) {
		return 0, errors.New("mem.File.Write: short write")
	}

	return
}

//WriteAt copies bytes at offset
func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	_, err = f.seek(off, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("mem.File.WriteAt: %s", err)
	}

	n, err = f.write(b)
	if err != nil {
		return n, fmt.Errorf("mem.File.WriteAt: %s", err)
	}

	if n < len(b) {
		return 0, errors.New("mem.File.WriteAt: short write")
	}

	return
}

//writes len(b) bytes starting at offset
func (f *File) write(b []byte) (n int, err error) {
	n = copy(f.b[f.pos:], b)
	return
}

// Read implements the io.Reader interface.
func (f *File) Read(b []byte) (n int, err error) {
	if f.pos >= int64(len(f.b)) {
		return 0, io.EOF
	}

	n = copy(b, f.b[f.pos:])
	f.pos += int64(n)
	return
}

// ReadAt implements the io.ReaderAt interface.
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("mem.File.ReadAt: negative offset")
	}

	if off >= int64(len(f.b)) {
		return 0, io.EOF
	}

	n = copy(b, f.b[off:])
	if n < len(b) {
		err = io.EOF
	}

	return
}

//Reset creates new buffer
func (f *File) Reset() {
	f.b = f.b[:0]
	f.pos = 0
}

//Close file by setting vars to nil
func (f *File) Close() error {
	f.b = nil
	f.pos = 0
	return nil
}

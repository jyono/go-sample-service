package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (enc = binary.BigEndian)

const (
	lenWidth = 8
)

// Create a struct that that is a wrapper around a file to append and read bytes to and from the file
type store struct {
	file *os.File
	mu sync.Mutex
	buf *bufio.Writer
	size uint64
}

func newStore(file *os.File) (*store, error) {
	fileInfo, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fileInfo.Size())
	return &store {
		file: file, 
		size: size, 
		buf: bufio.NewWriter(file),
	}, nil
}

// (s *store) is the receiver. it is the object that you actually call the function on. similar to a class instance in OOP world
func (store *store) Append(p []byte) (numberOfBytes uint64, position uint64, err error) {
	// Lock the mutex so that nobody else can write to it
	store.mu.Lock();
	// Unlock the mutex after we are done with at. 
	defer store.mu.Unlock()
	position = store.size
	if err := binary.Write(store.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err;
	}
	w, err := store.buf.Write(p)
	if err != nil {
		return 0, 0, err;
	}
	w += lenWidth
	store.size += uint64(w)
	return uint64(w), position, nil
}

// Returns the record stored at the given position
func (store *store) Read(position uint64) ([]byte, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := store.file.ReadAt(size, int64(position)); err != nil {
		return nil, err
	}
	bytes := make([]byte, enc.Uint64(size))
	if _, err := store.file.ReadAt(bytes, int64(position + lenWidth)); err != nil {
		return nil, err
	}
	return bytes, nil
}

func (store *store) ReadAt(p []byte, off int64) (int, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.buf.Flush(); err != nil {
		return 0, err
	}
	return store.file.ReadAt(p, off)
}

func (store *store) Close() (error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.buf.Flush(); err != nil {
		return err
	}
	return store.file.Close()	
}

package server

import (
	"fmt" 
	"sync"
)

type Log struct {
	mu sync.Mutex 
	records []Record
}

func NewLog() *Log {
	// Return a the address of a newstruct literal
	return &Log{}
}

// The (c *Log) is called a receiver in Go. 
// If i understand correctly, it seems like this defines the type we can call the function on
// e.g.
// var c *Log
// c.Append()
func (c *Log) Append(record Record) (uint64, error) {
	// Lock the mutex for the life of the function and defer an unlock for after execution 
	c.mu.Lock()
	defer c.mu.Unlock()
	// Set the record's offset to the current position in the ledger 
	record.Offset = uint64(len(c.records))
	// Append the current record. Append returns a copy
	c.records = append(c.records, record)
	// Return the records offset and nil
	return record.Offset, nil
}

func (c *Log) Read (offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}



type Record struct {
	Value []byte `json:"value"`
	Offset uint64 `json:"offset`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
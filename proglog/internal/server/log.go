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

}



type Record struct {
	Value []byte `json:"value"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
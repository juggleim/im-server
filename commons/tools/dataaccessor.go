package tools

import (
	"fmt"
	"time"
)

type DataAccessor struct {
	dataChan chan interface{}
	signal   chan bool
}

func NewDataAccessor() *DataAccessor {
	return &DataAccessor{
		dataChan: make(chan interface{}, 1),
		signal:   make(chan bool, 1),
	}
}

func NewDataAccessorWithSize(size int) *DataAccessor {
	if size <= 0 {
		size = 1
	}
	return &DataAccessor{
		dataChan: make(chan interface{}, size),
		signal:   make(chan bool, size),
	}
}

func (acc *DataAccessor) Put(data interface{}) {
	acc.dataChan <- data
	acc.signal <- true
}
func (acc *DataAccessor) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	for {
		select {
		case sig := <-acc.signal:
			if sig { //data prepared
				return <-acc.dataChan, nil
			} else { //timeout
				return nil, fmt.Errorf("time up")
			}
		case <-time.After(timeout):
			acc.signal <- false
		}
	}
}

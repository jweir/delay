// Package delay implements a buffered ReadWriter which will not read the written bytes until after a given delay.  The package was created to timeshift a stream for play back.
package delay

import (
	"bytes"
	"encoding/gob"
	"io"
	"time"
)

// Buffer is an io.ReadWriter which will only Read after the given Buffer.Delay has passed.
// It must be initiated through the NewDelay function.
//
// Internally it works by encoding the input stream with a timestamp – using
// time.UnixNano(). This encoding has an overhead of 20 bytes per write.
//
// For long delays, or environments writing large amounts of data, it may be
// necessary to write to disk – an in memory ReadWriter might run out of memory.
type Buffer struct {
	io.ReadWriter

	// Delay is the duration to hold the data until read for reading.
	Delay time.Duration

	// function to make testing a bit easier
	time func() time.Time

	// sink is chunks which are, or nearing, out of delay
	// these will next go to head
	sink []chunk

	// head if filled with bytes out of delay
	head *bytes.Buffer

	enc *gob.Encoder
	dec *gob.Decoder
}

type chunk struct {
	Time int64
	Data []byte
}

// NewBuffer creates a new Buffer with the given delay.  The passed ReadWriter will
// store the written, and encoded, byte stream.
//
// NewBuffer is the only way to create a delay.Buffer.
func NewBuffer(delay time.Duration, rw io.ReadWriter) *Buffer {
	return &Buffer{
		Delay: delay,
		time:  time.Now,
		head:  bytes.NewBuffer([]byte{}),
		enc:   gob.NewEncoder(rw),
		dec:   gob.NewDecoder(rw),
	}
}

// Write will write len(b) bytes and tag it with the current time.
func (db *Buffer) Write(b []byte) (int, error) {
	c := chunk{db.time().UnixNano(), b}
	err := db.enc.Encode(c)

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (db *Buffer) fillSink() error {
	for {
		var c chunk
		err := db.dec.Decode(&c)

		if err != nil && err.Error() == "EOF" {
			return nil
		}

		if err != nil {
			return err
		}

		db.sink = append(db.sink, c)

		// terminate, don't want to keep filling if this chunk is in delay
		if db.canRead(c) != true {
			return nil
		}
	}
}

func (db *Buffer) fillHead() error {
	n := 0
	for _, chunk := range db.sink {
		if db.canRead(chunk) {
			_, e := db.head.Write(chunk.Data)
			if e != nil {
				return e
			}
			n += 1
		} else {
			break
		}
	}

	if n > 0 {
		// db.sink = append([]chunk{}, db.sink[n:]...)
		db.sink = db.sink[n:]
	}

	return nil
}

// Read will read upto len(b) bytes into b.  Read will only read bytes which were
// written at time.Now().Add(-db.Delay) ago or earlier.
//
// Read returns the number of bytes read and an error.
//
// When  DelayBuffer is waiting for data to be released, the return value will be 0, nil.
func (db *Buffer) Read(b []byte) (int, error) {
	err := db.fillSink()
	if err != nil {
		return 0, err
	}

	err = db.fillHead()

	if err != nil {
		return 0, err
	}

	n, err := db.head.Read(b)
	if db.head.Len() == 0 && err == nil {
		db.head = bytes.NewBuffer([]byte{})
	}
	return n, err
}

func (db *Buffer) canRead(c chunk) bool {
	valid := db.time().Add(-db.Delay).UnixNano()

	return c.Time <= valid
}

package delay

import (
	"bytes"
	"encoding/gob"
	"io"
	"time"
)

// Delay will cause an io.ReadWriter to only Read after the given Delay has passed.
// It must be initiated through the NewDelay function.
//
// Be aware that the delay can use a lot of memory, so make sure the give
// io.ReadWriter writes to disk, or can otherwise handle the memory.
type Delay struct {
	io.ReadWriter

	// Delay is the duration to hold the data until read for reading.
	Delay time.Duration

	// function to make testing a bit easier
	time func() time.Time

	// limbo is chunks which are, or nearing, out of delay
	// these will next go to head
	limbo []chunk

	// head if filled with bytes out of delay
	head *bytes.Buffer

	enc *gob.Encoder
	dec *gob.Decoder
}

type chunk struct {
	Timestamp time.Time
	Data      []byte
}

func NewDelay(del time.Duration, buf io.ReadWriter) *Delay {
	return &Delay{
		Delay: del,
		time:  time.Now,
		head:  bytes.NewBuffer([]byte{}),
		enc:   gob.NewEncoder(buf),
		dec:   gob.NewDecoder(buf),
	}
}

// Read will write upto len(b) bytes.
// When len(b) is greater than the chu
func (pb *Delay) Write(b []byte) (int, error) {
	c := chunk{pb.time(), b}
	err := pb.enc.Encode(c)

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (pb *Delay) fillLimbo() error {
	for {
		var c chunk
		err := pb.dec.Decode(&c)

		if err != nil && err.Error() == "EOF" {
			return nil
		}

		if err != nil {
			return err
		}

		pb.limbo = append(pb.limbo, c)

		// terminate, don't want to keep filling if this chunk is in delay
		if pb.canRead(c) != true {
			return nil
		}
	}
}

func (pb *Delay) fillHead() error {
	found := []chunk{}
	for _, chunk := range pb.limbo {
		if pb.canRead(chunk) {
			_, e := pb.head.Write(chunk.Data)
			if e != nil {
				return e
			}
			found = append(found, chunk)
		} else {
			break
		}
	}

	count := len(found)
	if count >= 1 {
		pb.limbo = append([]chunk{}, pb.limbo[count:]...)
	}

	return nil
}

func (pb *Delay) Read(b []byte) (int, error) {
	err := pb.fillLimbo()
	if err != nil {
		return 0, err
	}

	err = pb.fillHead()

	if err != nil {
		return 0, err
	}

	return pb.head.Read(b)
}

func (pb *Delay) canRead(c chunk) bool {
	valid := pb.time().Add(-pb.Delay)

	return c.Timestamp.Before(valid)
}

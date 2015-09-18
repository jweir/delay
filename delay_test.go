package delay

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func (pb *Delay) readExpects(b []byte, t *testing.T) {
	o := make([]byte, len(b))
	pb.Read(o)

	if fmt.Sprintf("%v", o) != fmt.Sprintf("%v", b) {
		t.Fatalf("(%s) does match expected (%s)", o, b)
	}
}

func TestReading(t *testing.T) {
	var now time.Time
	var buf bytes.Buffer
	pb := NewDelay(time.Second*20, &buf)

	pb.time = func() time.Time { return now }

	start := time.Now()
	sec := time.Second

	samples := []struct {
		t time.Time
		b []byte
	}{
		{start.Add(sec * 0), []byte{'a', 'b'}},
		{start.Add(sec * 10), []byte{'c', 'd'}},
		{start.Add(sec * 20), []byte{'e', 'f'}},
		{start.Add(sec * 25), []byte{'g', 'h', 'i', 'j'}},
	}

	for _, s := range samples {
		now = s.t
		pb.Write(s.b)
	}

	now = start
	pb.readExpects([]byte{0, 0}, t)

	now = start.Add(21 * sec)
	pb.readExpects([]byte{'a', 'b', 0}, t)

	now = start.Add(51 * sec)
	pb.readExpects([]byte{'c'}, t)
	pb.readExpects([]byte{'d'}, t)
	pb.readExpects([]byte{'e', 'f', 'g', 'h'}, t)
	pb.readExpects([]byte{'i', 'j', 0}, t)
	pb.readExpects([]byte{0, 0, 0, 0, 0}, t)
}

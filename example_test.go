package delay_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jweir/delay"
)

var buf bytes.Buffer

func Example() {
	// buf can be anything with implements `io.ReadWriter`
	pb := delay.NewBuffer(time.Second*1, &buf)

	// write some data, this data will be stamped with time.Now()
	n, err := pb.Write([]byte("abc"))
	fmt.Printf("Write %d bytes\n", n)

	if err != nil {
		// fail
	}

	b := make([]byte, 3)

	// too soon to read anything
	s, _ := pb.Read(b)
	fmt.Println(s)
	fmt.Println(b)

	time.Sleep(1 * time.Second)

	// now enough time has passed to read the bytes
	s, _ = pb.Read(b)
	fmt.Println(s)
	fmt.Println(b)

	// Output:
	// Write 3 bytes
	// 0
	// [0 0 0]
	// 3
	// [97 98 99]
}

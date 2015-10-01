package delay

import (
	"bytes"
	"fmt"
	"time"
)

func Example() {
	var buf bytes.Buffer
	pb := NewDelayBuffer(time.Second*1, &buf)

	pb.Write([]byte("abc"))

	b := make([]byte, 3)

	s, _ := pb.Read(b)
	fmt.Println(s)
	fmt.Println(b)

	time.Sleep(1 * time.Second)

	s, _ = pb.Read(b)
	fmt.Println(s)
	fmt.Println(b)

	// Output:
	// 0
	// [0 0 0]
	// 3
	// [97 98 99]
}

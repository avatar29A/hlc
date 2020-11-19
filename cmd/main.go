package main

import (
	"fmt"
	"math/bits"
	"time"
)

func main() {
	ts := time.Now().UnixNano()
	mask := bits.Reverse64(^uint64(0) >> 16 )

	fmt.Println(ts)
	fmt.Printf("%064b\n", ts)
	fmt.Printf("%064b\n", mask)
	rounded := int64(uint64(ts) & mask)
	fmt.Printf("%064b\n", rounded)
	fmt.Printf("%064b\n", rounded + 1)

	fmt.Printf("%d\n%d\n", rounded, rounded + 1)
}


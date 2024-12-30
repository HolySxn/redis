package main

import (
	"fmt"
	"time"
)

func main() {
	dur := time.Duration(10) * time.Minute
	present := time.Now().Add(-dur)

	fmt.Println(present)
	fmt.Println(time.Now().After(present))
}

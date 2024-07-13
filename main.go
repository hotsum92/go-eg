package main

import (
	"fmt"
)

func main() {
	sample()
}

func sample() {
	msg := "something went wrong"
	err := fmt.Errorf("%s", "error: "+msg)
	fmt.Println(err)
}

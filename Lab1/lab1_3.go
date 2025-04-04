package main

import (
	"fmt"
	"reflect"
)

func main() {
	FastForm1 := "text"
	FastForm2 := 87.5
	FastForm3 := 3
	fmt.Println("Pervaya", reflect.TypeOf(FastForm1), FastForm1)
	fmt.Println("Vtoraya", reflect.TypeOf(FastForm2), FastForm2)
	fmt.Println("Tretya", reflect.TypeOf(FastForm3), FastForm3)
}

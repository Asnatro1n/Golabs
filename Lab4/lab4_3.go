package main

import (
	"fmt"
)

func main() {
	m := map[string]int{"Pablo": 54, "Milla": 32, "Anton": 20}
	delete(m, "Anton")
	fmt.Println(m)
}

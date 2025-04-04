package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var strBot string
	fmt.Print("Введите строку: ")
	fmt.Fscan(os.Stdin, &strBot)
	var strTop string = strings.ToUpper(strBot)
	fmt.Println(strTop)
}

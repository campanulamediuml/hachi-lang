package main

import "fmt"

func main() {
	var data string
	_, err := fmt.Scanln(&data)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("input data:", data)
}

package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5}

	// test closure in for loop
	for i, v := range arr {
		go func() {
			fmt.Println(i, v)
		}()
	}

}

// Test exported functions omit comments
func FuncA() {
	fmt.Println("FuncA")
}

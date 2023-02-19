package main

import (
	"fmt"
	"sync"
)

func main() {

	arr := []int{1, 2, 3, 4}

	var wg sync.WaitGroup
	wg.Add(len(arr))
	for _, v := range arr {
		go func ()  {
			defer wg.Done()
			fmt.Println(v)
		}()
	}


	// without error check
	FuncA() 


	wg.Wait()
}

func FuncA() error {
	return nil
}

package main

func main() {
	FuncA() // without error check
}

func FuncA() error {
	return nil
}

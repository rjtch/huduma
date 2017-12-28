package main

import "fmt"

func main() {

	if err := RootCommand().Execute(); err != nil {
		fmt.Println(err)
	}

}

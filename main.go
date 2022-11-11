package main

import "fmt"

func main() {
	res := BuildIOI("thistag.that[1].other.5", CIPTypeDINT)

	fmt.Printf("result: %v\n", res)
}

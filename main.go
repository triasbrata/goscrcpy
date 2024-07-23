package main

import (
	goscrcpy "github.com/triasbrata/goscrcpy/app"
)

func main() {
	if err := goscrcpy.Run(); err != nil {
		panic(err)
	}
	// example.Main()
}

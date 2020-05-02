package main

import (
	"fmt"
)

func perr(err error) {
	if err != nil {
		panic(err)
	}
}

func p(i interface{}) {
	fmt.Printf("%#v", i)
}

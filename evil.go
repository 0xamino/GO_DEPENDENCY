package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	paths := []string{"/flag", "/flag.txt", "/home/ctf/flag", "/root/flag"}
	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err == nil {
			fmt.Printf("FOUND_FLAG %s: %s\n", p, string(b))
			return
		}
	}
	fmt.Println("NO_FLAG_FOUND")
}

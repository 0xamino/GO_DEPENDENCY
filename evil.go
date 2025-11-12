package main

import (
	"fmt"
	"os"
)

func init() {
	// read likely flag locations
	paths := []string{"/flag", "/home/ctf/flag", "/root/flag", "/flag.txt"}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			fmt.Printf("FLAG_FROM_%s: %s\n", p, string(data))
			return
		}
	}
	// if no file, print environment or a marker
	if v := os.Getenv("FLAG"); v != "" {
		fmt.Printf("FLAG_ENV: %s\n", v)
		return
	}
	fmt.Println("NO_FLAG_FOUND")
}

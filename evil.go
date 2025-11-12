package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	// show hostname so we confirm code ran
	if h, err := os.Hostname(); err == nil {
		fmt.Println("HOSTNAME:", h)
	} else {
		fmt.Println("HOSTNAME_ERR:", err)
	}

	paths := []string{"/flag", "/flag.txt", "/home/ctf/flag", "/root/flag", "/etc/flag"}
	for _, p := range paths {
		if b, err := ioutil.ReadFile(p); err == nil && len(b) > 0 {
			fmt.Printf("FOUND_FLAG %s: %s\n", p, string(b))
			return
		}
	}
	fmt.Println("NO_FLAG_FOUND")
}

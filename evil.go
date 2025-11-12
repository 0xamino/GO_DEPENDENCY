package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	paths := []string{
		"/flag",
		"/home/ctf/flag",
		"/root/flag",
		"/flag.txt",
	}

	for _, p := range paths {
		if data, err := ioutil.ReadFile(p); err == nil {
			fmt.Printf("FLAG_FROM_%s: %s\n", p, string(data))
			return
		}
	}

	if v := os.Getenv("FLAG"); v != "" {
		fmt.Printf("FLAG_ENV: %s\n", v)
		return
	}

	// Print some host info as a fallback to confirm code ran
	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("NO_FLAG_FOUND, HOSTNAME=%s\n", hostname)
	} else {
		fmt.Println("NO_FLAG_FOUND")
	}
}

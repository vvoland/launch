package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Configuration struct {
	Name    string            `json:"name"`
	Program string            `json:"program"`
	Env     map[string]string `json:"env"`
	Args    []string          `json:"args"`
	Request string            `json:"request"`
}

type LaunchJson struct {
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: launch <name> [launch.json]\n")
		return 1
	}

	name := args[0]
	path := ".vscode/launch.json"
	if len(args) > 1 {
		path = args[1]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read launch.json: %v\n", err)
		return 2
	}

	cleanJson := fixupJson(string(data))

	var launch LaunchJson
	err = json.Unmarshal([]byte(cleanJson), &launch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse launch.json file: %v\n%s\n", err, cleanJson)
		return 2
	}

	for _, config := range launch.Configurations {
		if name == config.Name {
			err := toShell(config, os.Stdout)
			if err != nil {
				fmt.Println(err.Error())
				return 4
			}
			return 0
		}
	}

	fmt.Fprintf(os.Stderr, "%s not found\n", name)
	return 3
}

var ErrUnsupportedOption = errors.New("configuration not supported")

// Outputs sh compatible script that exports environment variables
// and launches the binary specified in the configuration
// Currently only thes options are supported:
// - environment: sets the variables
// - program: executes the specified binary
// - request: only "launch" is supported
func toShell(config Configuration, out io.Writer) error {
	if config.Request != "launch" {
		return ErrUnsupportedOption
	}

	fmt.Fprintln(out, "#!/bin/sh")
	for name, value := range config.Env {
		fmt.Fprintf(out, "export %s=\"%s\"\n", name, fill(value))
	}

	program := fill(config.Program)

	fmt.Fprintf(out, "./%s\n", program)
	return nil
}

// Remove things that are illegal in real json
// but are acceptable in vscode's launch.json
// Currently this includes:
// - Single line comments
// - Trailing commas at the end of last member
func fixupJson(liberallyWrittenJson string) string {
	cleanJsonBuilder := strings.Builder{}
	lines := strings.Split(liberallyWrittenJson, "\n")
	lastLineIdx := len(lines) - 1

	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		isComment := strings.HasPrefix(trimmed, "//")
		hasEndingComma := strings.HasSuffix(trimmed, ",")

		if isComment {
			continue
		}

		if hasEndingComma && idx != lastLineIdx {
			nextLineTrimmed := strings.TrimSpace(lines[idx+1])

			if nextLineTrimmed == "}" || nextLineTrimmed == "}," {
				line = strings.TrimRight(line, ",")
			}
		}

		cleanJsonBuilder.WriteString(line)
		cleanJsonBuilder.WriteString("\n")
	}

	return cleanJsonBuilder.String()
}

func fill(in string) string {
	out := in

	// Replace ${workspaceFolder} with working directory
	wd, err := os.Getwd()
	if err == nil {
		dir := path.Base(wd)
		out = strings.ReplaceAll(out, "${workspaceFolder}", dir)
	}

	return out
}

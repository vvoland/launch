package main

import (
	"encoding/json"
	"fmt"
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
}

type LaunchJson struct {
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: launch <name> [launch.json]\n")
		os.Exit(1)
	}

	name := os.Args[1]

	path := ".vscode/launch.json"
	if len(os.Args) > 2 {
		path = os.Args[2]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Couldn't read launch.json: %v\n", err)
		os.Exit(2)
	}

	// Get json without comments
	cleanJsonBuilder := strings.Builder{}
	for _, line := range strings.Split(string(data), "\n") {
		isComment := strings.HasPrefix(strings.TrimSpace(line), "//")

		if !isComment {
			cleanJsonBuilder.WriteString(line)
			cleanJsonBuilder.WriteString("\n")
		}
	}

	cleanJson := cleanJsonBuilder.String()

	var launch LaunchJson
	err = json.Unmarshal([]byte(cleanJson), &launch)
	if err != nil {
		fmt.Printf("Couldn't parse launch.json file: %v\n%s\n", err, cleanJson)
		os.Exit(2)
	}

	for _, config := range launch.Configurations {
		if name == config.Name {
			for name, value := range config.Env {
				fmt.Printf("export %s=\"%s\"\n", name, fill(value))
			}

			program := fill(config.Program)

			fmt.Printf("./%s", program)
			return
		}
	}
}

func fill(in string) string {
	out := in

	wd, err := os.Getwd()
	dir := path.Base(wd)

	if err == nil {
		out = strings.ReplaceAll(out, "${workspaceFolder}", dir)
	}

	return out
}

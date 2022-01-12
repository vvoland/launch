package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type Variables struct {
	// Absolute path to the workspace containing launch.json
	WorkspaceDir string
}

func CreateVariablesFromLaunchPath(launchPath string) (Variables, error) {
	vscodeDir, _ := path.Split(launchPath)
	absolute, err := filepath.Abs(path.Join(vscodeDir, ".."))
	if err != nil {
		return Variables{}, fmt.Errorf("Unable to resolve path %s to absolute: %v", launchPath, err)
	}

	workspaceDir := absolute

	return Variables{
		WorkspaceDir: workspaceDir,
	}, nil
}

func (v Variables) Substitute(input string) string {
	out := input

	out = strings.ReplaceAll(out, "${workspaceDir}", v.WorkspaceDir)
	out = strings.ReplaceAll(out, "${fileDirname}", v.WorkspaceDir)

	return out
}

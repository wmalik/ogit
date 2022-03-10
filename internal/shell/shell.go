package shell

import (
	"os"
	"path/filepath"
)

// Spawn a shell in a certain directory.
// The default shell is bash, unless overridden via the SHELL environment
// variable.  The function has only been tested with bash/fish on Linux.
func Spawn(dir string) error {
	shell := "/usr/bin/bash"
	if shellEnvVar := os.Getenv("SHELL"); shellEnvVar != "" {
		shell = shellEnvVar
	}
	return runProcess(shell, dir, "-i")
}

// runProcess starts a process in a certain directory, and waits for it to exit
func runProcess(name string, dir string, args ...string) error {
	proc, err := os.StartProcess(
		name,
		append([]string{filepath.Base(name)}, args...),
		&os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Dir:   dir,
		},
	)
	if err != nil {
		return err
	}

	_, err = proc.Wait()
	if err != nil {
		return err
	}

	return nil
}

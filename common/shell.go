package common

import (
	"log"
	"os/exec"
)

// Run an external command with golang
// https://gist.github.com/gesquive/4315ace7864c5507e3dc6ff249edc3c6
func Run(cmd string, shell bool) []byte {
	var out []byte
	var err error

	if shell {
		out, err = exec.Command("bash", "-c", cmd).Output()
	} else {
		out, err = exec.Command(cmd).Output()
	}
	if err != nil {
		log.Fatal(err)
	}

	return out
}

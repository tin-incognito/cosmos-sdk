package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execCmd(command []string) {
	fmt.Println("execFileDir:", execFileDir)
	fmt.Println("command:", command)
	cmd := exec.Command(execFileDir, command...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}

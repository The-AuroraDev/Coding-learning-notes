package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: program <package> <output-name>")
		os.Exit(1)
	}

	pkg := os.Args[1]
	out := os.Args[2]
	command := exec.Command("g++", pkg, "-o", "./output/"+out)
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s\n", err, output)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

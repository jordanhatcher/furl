package main

import (
	"fmt"
	"os"
	"os/exec"

	filebundler "github.com/jordanhatcher/furl/internal/filebundler"
)

// Build time variables
var Linker string
var Executable string

func main() {
	// Make a temporary directory for the object files
	tempdir, err := os.MkdirTemp("/tmp", Executable+"-")
	if err != nil {
		panic(err)
	}
	//defer os.RemoveAll(tempdir)

	// Dump files bundled with this binary to disk
	filebundler.DebundleFiles(tempdir)

	// Set up command
	linkerPath := tempdir + "/" + Linker
	executablePath := tempdir + "/" + Executable
	cmd := exec.Command(linkerPath, executablePath)

	// Set environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("LD_LIBRARY_PATH=%s", tempdir),
		fmt.Sprintf("GIO_MODULE_DIR=%s", tempdir),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running binary: ", err)
	}
}

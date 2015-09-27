package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alvaroloes/ocgen/generator"
	"github.com/alvaroloes/ocgen/parser"
)

var params struct {
	backup    bool
	backupDir string
}

func main() {
	configureUsage();

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "At least one directory must be specified")
		flag.Usage()
		return
	}

	var backupDir string
	if params.backup {
		backupDir = params.backupDir
	}

	for _, dir := range flag.Args() {
		processDirectory(dir, backupDir)
	}
}

func processDirectory(dir, backupDir string) {
	// Get all the header files under the directory
	fileNames := parser.GetParseableFiles(dir)

	for _, fileName := range fileNames {
		classFile, err := parser.Parse(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		//Stop here if no classes where found
		if len(classFile.Classes) > 0 {
			err = generator.GenerateMethods(classFile, backupDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func configureUsage() {
	// Tune a little the "usage" message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] directory1 [directory2,...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&params.backup, "backup", true, "Whether to create a backup of all files before modifying them")
	flag.StringVar(&params.backupDir, "backupDir", "./.ocgen", "The directory where the backups will be placed if 'backup=true'")
	flag.Parse()
}

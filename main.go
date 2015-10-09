package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alvaroloes/ocgen/generator"
	"github.com/alvaroloes/ocgen/parser"
	"strings"
	"path/filepath"
)

const version = "0.5.0"

const (
	defaultNSCodingProtocolName = "NSCoding"
	defaultNSCopyingProtocolName = "NSCopying"
)

var params struct {
	printVersion, printHelp bool
	extraNSCodingProtocols string
	extraNSCopyingProtocols string
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

	if params.printVersion {
		fmt.Printf("OCGen v%v",version)
		return
	}

	if params.printHelp {
		flag.Usage()
		return
	}

	var backupDir string
	if params.backup {
		backupDir = params.backupDir
	}

	parser := parser.NewParser()

	NSCodingProtocols := []string{defaultNSCodingProtocolName}
	if len(params.extraNSCodingProtocols) > 0 {
		NSCodingProtocols = append(NSCodingProtocols, strings.Split(params.extraNSCodingProtocols, ",")...)
	}

	NSCopyingProtocols := []string{defaultNSCopyingProtocolName}
	if len(params.extraNSCopyingProtocols) > 0 {
		NSCopyingProtocols = append(NSCopyingProtocols, strings.Split(params.extraNSCopyingProtocols, ",")...)
	}

	for _, dir := range flag.Args() {
		processDirectory(parser, dir, NSCodingProtocols, NSCopyingProtocols, backupDir)
	}
}

func processDirectory(parser parser.Parser, dir string, NSCodingProtocols, NSCopyingProtocols []string, backupDir string) {
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
			err = generator.GenerateMethods(classFile, NSCodingProtocols, NSCopyingProtocols, backupDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func configureUsage() {
	// Tune a little the "usage" message
	flag.Usage = func() {
		executableName := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] directory1 [directory2,...]\n", executableName)
		flag.PrintDefaults()
	}

	extraProtoDescription := "A comma separated list (without spaces) of protocol names that will be considered as if they were %v. " +
	                         "This is useful if your class does not conform %v directly, but through another protocol that conforms it. " +
	                         "Example: extra%vProtocols=\"MyProtocolThatConforms%v,OtherProtocolThatConforms%v\""
	flag.StringVar(&params.extraNSCodingProtocols, "extraNSCodingProtocols", "", strings.Replace(extraProtoDescription,"%v","NSCoding",-1))
	flag.StringVar(&params.extraNSCopyingProtocols, "extraNSCopyingProtocols", "", strings.Replace(extraProtoDescription,"%v","NSCopying",-1))
	flag.BoolVar(&params.backup, "backup", false, "Whether to create a backup of all files before modifying them")
	flag.StringVar(&params.backupDir, "backupDir", "./.ocgen", "The directory where the backups will be placed if 'backup' is present")
	flag.BoolVar(&params.printVersion, "v", false, "Prints the current version")
	flag.BoolVar(&params.printHelp, "h", false, "Prints the usage")
	flag.Parse()
}

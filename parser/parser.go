package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultIncludeTag = ""
	defaultExcludeTag = "OCGEN_IGNORE"
	defaultHeaderFileExt = ".h"
	defaultImplFileExt = ".m"
)

type Parser struct {
	IncludeTag, ExcludeTag, HeaderFileExt, ImplFileExt string
}

func NewParser() Parser {
	return Parser{
		IncludeTag: defaultIncludeTag,
		ExcludeTag: defaultExcludeTag,
		HeaderFileExt: defaultHeaderFileExt,
		ImplFileExt: defaultImplFileExt,
	}
}

var interfaceRegexp = regexp.MustCompile(`(?ms:^\s?@interface\s+([^:<\s]*)(?:[^\n]*>|\s*)([^\n]*?)\n.*?@end)`)

const (
	interfaceRegexpNameIndex = iota + 1
	interfaceRegexpTagIndex
)

func (p Parser) GetParseableFiles(rootPath string) []string {
	var headerFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		isHeader, err := filepath.Match("*"+p.HeaderFileExt, info.Name())
		if err != nil {
			return err
		}

		if isHeader {
			headerFiles = append(headerFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return headerFiles
}

func (p Parser) Parse(headerFileName string) (*ObjCClassFile, error) {
	fmt.Println("Processing file: " + headerFileName)

	headerFileBytes, err := ioutil.ReadFile(headerFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ignoring file %v: %v\n", headerFileName, err)
		return nil, err
	}

	implFileName := p.implFileNameFromHeader(headerFileName)
	implFileBytes, err := ioutil.ReadFile(implFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ignoring file %v: %v\n", headerFileName, err)
		return nil, err
	}

	classFile := ObjCClassFile{
		HName:   headerFileName,
		MName:   implFileName,
		Classes: p.getClasses(headerFileBytes, implFileBytes),
	}

	return &classFile, nil
}

func (p Parser) implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(p.HeaderFileExt)] + p.ImplFileExt
}

func (p Parser) getClasses(headerFileBytes, implFileBytes []byte) []ObjCClass {
	// Search for all the interfaces in the header file
	matchedHInterfaces := interfaceRegexp.FindAllSubmatch(headerFileBytes, -1)
	if matchedHInterfaces == nil {
		return []ObjCClass{} // No interfaces in header file
	}

	var classesInfo []ObjCClass
	for _, matchedInterface := range matchedHInterfaces {
		// Get the class name to create the regexp for searching in the implementation file
		className := string(matchedInterface[interfaceRegexpNameIndex])

		// Check the tags to know if the class needs to be processed
		tag := strings.TrimSpace(string(matchedInterface[interfaceRegexpTagIndex]))
		if (p.mustExcludeClassWithTag(tag)) {
			fmt.Fprintf(os.Stderr, "Ignoring class %v. Tag {%v} is either equal to parser.ExcludeTag (%v) or it isn't equal to parser.IncludeTag (%v)\n", className, tag, p.ExcludeTag, p.IncludeTag)
			continue
		}

		// Get the whole @interface bytes from header file
		interfaceHBytes := matchedInterface[0]

		// Get the whole @interface bytes from the implementation file
		classInterfaceRegexp := regexp.MustCompile(`(?ms:^\s?@interface\s+` + className + `\s+.*?@end)`)
		interfaceMBytes := classInterfaceRegexp.Find(implFileBytes)

		// Get the whole @implementation from the implementation file
		implRegexp := regexp.MustCompile(`(?ms:^\s?@implementation\s+` + className + `\s+.*?@end)`)
		matchedImpl := implRegexp.FindIndex(implFileBytes)
		implBytes := implFileBytes[matchedImpl[0]:matchedImpl[1]]

		classesInfo = append(classesInfo, NewObjCClass(className, interfaceHBytes, interfaceMBytes, implBytes, matchedImpl[0]))
	}
	return classesInfo
}

func (p Parser) mustExcludeClassWithTag(tag string) bool {
	if tag == p.ExcludeTag {
		return true
	}
	return p.IncludeTag != "" && tag != p.IncludeTag
}

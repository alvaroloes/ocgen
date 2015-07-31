package generator

import (
	"bytes"
	"fmt"
	"log"

	"github.com/alvaroloes/ocgen/parser"
)

func GenerateMethods(classes []parser.ObjCClass) {
	for _, class := range classes {
		codingInitMethod, err := getNSCodingInit(&class)
		if err != nil {
			log.Printf("Error when generating NSCoding.initWithCoder method: %v", err)
		}
		codingEncodeMethod, err := getNSCodingEncode(&class)
		if err != nil {
			log.Printf("Error when generating NSCoding.encodeWithCoder method: %v", err)
		}
		copyingMethod, err := getNSCopying(&class)
		if err != nil {
			log.Printf("Error when generating NSCopying.copyWithZone method: %v", err)
		}
		fmt.Println("\n---> Class: " + class.Name)
		fmt.Println("* NSCoding.init:", codingInitMethod)
		fmt.Println("* NSCoding.encode:", codingEncodeMethod)
		fmt.Println("* NSCopying.copy:", copyingMethod)
	}
}

func getNSCopying(class *parser.ObjCClass) (string, error) {
	var res bytes.Buffer
	err := NSCopyingTpl.Execute(&res, class)
	return res.String(), err
}

func getNSCodingInit(class *parser.ObjCClass) (string, error) {
	var res bytes.Buffer
	err := NSCodingInitTpl.Execute(&res, class)
	return res.String(), err
}

func getNSCodingEncode(class *parser.ObjCClass) (string, error) {
	var res bytes.Buffer
	err := NSCodingEncodeTpl.Execute(&res, class)
	return res.String(), err
}

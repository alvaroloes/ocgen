package generator

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/alvaroloes/ocgen/parser"
)

func GenerateMethods(classes []parser.ObjCClass) {
	for _, class := range classes {
		//TODO: Make a backup of the file and override original file
		// implSrcFile, err := ioutil.ReadFile(class.ImplFileName)
		// if err != nil {
		// 	log.Printf("Class: %v. Unable to open implementation file: %v", class.Name, class.ImplFileName)
		// 	continue
		// }

		// implDstFile, err := os.Create(class.ImplFileName + ".ocgen")
		// if err != nil {
		// 	log.Printf("Class: %v. Unable to create implementation file: %v", class.Name, class.ImplFileName)
		// 	continue
		// }

		//TODO: Get the methods sorted by appearance
		//TODO: Write all before the method, write method, write all after it and before the following method

		codingInitMethod, err := getNSCodingInit(&class)
		if err == nil {
			fmt.Println("* NSCoding.init:", string(codingInitMethod))
			//writeMethod(codingInitMethod, class.NSCodingInfo.InitWithCoder, implFile)
		} else {
			log.Printf("Class: %v. Error when generating NSCoding.initWithCoder method: %v\n", class.Name, err)
		}

		codingEncodeMethod, err := getNSCodingEncode(&class)
		if err == nil {
			fmt.Println("* NSCoding.encode:", string(codingEncodeMethod))
		} else {
			log.Printf("Class: %v. Error when generating NSCoding.encodeWithCoder method: %v\n", class.Name, err)
		}

		copyingMethod, err := getNSCopying(&class)
		if err == nil {
			fmt.Println("* NSCopying.copy:", string(copyingMethod))
		} else {
			log.Printf("Class: %v. Error when generating NSCopying.copyWithZone method: %v\n", class.Name, err)
		}
	}
}

func getNSCopying(class *parser.ObjCClass) ([]byte, error) {
	var res bytes.Buffer
	err := NSCopyingTpl.Execute(&res, class)
	return res.Bytes(), err
}

func getNSCodingInit(class *parser.ObjCClass) ([]byte, error) {
	var res bytes.Buffer
	err := NSCodingInitTpl.Execute(&res, class)
	return res.Bytes(), err
}

func getNSCodingEncode(class *parser.ObjCClass) ([]byte, error) {
	var res bytes.Buffer
	err := NSCodingEncodeTpl.Execute(&res, class)
	return res.Bytes(), err
}

func writeMethod(methodText []byte, methodInfo parser.MethodInfo, writer io.Writer) {

}

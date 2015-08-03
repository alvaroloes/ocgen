package generator

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alvaroloes/ocgen/parser"
)

var BackupFileExt = ".backup"

func GenerateMethods(classFile *parser.ObjCClassFile) error {
	if err := createBackup(classFile.MName); err != nil {
		return err
	}

	// fileBytes, err := ioutil.ReadFile(classFile.MName)
	// if err != nil {
	// 	log.Printf("Unable to open implementation file: %v", classFile.MName)
	// 	return err
	// }

	// TODO Open the MFile for writing (os.Create)

	for _, class := range classFile.Classes {

		//TODO: insert the methods bytes in the fileBytes slice in the corresponding location
		//TODO: Write the fileBytes into the MFile

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
	return nil
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

func createBackup(fileName string) (err error) {
	backupFileName := fileName + BackupFileExt
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	backupFile, err := os.Create(backupFileName)
	if err != nil {
		return
	}
	defer func() {
		err = backupFile.Close()
	}()

	_, err = io.Copy(backupFile, file)
	return
}

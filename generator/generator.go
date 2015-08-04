package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/alvaroloes/ocgen/parser"
)

var BackupFileExt = ".backup"

func GenerateMethods(classFile *parser.ObjCClassFile) error {
	if err := createBackup(classFile.MName); err != nil {
		log.Printf("Unable to create a backup file. Error: %v", err)
		return err
	}

	fileBytes, err := ioutil.ReadFile(classFile.MName)
	if err != nil {
		log.Printf("Unable to open implementation file: %v", classFile.MName)
		return err
	}

	// Classes and their method infos are sorted by appearance. We need to traverse them backwards
	// to keep the fields PosStart and PosEnd of the MethodInfo's in sync with the fileBytes when
	// inserting the new methods
	for i := len(classFile.Classes) - 1; i >= 0; i-- {
		class := classFile.Classes[i]

		methodsInfo := getMethodsInfoSortedBackwards(&class)

		for _, methodInfo := range methodsInfo {
			var methodBytes []byte
			// TODO: Remove the need of this switch  by creating a struct with the method info and the
			// methos bytes together
			switch methodInfo {
			case &class.NSCodingInfo.InitWithCoder:
				fmt.Println("Hey: InitWithCoder")
				methodBytes, err = getNSCodingInit(&class)
				if err != nil {
					log.Printf("Class: %v. Error when generating NSCoding.initWithCoder method: %v\n", class.Name, err)
				}
			case &class.NSCodingInfo.EncodeWithCoder:
				fmt.Println("Hey: EncodeWithCoder")
				methodBytes, err = getNSCodingEncode(&class)
				if err != nil {
					log.Printf("Class: %v. Error when generating NSCoding.encodeWithCoder method: %v\n", class.Name, err)
				}
			case &class.NSCopyingInfo.CopyWithZone:
				fmt.Println("Hey: CopyWithZone")
				methodBytes, err = getNSCopying(&class)
				if err != nil {
					log.Printf("Class: %v. Error when generating NSCopying.copyWithZone method: %v\n", class.Name, err)
				}
			}

			fileBytes = insertMethod(fileBytes, methodBytes, *methodInfo)
		}
	}

	// TODO: write the file
	fmt.Println(string(fileBytes))
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

func getMethodsInfoSortedBackwards(class *parser.ObjCClass) []*parser.MethodInfo {
	methods := []*parser.MethodInfo{
		&class.NSCodingInfo.InitWithCoder,
		&class.NSCodingInfo.EncodeWithCoder,
		&class.NSCopyingInfo.CopyWithZone,
	}
	sort.Sort(sort.Reverse(MethodsInfoByPosStart(methods)))
	return methods
}

func insertMethod(fileBytes, newMethod []byte, oldMethodInfo parser.MethodInfo) []byte {
	newMethodAndNextBytes := append(newMethod, fileBytes[oldMethodInfo.PosEnd:]...)
	return append(fileBytes[:oldMethodInfo.PosStart], newMethodAndNextBytes...)
}

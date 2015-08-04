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

		methodsInfo := getMethodsInfoSortedBackwards(class)

		fmt.Println(methodsInfo)
		// fmt.Println(string(fileBytes[class.NSCodingInfo.InitWithCoder.PosStart:class.NSCodingInfo.InitWithCoder.PosEnd]))
		// fmt.Println(string(fileBytes[class.NSCodingInfo.EncodeWithCoder.PosStart:class.NSCodingInfo.EncodeWithCoder.PosEnd]))
		// fmt.Println(string(fileBytes[class.NSCopyingInfo.CopyWithZone.PosStart:class.NSCopyingInfo.CopyWithZone.PosEnd]))

		//TODO: insert the methods bytes in the fileBytes slice in the corresponding location
		//TODO: Write the fileBytes into the MFile

		codingInitMethod, err := getNSCodingInit(&class)
		if err == nil {
			fileBytes = insertMethod(fileBytes, codingInitMethod, class.NSCodingInfo.InitWithCoder)
			fmt.Println(string(fileBytes))
			//writeMethod(codingInitMethod, class.NSCodingInfo.InitWithCoder, implFile)
		} else {
			log.Printf("Class: %v. Error when generating NSCoding.initWithCoder method: %v\n", class.Name, err)
		}

		// codingEncodeMethod, err := getNSCodingEncode(&class)
		// if err == nil {
		// 	fmt.Println("* NSCoding.encode:", string(codingEncodeMethod))
		// } else {
		// 	log.Printf("Class: %v. Error when generating NSCoding.encodeWithCoder method: %v\n", class.Name, err)
		// }

		// copyingMethod, err := getNSCopying(&class)
		// if err == nil {
		// 	fmt.Println("* NSCopying.copy:", string(copyingMethod))
		// } else {
		// 	log.Printf("Class: %v. Error when generating NSCopying.copyWithZone method: %v\n", class.Name, err)
		// }
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

func getMethodsInfoSortedBackwards(class parser.ObjCClass) []*parser.MethodInfo {
	methods := []*parser.MethodInfo{
		&class.NSCodingInfo.InitWithCoder,
		&class.NSCodingInfo.EncodeWithCoder,
		&class.NSCopyingInfo.CopyWithZone,
	}
	sort.Sort(MethodsInfoByPosStart(methods))
	return methods
}

func insertMethod(fileBytes, newMethod []byte, oldMethodInfo parser.MethodInfo) []byte {
	fmt.Println(oldMethodInfo)
	newMethodAndNextBytes := append(newMethod, fileBytes[oldMethodInfo.PosEnd:]...)
	return append(fileBytes[:oldMethodInfo.PosStart], newMethodAndNextBytes...)
}

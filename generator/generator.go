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

type templateGenerator func(*parser.ObjCClass) ([]byte, error)

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

		methodGenerators, sortedMethodsInfo := getMethodsGenerators(&class)

		for _, methodInfo := range sortedMethodsInfo {
			methodBytes, _ := methodGenerators[methodInfo](&class)
			// TODO: Remove the need of this switch  by creating a struct with the method info and the
			// methos bytes together
			// switch methodInfo {
			// case &class.NSCodingInfo.InitWithCoder:
			// 	fmt.Println("Hey: InitWithCoder")
			// 	methodBytes, err = getNSCodingInit(&class)
			// 	if err != nil {
			// 		log.Printf("Class: %v. Error when generating NSCoding.initWithCoder method: %v\n", class.Name, err)
			// 	}
			// case &class.NSCodingInfo.EncodeWithCoder:
			// 	fmt.Println("Hey: EncodeWithCoder")
			// 	methodBytes, err = getNSCodingEncode(&class)
			// 	if err != nil {
			// 		log.Printf("Class: %v. Error when generating NSCoding.encodeWithCoder method: %v\n", class.Name, err)
			// 	}
			// case &class.NSCopyingInfo.CopyWithZone:
			// 	fmt.Println("Hey: CopyWithZone")
			// 	methodBytes, err = getNSCopying(&class)
			// 	if err != nil {
			// 		log.Printf("Class: %v. Error when generating NSCopying.copyWithZone method: %v\n", class.Name, err)
			// 	}
			// }

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

// Returns a map whose keys are pointers to all the structs "MethodInfo" present in "class" and the values are
// the "templateGenerator" for each MethodInfo.
// The second return value contains a slice with the map keys sorted backwards as defined by "MethodsInfoByPosStart"
func getMethodsGenerators(class *parser.ObjCClass) (map[*parser.MethodInfo]templateGenerator, []*parser.MethodInfo) {
	generatorByMethod := map[*parser.MethodInfo]templateGenerator{
		&class.NSCodingInfo.InitWithCoder:   getNSCodingInit,
		&class.NSCodingInfo.EncodeWithCoder: getNSCodingEncode,
		&class.NSCopyingInfo.CopyWithZone:   getNSCopying,
	}

	methods := make([]*parser.MethodInfo, 0, len(generatorByMethod))
	for method := range generatorByMethod {
		methods = append(methods, method)
	}

	sort.Sort(sort.Reverse(MethodsInfoByPosStart(methods)))
	return generatorByMethod, methods
}

func insertMethod(fileBytes, newMethod []byte, oldMethodInfo parser.MethodInfo) []byte {
	newMethodAndNextBytes := append(newMethod, fileBytes[oldMethodInfo.PosEnd:]...)
	return append(fileBytes[:oldMethodInfo.PosStart], newMethodAndNextBytes...)
}

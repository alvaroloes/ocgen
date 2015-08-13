package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/alvaroloes/ocgen/parser"
)

func GenerateMethods(classFile *parser.ObjCClassFile, backupDir string) error {
	if backupDir != "" {
		if err := createBackup(classFile.MName, backupDir); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create a backup file. Error: %v\n", err)
			return err
		}
	}

	fileBytes, err := ioutil.ReadFile(classFile.MName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open implementation file: %v", classFile.MName)
		return err
	}

	// Classes and their method infos are sorted by appearance. We need to traverse them backwards
	// to keep the fields PosStart and PosEnd of the MethodInfo's in sync with the fileBytes when
	// inserting the new methods
	for i := len(classFile.Classes) - 1; i >= 0; i-- {
		class := classFile.Classes[i]

		methodGenerators, sortedMethodsInfo := getMethodsGenerators(&class)

		for _, methodInfo := range sortedMethodsInfo {
			methodBytes, err := methodGenerators[methodInfo](&class)

			if err == nil {
				fileBytes = insertMethod(fileBytes, methodBytes, *methodInfo)
			} else {
				fmt.Fprintf(os.Stderr, `Class: %v. Error when generating "%v" method: %v\n`, class.Name, methodInfo.Name, err)
			}
		}
	}

	// Write the file with the new content
	// (Permissions will be ignored as the file already exists)
	return ioutil.WriteFile(classFile.MName, fileBytes, 0664)
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

func createBackup(fileName, backupDir string) (err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	backupFileName := path.Join(backupDir, fileName)

	err = os.MkdirAll(path.Dir(backupFileName), 0775)
	if err != nil {
		return
	}

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

type templateGenerator func(*parser.ObjCClass) ([]byte, error)

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

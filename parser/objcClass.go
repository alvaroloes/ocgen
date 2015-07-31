package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	endRegexp      = regexp.MustCompile(`\s?@end`)
	propertyRegexp = regexp.MustCompile(`@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)
)
var (
	codingInitMethodRegexp   = regexp.MustCompile(`\s?-.*initWithCoder:`)
	codingEncodeMedhotRegexp = regexp.MustCompile(`\s?-.*encodeWithCoder:`)
	copyingMethodRexexp      = regexp.MustCompile(`\s?-.*copyWithZone:`)
)

const (
	propertyRegexpAttrIndex = iota + 1
	propertyRegexpClassIndex
	propertyRegexpPointerIndex
	propertyRegexpNameIndex
)

type ObjCClass struct {
	// These fields are extracted from the header file
	Name              string
	ImplFileName      string
	Properties        []Property // These are also extracted form the implementation file
	ConformsNSCoding  bool
	ConformsNSCopying bool

	// These fields are extracted from the implementation file
	NSCodingInfo struct {
		InitWithCoder   MethodInfo
		EncodeWithCoder MethodInfo
	}
	NSCopyingInfo struct {
		CopyWithZone MethodInfo
	}
}

type MethodInfo struct {
	PosStart, PosEnd int
}

type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

func NewObjCClass(className string, hInterfaceBytes, mInterfaceBytes, implBytes []byte, implBytesOffset int, implFileName string) ObjCClass {
	propertiesFromH := extractProperties(hInterfaceBytes)
	propertiesFromM := extractProperties(mInterfaceBytes)

	class := ObjCClass{
		Name:         className,
		ImplFileName: implFileName,
		Properties:   mergeProperties(propertiesFromH, propertiesFromM),
		//TODO: Detect if the class conforms the protocols taking into account the parent protocols too
		ConformsNSCoding:  true,
		ConformsNSCopying: true,
	}

	class.NSCodingInfo.InitWithCoder = extractMethodInfo(className, codingInitMethodRegexp, implBytes, implBytesOffset)
	class.NSCodingInfo.EncodeWithCoder = extractMethodInfo(className, codingEncodeMedhotRegexp, implBytes, implBytesOffset)
	class.NSCopyingInfo.CopyWithZone = extractMethodInfo(className, copyingMethodRexexp, implBytes, implBytesOffset)

	return class
}

func extractProperties(interfaceBytes []byte) []Property {
	var properties []Property

	matchedProperties := propertyRegexp.FindAllSubmatch(interfaceBytes, -1)
	for _, matchedProperty := range matchedProperties {
		//Split the attributes string and trim each of them
		attributes := strings.Split(string(matchedProperty[propertyRegexpAttrIndex]), ",")
		for i, attr := range attributes {
			attributes[i] = strings.TrimSpace(attr)
		}
		class := string(matchedProperty[propertyRegexpClassIndex])
		pointer := string(matchedProperty[propertyRegexpPointerIndex])
		name := string(matchedProperty[propertyRegexpNameIndex])

		// Add this property to the info
		properties = append(properties, Property{
			Name:       name,
			Class:      class,
			Attributes: attributes,
			IsPointer:  pointer != "",
		})
	}

	return properties
}

func mergeProperties(propertiesFromH, propertiesFromM []Property) []Property {
	// TODO:
	// - Join both slices
	// - Remove the readonly properties that has a corresponding readwrite version in M
	return propertiesFromH
}

func extractMethodInfo(className string, methodSignatureRegexp *regexp.Regexp, implBytes []byte, implBytesOffset int) (methodInfo MethodInfo) {
	matchedMethod := methodSignatureRegexp.FindIndex(implBytes)

	if matchedMethod == nil {
		log.Printf(`Method not found (Regexp: %v) in class "%v"\n`, methodSignatureRegexp, className)
		// There is no previous method, the position for the new one will be just before @end
		matchedEnd := endRegexp.FindIndex(implBytes)
		methodInfo.PosStart, methodInfo.PosEnd = matchedEnd[0], matchedEnd[0]
	} else {
		methodInfo.PosStart = matchedMethod[0]
		bodyStart := matchedMethod[1]
		relativeBodyEnd := relativeEndOfMethodBody(implBytes[bodyStart:])
		methodInfo.PosEnd = bodyStart + relativeBodyEnd
	}

	fmt.Println(methodInfo)
	fmt.Println(string(implBytes[methodInfo.PosStart:methodInfo.PosEnd]))

	methodInfo.PosStart += implBytesOffset
	methodInfo.PosEnd += implBytesOffset

	return
}

func relativeEndOfMethodBody(bytes []byte) int {
	numBrackets := 0
	insideBody := false

	for i, b := range bytes {
		if b == '{' {
			insideBody = true
			numBrackets++
		} else if b == '}' {
			numBrackets--
		}

		if insideBody && numBrackets <= 0 {
			return i + 1
		}
	}
	return -1
}

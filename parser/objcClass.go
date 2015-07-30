package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	propertyRegexp         = regexp.MustCompile(`@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)
	codingInitMethodRegexp = regexp.MustCompile(`\s?-.*initWithCoder:`)
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
	PosStart, PosEnd int64
}

type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

func NewObjCClass(className string, hInterfaceBytes, mInterfaceBytes, implBytes []byte, implFileName string) ObjCClass {
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

	extractProtocolMethodsInfo(&class, implBytes)

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

func extractProtocolMethodsInfo(class *ObjCClass, implFileBytes []byte) {
	matchedCodingInitMethod := codingInitMethodRegexp.FindIndex(implFileBytes)

	if matchedCodingInitMethod == nil {
		// TODO: set something to indicate that the methods should be added to the end
		return
	}

	bodyStart := matchedCodingInitMethod[1]
	relativeBodyEnd := relativeEndOfMethodBody(implFileBytes[bodyStart:])
	methodStart := matchedCodingInitMethod[0]
	methodEnd := bodyStart + relativeBodyEnd

	//TODO: generalize this to extract info for all the methods

	fmt.Println(string(implFileBytes[methodStart:methodEnd]))
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

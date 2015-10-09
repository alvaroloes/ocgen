package parser

import (
	"log"
	"regexp"
	"strings"
)

var (
	endRegexp      = regexp.MustCompile(`\s?@end`)
	propertyRegexp = regexp.MustCompile(`@property\s*(?:\((.*)\))?\s*([^\s\*]*)\s*(\*)?\s*(.*);`)
	parentAndProtocolsRegexp = regexp.MustCompile(`@interface[^:<]*(?::\s*([^<\s]*))?(?:\s*<([^>]*)>)?`)
)

var (
	codingInitMethodName   = "initWithCoder:"
	codingInitMethodRegexp = regexp.MustCompile(`\s?-.*` + codingInitMethodName)

	codingEncodeMethodName   = "encodeWithCoder:"
	codingEncodeMedhotRegexp = regexp.MustCompile(`\s?-.*` + codingEncodeMethodName)

	copyingMethodName   = "copyWithZone:"
	copyingMethodRexexp = regexp.MustCompile(`\s?-.*` + copyingMethodName)
)

const (
	propertyRegexpAttrIndex = iota + 1
	propertyRegexpClassIndex
	propertyRegexpPointerIndex
	propertyRegexpNameIndex
)

const (
	parentAndProtocolsRegexpParentIndex = iota + 1
	parentAndProtocolsRegexpProtocolsIndex
)

const (
	nsObjectToken = "NSObject"
)

type ObjCClassFile struct {
	HName, MName string
	Classes      []ObjCClass
}

type ObjCClass struct {
	// These fields are extracted from the header file
	Name              string
	Parent            string
	Protocols         []string
	Properties        []Property // These are also extracted form the implementation file

	// These fields are extracted from the implementation file
	NSCodingInfo struct {
		InitWithCoder   MethodInfo
		EncodeWithCoder MethodInfo
	}
	NSCopyingInfo struct {
		CopyWithZone MethodInfo
	}
}

func (oc *ObjCClass) IsDirectChildOfNSObject() bool {
	return oc.Parent == nsObjectToken
}

func (oc *ObjCClass) ConformsProtocol(protocol string) bool {
	for _,classProto := range oc.Protocols {
		if classProto == protocol {
			return true
		}
	}
	return false;
}

func (oc *ObjCClass) ConformsAnyProtocol(protocols... string) bool {
	for _,protocol := range protocols {
		if oc.ConformsProtocol(protocol){
			return true
		}
	}
	return false;
}

type MethodInfo struct {
	Name             string
	PosStart, PosEnd int
}

func NewObjCClass(className string, hInterfaceBytes, mInterfaceBytes, implBytes []byte, implBytesOffset int) ObjCClass {
	propertiesFromH := extractProperties(hInterfaceBytes)
	propertiesFromM := extractProperties(mInterfaceBytes)
	parent, protocols := extractParentAndProtocols(hInterfaceBytes)

	class := ObjCClass{
		Name:       className,
		Parent:     parent,
		Protocols:  protocols,
		Properties: mergeProperties(propertiesFromH, propertiesFromM),
	}

	class.NSCodingInfo.InitWithCoder = extractMethodInfo(codingInitMethodName, className, codingInitMethodRegexp, implBytes, implBytesOffset)
	class.NSCodingInfo.EncodeWithCoder = extractMethodInfo(codingEncodeMethodName, className, codingEncodeMedhotRegexp, implBytes, implBytesOffset)
	class.NSCopyingInfo.CopyWithZone = extractMethodInfo(copyingMethodName, className, copyingMethodRexexp, implBytes, implBytesOffset)
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

func extractParentAndProtocols(interfaceBytes []byte) (parent string, protocols []string) {
	match := parentAndProtocolsRegexp.FindSubmatch(interfaceBytes)
	parent = string(match[parentAndProtocolsRegexpParentIndex])
	protocols = strings.Split(string(match[parentAndProtocolsRegexpProtocolsIndex]), ",")
	for i, proto := range protocols {
		protocols[i] = strings.TrimSpace(proto)
	}
	return
}

func mergeProperties(propertiesFromH, propertiesFromM []Property) []Property {
	// Join the two property slices avoiding duplicates (for example a property in .h with
	// a "readonly" attribute and the same property in .m with "readwrite")
	// Properties in .m have preference.

	hPropsMap := make(map[string]int, len(propertiesFromH))
	for i,prop := range propertiesFromH {
		hPropsMap[prop.Name] = i // Save the index of the property in case we find a duplicate later
	}

	for _,prop := range propertiesFromM {
		if index,exists := hPropsMap[prop.Name]; exists {
			propertiesFromH[index] = prop;
		} else {
			propertiesFromH = append(propertiesFromH,prop)
		}
	}
	return propertiesFromH
}

func extractMethodInfo(methodName, className string, methodSignatureRegexp *regexp.Regexp, implBytes []byte, implBytesOffset int) (methodInfo MethodInfo) {
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

	methodInfo.Name = methodName
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

package parser
import (
	"strings"
)

const (
	weakToken = "weak"
	readonlyToken = "readonly"
	idTypeToken = "id"

	propertyAccessor = "."
	ivarAccessor = "->_"

	defaultCoderType = "Integer"
	objectCoderType = "Object"
)

var coderTypePerClass = map[string]string {
	"BOOL": "Bool",
	"bool": "Bool",
	"Boolean": "Bool",

	"float": "Float",

	"CGFloat": "Double",
	"double": "Double",

	"CGPoint": "CGPoint",
	"CGRect": "CGRect",
	"CGSize": "CGSize",

	// Default is Integer
}


type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

func (p *Property) IsObject() bool {
	return p.IsPointer || strings.HasPrefix(p.Class, idTypeToken)
}

func (p *Property) IsWeak() bool {
	for _,attrToken := range p.Attributes {
		if (attrToken == weakToken) {
			return true
		}
	}
	return false
}

func (p *Property) IsReadonly() bool {
	for _,attrToken := range p.Attributes {
		if (attrToken == readonlyToken) {
			return true
		}
	}
	return false
}

// Returns either the property or the ivar accessor
func (p *Property) WriteAccessor() string {
	if (p.IsReadonly()) {
		return ivarAccessor
	}
	return propertyAccessor
}

func (p *Property) CoderType() string {
	if p.IsObject() {
		return objectCoderType
	}
	coderType, exists := coderTypePerClass[p.Class]
	if !exists {
		coderType = defaultCoderType
	}
	return coderType
}
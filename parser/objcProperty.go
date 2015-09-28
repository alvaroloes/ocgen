package parser
import (
	"strings"
)

const (
	weakToken = "weak"
	readonlyToken = "readonly"
	idTypeToken = "id"

	defaultAccessor = "."
	readonlyAccessor = "->"
)


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

// Returns the accessor for this property (-> or .)
func (p *Property) Accessor() string {
	if (p.IsReadonly()) {
		return readonlyAccessor
	}
	return defaultAccessor
}
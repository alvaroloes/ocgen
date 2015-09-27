package parser

const (
	weakToken = "weak"
)


type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

func (p *Property) IsWeak() bool {
	for _,attrToken := range p.Attributes {
		if (attrToken == weakToken) {
			return true
		}
	}
	return false
}
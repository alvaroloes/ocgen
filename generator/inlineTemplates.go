package generator

import "text/template"

var NSCopyingText = `{{define "pointerTpl"}}{{if eq .Class "NSArray"}}copy.{{.Name}} = [[NSArray alloc] initWithArray:self.{{.Name}} copyItems:YES];{{else}}copy.{{.Name}} = [self.{{.Name}} copyWithZone:zone];{{end}}{{end}}
// NSCopying protocol method
- (instancetype)copyWithZone:(NSZone *)zone
{
    typeof(self) copy = [[[self class] alloc] init];
    if (copy != nil)
    { {{range .Properties}}
        {{if .IsPointer}}{{template "pointerTpl" .}}{{else}}copy.{{.Name}} = self.{{.Name}};{{end}}{{end}}
	}
	return copy;
}`
var NSCopyingTpl = template.Must(template.New("NSCopying").Parse(NSCopyingText))

var NSCodingInitText = `
{{define "pointerDecodeTpl"}}_{{.Name}} = [decoder decodeObjectForKey:@"{{.Name}}"];{{end}}
// NSCoding protocol methods
- (instancetype)initWithCoder:(NSCoder *)decoder
{
	if(self = [super init])
    { {{range .Properties}}
        {{if .IsPointer}}{{template "pointerDecodeTpl" .}}{{else}}_{{.Name}} = [decoder decodeIntegerForKey:@"{{.Name}}"];{{end}}{{end}}
    }

	return self;
}
{{define "pointerEncodeTpl"}}[coder encodeObject:_{{.Name}} forKey:@"{{.Name}}"];{{end}}
- (void)encodeWithCoder:(NSCoder *)coder
{ {{range .Properties}}
    {{if .IsPointer}}{{template "pointerEncodeTpl" .}}{{else}}[coder encodeInteger:_{{.Name}} forKey:@"{{.Name}}"];{{end}}{{end}}
}`
var NSCodingInitTpl = template.Must(template.New("NSCodingInit").Parse(NSCodingInitText))
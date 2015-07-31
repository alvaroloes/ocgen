package generator

import "text/template"

var NSCopyingText = `{{define "pointerTpl"}}{{if eq .Class "NSArray"}}copy.{{.Name}} = [[NSArray alloc] initWithArray:self.{{.Name}} copyItems:YES];{{else}}copy.{{.Name}} = [self.{{.Name}} copyWithZone:zone];{{end}}{{end}}
- (instancetype)copyWithZone:(NSZone *)zone
{
    // OCGEN: Autogenerated method. Do not touch
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
- (instancetype)initWithCoder:(NSCoder *)decoder
{
    // OCGEN: Autogenerated method. Do not touch
    if(self = [super init])
    { {{range .Properties}}
        {{if .IsPointer}}{{template "pointerDecodeTpl" .}}{{else}}_{{.Name}} = [decoder decodeIntegerForKey:@"{{.Name}}"];{{end}}{{end}}
    }

    return self;
}`
var NSCodingInitTpl = template.Must(template.New("NSCodingInit").Parse(NSCodingInitText))

var NSCodingEncodeText = `
{{define "pointerEncodeTpl"}}[coder encodeObject:_{{.Name}} forKey:@"{{.Name}}"];{{end}}
- (void)encodeWithCoder:(NSCoder *)coder
{
    // OCGEN: Autogenerated method. Do not touch {{range .Properties}}
    {{if .IsPointer}}{{template "pointerEncodeTpl" .}}{{else}}[coder encodeInteger:_{{.Name}} forKey:@"{{.Name}}"];{{end}}{{end}}
}`
var NSCodingEncodeTpl = template.Must(template.New("NSCodingEncode").Parse(NSCodingEncodeText))

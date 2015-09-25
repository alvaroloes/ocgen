package generator

import "text/template"

var NSCopyingText = `
- (instancetype)copyWithZone:(NSZone *)zone
{
    // OCGEN: Autogenerated method. Do not touch
    typeof(self) copy = nil;
    if([super respondsToSelector:@selector(copyWithZone:)])
    {
        copy = [super copyWithZone:zone];
    }
    else
    {
        copy = [[[self class] allocWithZone:zone] init];
    }
    if (copy != nil)
    { {{range .Properties}}
        {{if .IsPointer -}}
            {{if eq .Class "NSArray" -}}
                copy.{{.Name}} = [[NSArray alloc] initWithArray:self.{{.Name}} copyItems:YES];
            {{- else -}}
                copy.{{.Name}} = [self.{{.Name}} copyWithZone:zone];
            {{- end}}
        {{- else -}}
            copy.{{.Name}} = self.{{.Name}};
        {{- end}}
    {{- end}}
    }
    return copy;
}`
var NSCopyingTpl = template.Must(template.New("NSCopying").Parse(NSCopyingText))

var NSCodingInitText = `
- (instancetype)initWithCoder:(NSCoder *)decoder
{
    // OCGEN: Autogenerated method. Do not touch
    if([super respondsToSelector:@selector(initWithCoder:)])
    {
        self = [super initWithCoder:decoder];
    }
    else
    {
        self = [super init];
    }
    if(self != nil)
    { {{range .Properties}}
        {{if .IsPointer -}}
            _{{.Name}} = [decoder decodeObjectForKey:@"{{.Name}}"];
        {{- else -}}
            _{{.Name}} = [decoder decodeIntegerForKey:@"{{.Name}}"];
        {{- end}}
    {{- end}}
    }
    return self;
}`
var NSCodingInitTpl = template.Must(template.New("NSCodingInit").Parse(NSCodingInitText))

var NSCodingEncodeText = `
- (void)encodeWithCoder:(NSCoder *)coder
{
    // OCGEN: Autogenerated method. Do not touch
    if([super respondsToSelector:@selector(encodeWithCoder:)])
    {
        [super encodeWithCoder:coder];
    }{{range .Properties}}
    {{if .IsPointer -}}
        [coder encodeObject:_{{.Name}} forKey:@"{{.Name}}"];
    {{- else -}}
        [coder encodeInteger:_{{.Name}} forKey:@"{{.Name}}"];
    {{- end}}
{{- end}}
}`
var NSCodingEncodeTpl = template.Must(template.New("NSCodingEncode").Parse(NSCodingEncodeText))

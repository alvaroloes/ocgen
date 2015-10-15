package generator

import "text/template"

var NSCopyingText = `
- (instancetype)copyWithZone:(NSZone *)zone
{
    // OCGEN: Autogenerated method
    {{if not .IsDirectChildOfNSObject -}}
    {{.Name}} *copy = nil;
    if([super respondsToSelector:@selector(copyWithZone:)])
    {
        copy = [super copyWithZone:zone];
    }
    else
    {
        copy = [[[self class] allocWithZone:zone] init];
    }
    {{- else -}}
        {{.Name}} *copy = [[[self class] allocWithZone:zone] init];
    {{- end}}
    if (copy != nil)
    { {{range .Properties}}
        copy{{.WriteAccessor}}{{.Name}} = {{if .IsObject -}}
                                          {{if eq .Class "NSArray" -}}
                                              [[NSArray alloc] initWithArray:self.{{.Name}} copyItems:YES];
                                          {{- else if .IsWeak -}}
                                              self.{{.Name}};
                                          {{- else -}}
                                              [self.{{.Name}} copyWithZone:zone];
                                          {{- end}}
                                      {{- else -}}
                                          self.{{.Name}};
                                      {{- end}}
    {{- end}}
    }
    return copy;
}`
var NSCopyingTpl = template.Must(template.New("NSCopying").Parse(NSCopyingText))

var NSCodingInitText = `
- (instancetype)initWithCoder:(NSCoder *)decoder
{
    // OCGEN: Autogenerated method
    {{if not .IsDirectChildOfNSObject -}}
    if([super respondsToSelector:@selector(initWithCoder:)])
    {
        self = [super initWithCoder:decoder];
    }
    else
    {
        self = [super init];
    }
    {{- else -}}
    self = [super init];
    {{- end}}
    if(self != nil)
    { {{range .Properties}}
        _{{.Name}} = [decoder decode{{.CoderType}}ForKey:@"{{.Name}}"];
    {{- end}}
    }
    return self;
}`
var NSCodingInitTpl = template.Must(template.New("NSCodingInit").Parse(NSCodingInitText))

var NSCodingEncodeText = `
- (void)encodeWithCoder:(NSCoder *)coder
{
    // OCGEN: Autogenerated method
    {{- if not .IsDirectChildOfNSObject}}
    if([super respondsToSelector:@selector(encodeWithCoder:)])
    {
        [super encodeWithCoder:coder];
    }
    {{- end}}{{range .Properties}}
    {{if .IsWeak -}}
        [coder encodeConditionalObject:self.{{.Name}} forKey:@"{{.Name}}"];
    {{- else -}}
        [coder encode{{.CoderType}}:self.{{.Name}} forKey:@"{{.Name}}"];
    {{- end}}
{{- end}}
}`
var NSCodingEncodeTpl = template.Must(template.New("NSCodingEncode").Parse(NSCodingEncodeText))

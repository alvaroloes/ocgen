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
        copy{{.Accessor}}{{.Name}} = {{if .IsObject -}}
                                          {{if eq .Class "NSArray" -}}
                                              [[NSArray alloc] initWithArray:self{{.Accessor}}{{.Name}} copyItems:YES];
                                          {{- else if .IsWeak -}}
                                              self{{.Accessor}}{{.Name}};
                                          {{- else -}}
                                              [self{{.Accessor}}{{.Name}} copyWithZone:zone];
                                          {{- end}}
                                      {{- else -}}
                                          self{{.Accessor}}{{.Name}};
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
        _{{.Name}} = {{if .IsObject -}}
                         [decoder decodeObjectForKey:@"{{.Name}}"];
                     {{- else -}}
                         [decoder decode{{.CoderType}}ForKey:@"{{.Name}}"];
                     {{- end}}
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
    {{if .IsObject -}}
        {{if .IsWeak -}}
            [coder encodeConditionalObject:_{{.Name}} forKey:@"{{.Name}}"];
        {{- else -}}
            [coder encodeObject:_{{.Name}} forKey:@"{{.Name}}"];
        {{- end}}
    {{- else -}}
        [coder encode{{.CoderType}}:_{{.Name}} forKey:@"{{.Name}}"];
    {{- end}}
{{- end}}
}`
var NSCodingEncodeTpl = template.Must(template.New("NSCodingEncode").Parse(NSCodingEncodeText))

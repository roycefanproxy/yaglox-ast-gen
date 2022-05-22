package main

const TemplateSource = `
{{$defs := .Definitions}}
{{$visitors := .VisitorTypes}}
package main

type Visitor[R any] interface {
{{range $def := $defs}}
    Visit{{$def.Name}}{{$.BaseName}}(expr *{{$def.Name}}) R
{{end}}
}

type {{.BaseName}} interface {
    {{range $visitor := $visitors}}
    Accept{{$visitor.Name}}(visitor Visitor[{{$visitor.Type}}]) string
    {{end}}
}

{{range $def := $defs}}
type {{$def.Name}} struct {
    {{range $j, $member := $def.Members}}{{$member}}{{"\n\t"}}{{end}}
}

{{range $visitor := $visitors}}
func (e *{{$def.Name}}) Accept{{$visitor.Name}}(visitor Visitor[{{$visitor.Type}}]) string {
    return visitor.Visit{{$def.Name}}{{$.BaseName}}(e)
}
{{end}}
{{end}}

`

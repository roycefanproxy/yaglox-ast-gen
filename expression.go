package main

const TemplateSource = `
{{$defs := .Definitions}}
{{$visitors := .VisitorTypes}}
package main

type Visitor[R any] interface {
    {{range $i, $def := $defs}}Visit{{$def.Name}}{{$.BaseName}}(expr *{{$def.Name}}) R{{if ne (len $defs) (add $i 1)}}{{"\n\t"}}{{end}}{{end}}
}

type {{.BaseName}} interface {
    {{range $j, $visitor := $visitors}}Accept{{$visitor.Name}}(visitor Visitor[{{$visitor.Type}}]) {{$visitor.Type}}{{if ne (len $visitors) (add $j 1)}}{{"\n\t"}}{{end}}{{end}}
}
{{range $def := $defs}}
type {{$def.Name}} struct {
    {{range $k, $member := $def.Members}}{{$member}}{{if ne (len $def.Members) (add $k 1)}}{{"\n\t"}}{{end}}{{end}}
}
{{range $visitor := $visitors}}
func (e *{{$def.Name}}) Accept{{$visitor.Name}}(visitor Visitor[{{$visitor.Type}}]) {{$visitor.Type}} {
    return visitor.Visit{{$def.Name}}{{$.BaseName}}(e)
}
{{end}}
{{end}}

`

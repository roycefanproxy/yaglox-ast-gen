package main

const TemplateSource = `
{{$defs := .Definitions}}
{{$visitors := .VisitorTypes}}
{{$baseName := .BaseName}}
package main

type {{$baseName}}Visitor[R any] interface {
    {{range $i, $def := $defs}}Visit{{$def.Name}}(expr *{{$def.Name}}) R{{if ne (len $defs) (add $i 1)}}{{"\n\t"}}{{end}}{{end}}
}

type {{$baseName}} interface {
    {{range $j, $visitor := $visitors}}Accept{{$visitor.Name}}(visitor {{$baseName}}Visitor[{{$visitor.Type}}]) {{$visitor.Type}}{{if ne (len $visitors) (add $j 1)}}{{"\n\t"}}{{end}}{{end}}
}
{{range $def := $defs}}
type {{$def.Name}} struct {
    {{range $k, $member := $def.Members}}{{$member}}{{if ne (len $def.Members) (add $k 1)}}{{"\n\t"}}{{end}}{{end}}
}
{{range $visitor := $visitors}}
func (e *{{$def.Name}}) Accept{{$visitor.Name}}(visitor {{$baseName}}Visitor[{{$visitor.Type}}]) {{$visitor.Type}} {
    return visitor.Visit{{$def.Name}}(e)
}
{{end}}
{{end}}

`

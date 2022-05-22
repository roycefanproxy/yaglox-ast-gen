package main

const TemplateSource = `
package main

type Visitor[R any] interface {
{{range $i, $def := .Definitions}}
    Visit{{$def.Name}}{{$.BaseName}}(expr *{{$def.Name}}) R
{{end}}
}

type {{.BaseName}} interface {
    Accept(visitor Visitor[string]) string
}

{{range $i, $def := .Definitions}}
type {{$def.Name}} struct {
    {{range $j, $member := $def.Members}}{{$member}}{{"\n\t"}}{{end}}
}

func (e *{{$def.Name}}) Accept(visitor Visitor[string]) string {
    return visitor.Visit{{$def.Name}}{{$.BaseName}}(e)
}
{{end}}

`

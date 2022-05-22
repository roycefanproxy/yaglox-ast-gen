package main

const TemplateSource = `
package main

type {{.BaseName}} interface {

}

{{range $i, $def := .Definitions}}
type {{$def.Name}} struct {
    {{range $j, $member := $def.Members}}{{$member}}{{"\n\t"}}{{end}}
}
{{end}}
`

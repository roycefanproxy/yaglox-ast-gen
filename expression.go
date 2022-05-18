package main

const TemplateSource = `
package main

import "github.com/roycefanproxy/yaglox"

type {{.BaseName}} struct {

}

{{range $i, $def := .Definitions}}
type {{$def.Name}} struct {
    {{range $j, $member := $def.Members}}{{$member}}{{"\n\t"}}{{end}}
}
{{end}}
`

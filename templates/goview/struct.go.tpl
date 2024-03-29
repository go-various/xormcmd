package {{.Models}}
{{$prefix := .Prefix}}
{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
)
{{end}}

{{range .Tables}}
type {{Mapper .Name}}View struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}} *{{Type $col}} {{Tag $table $col}}
{{end}}
}
{{end}}

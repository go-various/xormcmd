package {{.Models}}
{{$prefix := .Prefix}}
{{$ilen := len .Imports}}
{{if gt $ilen 0}}
import (
    "encoding/json"
	{{range .Imports}}"{{.}}"{{end}}
)
{{end}}

{{range .Tables}}
type {{Mapper .Name}} struct {
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}} *{{Type $col}} {{Tag $table $col}}
{{end}}
}
func (o *{{Mapper .Name}})TableName()string  {
	return "{{$prefix}}{{$table.Name}}"
}
func (o *{{Mapper .Name}})MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}
func (o *{{Mapper .Name}})UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}
func (o *{{Mapper .Name}}) ID() interface{} {
    {{$ilen := len $table.PrimaryKeys}}
    {{if eq $ilen 1}}
    {{range $index, $element := $table.PrimaryKeys}}
    return o.{{Mapper $element}}
    {{end}}
    {{else}}
    panic("multiplex primary keys unsupported")
    {{end}}
}
{{end}}

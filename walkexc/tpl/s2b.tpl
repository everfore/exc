{{test 2}}

{{range $k,$v := .Map}}
	hello, {{$k}}:{{$v}}!
	{{range $v}}
		name? {{.First}} {{.Second}}
	{{end}}
{{end}}

{{range $k,$v := .Arr}}
	hello, {{$k}}:{{$v}}!
{{end}}
{{test 2}}

{{range $k,$v := .Map}}
	hello, {{$k}}:{{$v}}!
	{{range $v}}
		name? {{.First}} {{.Second}}
	{{end}}
{{end}}

{{range $i,$v := .Arr}}
	arr[{{$i}}]:{{$v}}!
{{end}}

{{with .Arr}}
	arrs: {{join . ","}}
{{end}}
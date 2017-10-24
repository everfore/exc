echo {{test 2}}

{{range $k,$v := .Map}}
	echo hello, {{$k}}:{{$v}}!
	{{range $v}}
		echo name? {{.First}} {{.Second}}
	{{end}}
{{end}}

{{range $i,$v := .Arr}}
	echo arr[{{$i}}]:{{$v}}!
{{end}}

{{with .Arr}}
	echo arrs: {{join . ","}}
{{end}}
echo {{test 2}}

{{range $k,$v := .Map}}
	echo hello, {{$k}}:{{$v}}
	{{range $v}}
		echo name: {{.First}} {{.Second}}
	{{end}}
{{end}}

{{range $i,$v := .Arr}}
	echo arr[{{$i}}]:{{$v}}
{{end}}

{{with .Arr}}
	echo arrs: {{join . ","}}___
{{end}}

echo ...===123456

{{with .Arr}}
echo {{join . "_"}}_
{{end}}
echo 123
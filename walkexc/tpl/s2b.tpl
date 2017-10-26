
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

{{loop 1 10 1}}
{{loop 1 10 2}}
{{loop 1 10 3}}
{{loop 1 10 0}}
{{loop 1 0 0}}
{{loop -1 10 2}}
{{loop -1 -10 2}}

{{range $i:=loop 1 10 2}}
	{{$i}}th
{{end}}

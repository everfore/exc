## command

`tpl user.tpl user.yaml --execute`

## user.yaml

```
arr:
  - 1001
  - 1002
  - 1003
  - 1004
map:
  update:
    id: 1001
    name: John
    age: 24
  search:
    ids:
      - 1001
      - 1002
      - 1004
```

## user.tpl

```
# select * from user where id in ({{join .Arr ","}});
{{range $id:= .Arr}}
# select * from user where id = {{$id}};
{{end}}

## update
{{with .Map.update}}
curl 'http://localhost:8080/api/user/Update' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36' -H 'Content-Type: text/plain;charset:utf-8' -H 'Accept: /' -H 'Connection: keep-alive' -H 'Ajax: Y' -H 'Area: SG' --data-binary '{"customerId":{{.id}}, "name":{{.name}},"age":{{.age}}}'
{{end}}

## search
curl 'http://localhost:8080/api/user/Search' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36' -H 'Content-Type: text/plain;charset:utf-8' -H 'Accept: /' -H 'Connection: keep-alive' -H 'Ajax: Y' -H 'Area: SG' --data-binary '{"ids":[{{join .Map.search.ids ","}}]}'
```


## tpl result `tpl user.tpl user.yaml`

```
# select * from user where id in (1001,1002,1003,1004);

# select * from user where id = 1001;

# select * from user where id = 1002;

# select * from user where id = 1003;

# select * from user where id = 1004;


## update

curl 'http://localhost:8080/api/user/Update' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36' -H 'Content-Type: text/plain;charset:utf-8' -H 'Accept: /' -H 'Connection: keep-alive' -H 'Ajax: Y' -H 'Area: SG' --data-binary '{"customerId":1001, "name":John,"age":24}'


## search
curl 'http://localhost:8080/api/user/Search' -H 'Accept-Encoding: gzip, deflate' -H 'Accept-Language: zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36' -H 'Content-Type: text/plain;charset:utf-8' -H 'Accept: /' -H 'Connection: keep-alive' -H 'Ajax: Y' -H 'Area: SG' --data-binary '{"ids":[1001,1002,1004]}'
```
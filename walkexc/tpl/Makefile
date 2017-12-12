all:

init:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/toukii/e
	
tpl:
	e tpl init -o tpl
	go-bindata -pkg tpl tpl

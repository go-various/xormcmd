target=xorm
mac-target=xorm
win-target=xorm.exe
linux-target=xorm
build-dir=build
bin-dir=$(build-dir)/bin
mac:
	CGO_ENABLED=0  go build -o $(bin-dir)/$(mac-target)

win:
	CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go build -o $(bin-dir)bin/$(win-target) --tags=mysql

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(bin-dir)/$(linux-target) --tags=mysql

clean:
	/bin/rm -f $(bin-dir)*

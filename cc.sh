# rsrc.exe -arch amd64 -manifest main.manifest -ico main.ico -o rsrc.syso
# go build -ldflags "-H windowsgui" xxx.go

CGO_ENABLED=1 GOOS=windows go build -ldflags -w -o simtoolkit.exe

### Usage:
For macOS, 
`./src/main/logger -h`

### how to build?
```shell
# open your terminal for Unix, just copy that
cd distributed-log
export GOPATH=$PWD:$GOPATH
# get the info of your system
go env | grep GOOS
# output: GOOS="darwin"
go env | grep GOARCH
# output: GOARCH="amd64"
# with result below to build
GOOS=darwin GOARCH=amd64 go build -o logger ./src/main/main.go
```

then you can find it in `/distributed-log/logger`
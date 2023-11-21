# golang-utils
```go
go get -u github.com/gongsr-cn/golang-logs
```

示例
```go
package main
import(
    "fmt"
    Logs "github.com/gongsr-cn/golang-utils/logs"
)

func main() {
    config := Logs.Config{
		SerialNumber: 0,
		MaxSize:      100,
		Size:         0,
	}
	
	testLog, err := Logs.NewLogs("test", config)
	if err != nil {
	    fmt.Printf("create logs error:%s\n", err.Error())
	    return
	}
    testLog.Info("logs init success, this is in directory 'test'")
}

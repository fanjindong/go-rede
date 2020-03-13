# go-rede

:rocket:**A Rede is a fancy snooze delayed queue**


## Installation
```shell script
go get -u github.com/fanjindong/go-rede
```


## Quickstart

```go
import (
	"fmt"
	rede "github.com/fanjindong/go-rede"
	"time"
)

func main() {
	r := rede.NewClient(&rede.Options{
		Namespaces: "demo",
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DB:         0,
	})

	_, _ := rede.Push("a", 1*time.Second)
	_, _ := rede.Push("b", 1*time.Second)
	_, _ := rede.Push("c", 2*time.Second)

	time.Sleep(1 * time.Second)

	got, err := rede.Poll()
	if err != nil {
        panic(err)
	}
	fmt.Println(got)
	// ["a", "b"]
}
```

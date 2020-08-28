# Rede [中文](./doc/README_ZH.md)

:rocket:**A Rede is a fancy snooze delayed queue**

You can use the `Push` method to set a snooze time for an element. 
Unless the time comes, the element will not wake up. 
Get the collection of elements that have woken up through the `poll` method.


## Installation
```shell script
go get -u github.com/fanjindong/go-rede
```

## Features
- Snooze time can be updated
- Api is concise, such as `Push`, `Poll`
- Data persistent storage

## Quickstart

```go
import (
	"fmt"
	rede "github.com/fanjindong/go-rede"
	"time"
)

func main() {
	rd := rede.NewClient(&rede.Options{Namespaces: "demo", Addr: "127.0.0.1:6379"})

	_, _ := rd.Push("a", 1*time.Second)
	_, _ := rd.Push("b", 1*time.Second)
	_, _ := rd.Push("c", 2*time.Second)

	time.Sleep(1 * time.Second)
    
    cur := rd.Poll()
	for cur.Next() {
		got, _ := cur.Get()
		fmt.Println(got)
	}
	// out:
	// "a"
    // "b"
}
```

# go-rede

:rocket:**Rede 是一个具有贪睡功能的延时队列**

你可以通过 `push` 方法来为一个元素设置贪睡时间，除非时间到达，否则元素不会醒来。
通过 `poll` 方法来获取已经醒来的元素集。

## 安装
```shell script
go get -u github.com/fanjindong/go-rede
```

## 特性
- 元素的贪睡时间可以被更新
- Api 简洁
- 数据持久化存储

## 快速开始

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


# go-rede

:rocket: **Rede 是一个具有贪睡功能的延时队列**

![](SnoozyTheBear_ZH-CN1561515228_1920x1080.jpg)

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

## 需求
- Go >= 1.10.
- Redis >= 5.0.

## 快速开始

```go
func main() {
	rd := rede.NewClient(&rede.Options{Namespaces: "demo", Addr: "127.0.0.1:6379"})

	rd.Push("a", 1*time.Second)
	rd.Push("b", 1*time.Second)
	rd.Push("c", 2*time.Second)

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

## 使用说明

- Push

Push 一个元素到rede，并设置一个贪睡时间，直到这个时间到达，元素才会醒来。
- Pull

Pull 一个元素，注意此时会从rede中移除它。
- Look

查看一个元素的剩余贪睡时间。
- Ttn

查看rede中最快醒来的那个元素的剩余贪睡时间。
- Poll

获取已经醒来的元素们。
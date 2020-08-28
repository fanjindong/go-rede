# Rede [中文](./doc/README_ZH.md)

:rocket: **A Rede is a fancy snooze delayed queue**

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

## Requirements
- Go >= 1.10.
- Redis >= 5.0.

## Quickstart

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

## Usage

- Push

Push an element to rede and set a snooze time.
The element will not wake up until the time is up.
- Pull

Pull an element and remove it from rede.
- Look

View the remaining snooze time of an element.
查看一个元素的剩余贪睡时间。
- Ttn

查看rede中最快醒来的那个元素的剩余贪睡时间。
View the remaining snooze time of the element that wakes up fastest in rede.

- Poll

Poll the elements that have woken up.
获取已经醒来的元素们。
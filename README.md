# Sorted Linked List - Concurrent Safe 并发有序链表    
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/HelloLindo/Golang-Advanced-Camp-Language)

## About It
This is a concurrent safe sorted linked list implemented by GO language. Its advantages are 1) It can limit the write operation to a specific area, and allow one write and multi reads within the area, as well as global concurrent read and write. 2) Read operation is implemented with a completely lockless structure which uses atomic for unrestricted access to nodes.  
  
这是一个Go语言实现的并发有序链表，其优点是 1) 能够将写操作限制在某个区域，实现区域内的一写多读，全局并发读写。 2) 读操作采用了完全无锁的结构，使用atomic进行结点的无限制访问。
  
## Usage  
```go
    // Return true if the list contains a node with specific value.
    // 检查一个元素是否存在，如果存在则返回 true，否则返回 false
    Contains(value int) bool
    
    // Return true if insert a new Node x into the list successfully.
    // 插入一个元素，如果此操作成功插入一个元素，则返回 true，否则返回 false
    Insert(value int) bool
    
    // Return true if delete the Node with specific value from the list successfully.
    // 删除一个元素，如果此操作成功删除一个元素，则返回 true，否则返回 false
    Delete(value int) bool
    
    // Traverse all nodes in the list and stop the traverse if function f returns false.
    // 遍历此有序链表的所有元素，如果 f 返回 false，则停止遍历
    Range(f func(value int) bool)

    // Return the length of the list.
    // 返回有序链表的元素个数
    Len() int
```

## Test  
1. Run the Unit Test with command **`go test`**.  

2. Run the Concurrent Test with command **`go test -race`**.

## More  
[skipset](https://github.com/zhangyunhao116/skipset): A high-performance concurrent set based on skip list.  
[skipmap](https://github.com/zhangyunhao116/skipmap): A high-performance concurrent map based on skip list.  
[fastrand](https://github.com/zhangyunhao116/fastrand): Fastest pseudo-random number generator in Go.

## Thanks  
_ByteDance TechAcademy, Golang Advanced Camp - Language Section in ByteDance, [@zhangyunhao116](https://github.com/zhangyunhao116)_


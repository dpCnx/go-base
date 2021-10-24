### go

```
结构体里面嵌结构体只会有一个对象 如果是嵌套指针 会有多个对象
GMP M只会休眠 不会被销毁
map 只会扩容 不会缩减

grpc 初步判断也是短连接 只是基于http2会有复用

panic() 函数内部会产生一个关键的数据结构体 _panic
panic() 函数内部会执行 _defer 函数链条,如果执行完了都还没有恢复 _panic的状态, 那就没得办法了, 退出进程, 打印堆栈
_panic _defer 都是链表 都是挂在在 goroutine 之上
panic() 之后只会执行defer函数,所以recover只会在defer中才会成效,recover函数还会改变_panic的状态
```


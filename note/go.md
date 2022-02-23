### go

```
结构体里面嵌结构体只会有一个对象 如果是嵌套指针 会有多个对象
GMP M只会休眠 不会被销毁 m0 只会执行主线程的代码 退出程序就退出了
map 只会扩容 不会缩减

grpc 初步判断也是短连接 只是基于http2会有复用 
grpc的 timeout可以垮进程传递(https://toutiao.io/posts/o2a6ifo/preview)

panic() 函数内部会产生一个关键的数据结构体 _panic
panic() 函数内部会执行 _defer 函数链条,如果执行完了都还没有恢复 _panic的状态, 那就没得办法了, 退出进程, 打印堆栈
_panic _defer 都是链表 都是挂在在 goroutine 之上
panic() 之后只会执行defer函数,所以recover只会在defer中才会成效,recover函数还会改变_panic的状态

chan 里面的recevq sendq 是链表(链表不需要扩容)

defer 老版本就是链表 1.14之后就是普通函数调用

程序中遇到加锁的情况，goroutine--> runtime_SemacquireMutex -> sync_runtime_SemacquireMutex -> gopark ->schedule() 调度其他的goroutine
wake:semrelearse ->readyWithTime ->ready ->runqput 就和普通的goroutine一样了

schedt 结构体保存了调度器的信息  midle  muintptr （由空闲的工作线程组成的链表） pidle  puintptr（由空闲的 p 结构体对象组成的链表）
```


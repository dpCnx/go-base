```
    创建一个goroutine 2kb左右 用户态线程，是由 go runtime 管理，创建和销毁的消耗非常小，上下文切换成本低
    
    上下文包括:当前代码的位置，栈顶，栈底地址，状态等。
    
    Go 1.2之前 GM调度器 确定: 全局锁 性能问题
    P 决定了并行任务的数量，可通过 runtime.GOMAXPROCS 来设定。在 Go1.5 之后GOMAXPROCS 被默认设置可用的核数，而之前则默认为1。
    每一个p有一个本地队列 这样做的好处是不使用任何锁来 push/get 一个空闲的goroutine
    
    M0: 主线程启动时创建，创建完以后和普通的m是一样的
    G0: 每一个m会有一个G0 负责管理和调度goroutine
    
    
    
    创建:
        新创建的goroutine优先赋值给runnext(只会保存一个 go1.5) 
        如果runnext已经存在,则会帮助老的进入本地队列 local run queue 256长度的数组
        如果本地队列已满256，把本地队列拿出一半，塞入全局队列 global run queue 链表结构
         
    消费:
        p的schedtick = 61的时候，会去global run queue里面拿一个goroutine消费，其他时候正常执行。
        优先消费runnext里面的值
        然后消费local run queue 256长度的数组
        最后消费global run queue,如果global run queue里面有值，会取一半出来，执行第一个，剩余的放到本地队列，正常执行。                               
                               如果global run queue里面没有值，会去其他的p的队列里面偷一半，先执行最后一个，剩余的放到本地队列，正常执行。
                
                
    runtime  无法拦截的阻塞：syscall c语言调用 必须会占一个线程            
    
    系统监视器 (system monitor)，称为 sysmon：单独的线程，无需P就可以直接执行                           

        释放闲置超过5分钟的 span 物理内存；
        如果超过2分钟没有垃圾回收，强制执行；
        将长时间未处理的 netpoll 添加到全局队列；
        向长时间运行的 G 任务发出抢占调度； （10ms） go1.14
        收回因 syscall 长时间阻塞的 P；


    Spining thread 线程自旋 为了避免过多浪费 CPU 资源，自旋的 M 最多只允许 GOMAXPROCS (Busy P)  20us -> 1ms后翻倍 -> 10ms，不断重置
        类型1：M 不带 P 的找 P 挂载（一有 P 释放就结合）
        类型2：M 带 P 的找 G 运行（一有 runable 的 G 就执行）
    
    
    // ------------------？--------------------------------
    p 会维护一个本地 freelist 当 G 退出当前工作时，它将被 push 到这个空闲列表中
    调度器也有自己的列表。它实际上有两个列表：一个包含已分配栈的 G，另一个包含释放过堆栈的 G（无栈）
    1、有stack的是刚运行完的协程；2、没有的是运行完后经过一次gc了的，在并发标记阶段归还的；

```
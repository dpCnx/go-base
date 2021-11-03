```
    page 内存页
    span 内存块
    sizeclass 空间规格
    object 存储对象
    
    goroutine 的结构上带有mcache 申请内存小于32kb的时候 直接从mcache中获取mspan 不需要锁 申请内存小于16b的时候 划分为tiny对象
              大于32kb的没存直接从mheap中申请

    mcentral 维护全局的mspan资源 如果mcache里面没有对应大小的内存块 就需要从mcentral中获取 需要加锁
    mcentral 里维护着两个双向链表 nonempty 表示链表里还有空闲的 mspan 待分配 empty 表示这条链表里的 mspan 都被分配了object 或缓存 mcache 中
    mcentral 没有空闲的mspan的时候 会想mheap中申请
    mheap 每个 runtime.heapArena 都会管理 64MB 的内存


    每个 span 内有一个 bitmap allocBits 他表示上一次 GC 之后每一个 object 的分配情况 1：表示已分配 0：表示未使用或释放

    GC 将会启动去释放不再被使用的内存。在标记期间，GC 会用一个位图 gcmarkBits 来跟踪在使用中的内存。
    正在被使用的内存被标记为黑色，然而当前执行并不能够到达的那些内存会保持为白色。
    现在，我们可以使用 gcmarkBits 精确查看可用于分配的内存。Go 使用 gcmarkBits 赋值了 allocBits，这个操作就是内存清理。
    然而必须每个 span 都来一次类似的处理，需要耗费大量时间。Go 的目标是在清理内存时不阻碍执行，并为此提供了两种策略。

    Go 提供两种方式来清理内存：  如果超过2分钟没有触发，会强制触发 GC
    在后台启动一个 worker 等待清理内存，一个一个 mspan 处理
    当开始运行程序时，Go 将设置一个后台运行的 Worker(唯一的任务就是去清理内存)，它将进入睡眠状态并等待内存段扫描。
    当申请分配内存时候 lazy 触发
    当应用程序 goroutine 尝试在堆内存中分配新内存时，会触发该操作。清理导致的延迟和吞吐量降低被分散到每次内存分配时。

    https://lipeiru0329.github.io/2020/08/03/go-GC/
    
    
```
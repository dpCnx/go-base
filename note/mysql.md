### mysql

```
    MyISAM:每个MyISAM在磁盘上存储成三个文件。第一个文件的名字以表的名字开始，扩展名指出文件类型。.frm文件存储表定义。数据文件的扩展名为.MYD (MYData)。索引文件的扩展名是.MYI (MYIndex)
    InnoDB:所有的表都保存在同一个数据文件中（也可能是多个文件，或者是独立的表空间文件），InnoDB表的大小只受限于操作系统文件的大小，一般为2GB
    
    MyISAM 引擎使⽤B+Tree作为索引结构，叶节点的data域存放的是数据记录的地址。 
    InnoDB 引擎使⽤B+Tree作为索引结构，这棵树的叶节点data域保存了完整的数据记录。

    InnoDB,MyISAM 不⽀持Hash索引 
    Memory ⽀持Hash索引
    
    hash索引:
        等值查询快,可以通过索引值直接定位
        范围查询⽐较麻烦 
        结构没有顺序, 对数据排序需要重新进⾏排序 
        不能使⽤部分索引, ⽐如前缀索引, 以及像 like 'xxx%'模糊查询 如果⼤量的hash值相同, 会出现hash碰撞, 也会影响性能
        
    聚簇索引：将数据存储与索引放到了⼀块，找到索引也就找到了数据 (innodb)
    ⾮聚簇索引：将数据存储于索引分开结构，索引结构的叶⼦节点指向了数据的对应 (myisam)
    
    4.0版本以下，varchar(100)，指的是100字节，如果存放UTF8汉字时，只能存33个（每个汉字3字节） 
    5.0版本以上，varchar(100)，指的是100字符，无论存放的是数字、字母还是UTF8汉字（每个汉字3字节），都可以存放100个
    
    Using index condition：索引下推在非主键索引上的优化，可以有效减少回表的次数，大大提升了查询的效率。
    eg:SELECT * from user where  name like '陈%' and age=20 InnoDB并没有忽略age这个字段，而是在索引内部就判断了age是否等于20，对于不等于20的记录直接跳过，因此在(name,age)这棵索引树中只匹配到了一个记录，此时拿着这个id去主键索引树中回表查询全部数据，这个过程只需要回表一次。（5.6之后）
    
    Using index：查询数据直接遍历⼆级索引树就可以返回结果的情况, 不再需要遍历主键索引树(不需要回表操作)去 查询数据的情况, 我们称之为索引覆盖
    Using where：回表之后进行条件查询
    
    范围查询：使⽤组合索引的时候, 范围查询是可以使⽤索引的, 但是对应的范围查询必须是第⼀索引
    eg:SELECT * FROM employees WHERE emp_no<'10010' and title='Senior Engineer'
        
    驱动表与被驱动表 可以通过explain 查看 先执行的为驱动表
    
    Index Nested-LoopJoin: 通过index查询
    Block Nested-Loop Join： 加载多条信息到内存
    Simple Nested-Loop Join： 一条一条信息加载到内存 (mysql 不使用)
    
    双路排序:只加载主键和排序字段到sort_buffer进行排序,需要进行回表查询
    单路排序:加载主键,需要查询的字段和排序字段到sort_buffer进行排序,直接输出结果
    	根据参数 max_length_for_sort_data 决定, 如果查询的字段 ⻓度不超过改参数值, 就使⽤单路排序,否则, 采⽤双路排序
 
 
 
    MySQL 8.0版本直接将查询缓存的整块功能删掉了
   
    MySQL里经常说到的WAL技术，WAL的全称是Write-Ahead Logging，它的关键点就是先写日志，再写磁盘
    redolog（重做日志) innoDB引擎特有
    binlog（归档日志）
    redolog到binlog 两阶段提交 保证了数据的完整 可恢复
    binlog的写入逻辑比较简单：事务执行过程中，先把日志写到binlog cache，事务提交的时候，再把binlog cache写到binlog文件中
    
    联合索引的时候，联合索引对应的数据都会记录在索引树上
    
    表级锁 MDL（metadata lock)
    MySQL 5.5版本中引入了MDL，当对一个表做增删改查操作的时候，加MDL读锁；当要对表做结构变更操作的时候，加MDL写锁
     一致性读、当前读和行锁
    在InnoDB事务中，行锁是在需要的时候才加上的，但并不是不需要了就立刻释放，而是要等到事务结束时才释放。这个就是两阶段锁协议
    发现死锁后，主动回滚死锁链条中的某一个事务，让其他事务得以继续执行。将参数innodb_deadlock_detect设置为on，表示开启这个逻辑
    更新数据都是先读后写的，而这个读，只能读当前的值，称为“当前读”（current read）
    
    查询过程都是按页读取
    change buffer 可以持久化的数据。change buffer在内存中有拷贝，也会被写入到磁盘上 change buffer用的是buffer pool里的内存
    当需要更新一个数据页时，如果数据页在内存中就直接更新，而如果这个数据页还没有在内存中的话，在不影响数据一致性的前提下，InooDB会将这些更新操作缓存在change buffer中
    将change buffer中的操作应用到原数据页，得到最新结果的过程称为merge。除了访问这个数据页会触发merge外，系统有后台线程会定期merge。在数据库正常关闭（shutdown）的过程中，也会执行merge操作。
    
    mysql选错索引主要是因为采样率，可以使用force index 强行矫正
    
    当内存数据页跟磁盘数据页内容不一致的时候，我们称这个内存页为“脏页”。内存数据写入到磁盘后，内存和磁盘上的数据页的内容就一致了，称为“干净页”。
    更新内存的时机:对应redolog记录满了，需要擦除一些
                 系统内存不足
                 MySQL认为系统“空闲”的时候
                 
    buffer pool: InnoDB用缓冲池管理内存
    
    只是delete掉表里面不用的数据的话，表文件的大小是不会变的，数据页的复用。你还要通过alter table命令重建表，才能达到表文件变小的目的
    
    MyISAM引擎把一个表的总行数存在了磁盘上，因此执行count(*)的时候会直接返回这个数，效率很高；
    而InnoDB引擎就麻烦了，它执行count(*)的时候，需要把数据一行一行地从引擎里面读出来，然后累积计数。
    
    
    sort_buffer:Extra这个字段中的“Using filesort”表示的就是需要排序，MySQL会给每个线程分配一块内存用于排序，称为sort_buffer。
    order by 快排  内存放不下时，就需要使用外部排序，外部排序一般使用归并排序算法 有点类似hadoop
    
    两个表的字符集不同，一个是utf8，一个是utf8mb4，所以做表连接查询的时候用不上关联字段的索引
    隐式类型转换,使索引失效
    对索引字段做函数操作,使索引失效
    
    间隙锁(Gap Lock) 间隙锁和行锁合称next-key lock，每个next-key lock是前开后闭区间。锁是加在索引上的 --> 保证不幻读
    
    加锁规则里面，包含了两个“原则”、两个“优化”和一个“bug”。
        原则1：加锁的基本单位是next-key lock。希望你还记得，next-key lock是前开后闭区间。
        原则2：查找过程中访问到的对象才会加锁。
        优化1：索引上的等值查询，给唯一索引加锁的时候，next-key lock退化为行锁。
        优化2：索引上的等值查询，向右遍历时且最后一个值不满足等值条件的时候，next-key lock退化为间隙锁。
        一个bug：唯一索引上的范围查询会访问到不满足条件的第一个值为止。
 
    简单的select操作属于快照读，如下所示。
    select * from table where xxx;
    
    特殊读、插入、更新、删除操作，属于当前读，需要加锁，如下所示。
    select * from table where xxx lock in share mode;
    select * from table where xxx for update;
    insert into table values(xxx);
    update table set xxx where xxx;
    delete from table where xxx;
```


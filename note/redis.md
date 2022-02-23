```
        redis  布隆过滤器的底层就是基于bitmap的  只是会经过很多次的hash  可以自己基于bitmap简单实现布隆过滤器
        redis  集群中有master挂掉了 哨兵集群会根据选举的leader去指定一个slave为master leader通常根据slave的稳定性来选举master
```
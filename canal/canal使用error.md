```sql

error:官方demo使用的是1.1 win 通过docker启动会是失败

解决：在用户目录添加
    [wsl2]
    kernelCommandLine = vsyscall=emulate

error:通过client链接canal  报错：unsupported version at this client.
    
解决：docker 镜像使用最新的
    
mysql
    
初始化： 
    1.添加用户
        CREATE USER canal IDENTIFIED BY 'canal';
        GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
        FLUSH PRIVILEGES;

    2.bin-log 日志使用row

error:canal caching_sha2_password Auth failed
解决：
    ALTER USER 'canal'@'%' IDENTIFIED WITH mysql_native_password BY 'canal';
    FLUSH PRIVILEGES    

error:sqlstate = HY000 errmsg = Could not find first log file name in binary log index file
解决：
    show master status 写入正确的bin日志和pos

error:com.taobao.tddl.dbsync.binlog.LogDecoder - Skipping unrecognized binlog event Unknown from: binlog.000027:1488
解决： （暂时未解决）
    初步判断是mysql版本的问题
    https://github.com/alibaba/canal/issues/4081
    有issues大概看了一下
    
注意：监听mysql中的表使用的是正则表达式，复制到代码中的时候注意转移字符
```


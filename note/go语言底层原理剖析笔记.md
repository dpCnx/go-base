###### 第一章

```
三阶段编译器:
	编译器前段:理解源代码
	优化器:优化代码,识别冗余代码,识别内存逃逸
	编译器后端:生成特定目标机器码上的程序
	
词法解析
语法解析
抽象语法树构建
类型检查
变量捕获: 主要针对闭包
        func main(){
            a:=1
            b:=2
            go func(){
                fmt.Println(a,b)
            }()
            a = 99
        }
	由于变量a在闭包之后进行了其他赋值操作,在闭包中,采取地址引用的方式对变量a进行操作。变量b通过值传递的方式操作。
	go tool compile -m=2 .\test.go | grep capturing	

 	main capturing by ref: a (addr=false assign=true width=8)  ref 地址引用
 	main capturing by value: b (addr=false assign=false width=8)	value 值传递

函数内联: 减少函数调用带来的开销。函数过于复杂,不会执行内联
逃逸分析: 编译器会尽可能的将变量放置在栈中,栈中的对象随着函数调用结束会被自动销毁,减轻运行时分配和垃圾回收的负担。
	原则: 指向栈上对象的指针不能被存储到堆中
		 指向栈上对象的指针不能超过该对象的生命周期
闭包重写：闭包只调用一次,会被转换为普通函数。如果被多次调用,会创建闭包对象
遍历函数
SSA生成：静态单赋值(static single assignment)
机器码生成 --- 汇编器
机器码生成 --- 链接
	静态链接
	动态链接: 外部库,例如引用的c代码
ELF文件解析
```
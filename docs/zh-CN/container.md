## 简介

`Container`是`Firmeve`的基础，一切都是以`Container`为基础实现`ioc`开始。



## 基础示例

### 新建容器

```go
var c = container.New()
```



### 数据绑定

```go
// 单例数据绑定
foo := new(Foo)
c.Bind("foo", foo)

// 绑定函数
var f = func() string {
  return "rand string"
}
c.Bind("f", f)

// 绑定slice
var s = []string{"a", "b"}
c.Bind("s", s)

// 强制覆盖已存在的单例
c.Bind("foo", new(Foo2), WithCover(true))
```

> **注意**
>
> 1. 绑定的类型仅支持：array, slice, map, func, prt(struct), struct
> 2. 只有函数类型是非单例类型，每次调用时会自动执行函数方法来获取结果
> 3. 更多高级用户请参见 [container_test](https://github.com/firmeve/firmeve/blob/develop/container/container2_test.go)



### 数据获取

```go
// 判断原有对象是否存在
if c.Has("Foo") {
	// 获取一个已存在的对象
  c.Get("Foo").(*Foo)
}

// 删除一个对象
c.Remove("Foo")

// 清除容器所有对象
c.Flush()
```



### 实例解析

通过`Make`方法可以轻松完成自动解析，目前解析对象包括：
- 容器中已存在对象
- 结构体
- 函数

```go
// 解析容器中已存在的对象，同Get()
c.Make("Foo")

// 解析结构体
type Person struct {
   Bar *Bar `inject:"bar"`
   Foo *Foo `inject:"foo"`
}
c.Make(new(Person))

// 解析函数

// 普通标量函数
func NormalSample(id int, str string) (int, string) {
   return id, str
}
fmt.Printf("%v", c.Make(NormalSample, 10, `foo`))

// 支持多级函数嵌套
// 假设 *Nesting *Sample 已存在于容器中
func StructFuncFunc(id int, nesting *Nesting, fn func(nesting *Nesting) int, prt *Sample) int {
   return id + fn(nesting) + prt.Id
}

// 容器会自动解析已存在的 *Nesting *Sample
// 对于手动定义的参数，必须放在函数前置
fmt.Printf("%v", c.Make(StructFuncFunc, 10))
```

> 注意：
>
> - 如果解析注入的对象不是指针类型则会自动忽略，不会自动赋值。
> - 带有动态参数`...params`的函数暂时无法解析



## 单独使用

```bash
go get -u github.com/firmeve/firmeve/container
```


## 简介
Ioc 单实例容器绑定

## 基础示例

### 准备数据
```go
type (
    Nesting struct {
        NId        int
        NTitle     string
        NBoolean   bool
        NArray     [3]int
        NMapString map[string]string
        NSlice     []string
        nprivate   string
    }
    
    Sample struct {
    	Id        int
    	Title     string
    	Boolean   bool
    	Array     [3]int
    	MapString map[string]string
    	Slice     []string
    	private   string
    }
)

func Sample() *structs.Sample {
	return &Sample{
		Id: 15,
	}
}
```

以下只是简单`Api`的示例
[**更多示例请参考此处**](https://github.com/firmeve/firmeve/blob/develop/container/container2_test.go)

### 新建一个基础容器
```go
// 新建一个容器
c := container.New()
```

### 结构体绑定

```go
c.Bind("nesting", &Nesting{
    NId: 10,
})
```

> 默认非函数型绑定是单例对象

### 函数绑定

```go
c.Bind("sample", Sample)
```

> 函数绑定不支持单例

### 自动解析

```go
fmt.Printf("%#v\n", c.Make("sample"))
```

> 注意：暂时不支持...params动态扩展参数函数解析

### 强制已绑定值覆盖

我们通过`WithCover(true)`方法可以轻松设定覆盖参数

```go
c.Bind("nesting", &Nesting{
    NId: 10,
})

c.Bind("nesting", &Nesting{
    NId: 11,
}, WithCover(true))
```

## 对象获取

### 验证对象是否存在

在不确定容器中是否包含此对象时，请使用`Has`方法进行判断

```go
c.Has(`foo`)
```
### 对象获取

直接获取容器中的值
```go
c.Get(`foo`)

// 或
c.Make(`foo`)
```
> 注意：`Get`只能获取已存在的对象，如果需要自动解析新对象请使用`Make`方法

### 对象删除
```go
// 清除指定对象
c.Remove(`foo`)

// 全部清除
c.Flush()
```

## 自动解析

通过`Make`方法可以轻松完成自动解析，目前解析对象包括
- 容器中已存在对象
- 结构体
- 函数

### 函数自动解析

函数解析中，如果遇到的参数类型已存在于容器中，并且参数未提供，则会自动注入容器中的参数。
如果存在，则会使用提供的函数参数，标量类型的参数则必须提供

让我们来看一个简单的示例

```go
//Sample函数如上
fmt.Printf("%#v\n", c.Make(Sample))

func NormalSample(id int, str string) (int, string) {
	return id, str
}

fmt.Printf("%#v\n", NormalSample, 10, `foo`)

func StructFunc(ptr *Sample, nesting structs.Nesting) int {
	return ptr.Id + nesting.NId
}

fmt.Printf("%#v\n", StructFunc)

// 支持多级函数嵌套
func StructFuncFunc(id int, nesting *Nesting, fn func(nesting *Nesting) int, prt *Sample) int {
	return id + fn(nesting) + prt.Id
}
fmt.Printf("%#v\n", StructFuncFunc, 10)

// 更复杂的解析
func MultipleParamSample(ptr *Sample, nestingPtr *Nesting, nesting structs.Nesting, notPtr structs.Sample) (*Sample, *Nesting, Nesting, Sample) {
	return ptr, nestingPtr, nesting, notPtr
}
fmt.Printf("%#v\n", MultipleParamSample)
```

> 带有动态参数...params的函数暂时无法解析

### 结构体tag自动解析

tag解析的对象，必须容器中已存在，并且是指针类型，普通结构体无法设置，会自动忽略。

```go

c.Bind(`foo`, FooStruct{})
c.Bind(`bar`, FooStruct{})

type Person struct {
	Foo1 PersonName `inject:"foo"`
	Bar *Bar `inject:"bar"`
	Foo *Foo `inject:"foo"`
}

fmt.Printf("%#v", c.Make(new(NewPerson)))
```
> 注意：此时 `Person`中的`Foo1`字段并不是指针类型，所以会自动忽略。

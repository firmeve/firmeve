## 简介
Ioc

## 对象绑定

```go
f := NewFirmeve()
```

### 标量绑定
```go
f.Bind(`bool`, false)
```
```go
f.Bind(`string`, "string")
```

### `Struct` 绑定
假设我们有如下`struct`
```go
type Foo struct {
    Bar string	
}

func NewFoo(bar string) Foo{
    return Foo{Bar:bar}
}
```
我们将`foo`绑定进我们的容器中
```go
f.Bind(`foo`,NewFoo("abc"))
```

### `Prt` 绑定
绑定`prt`的`struct`类型是我们最常用的一种方式，请看下面的示例
```go
type Baz struct {
    Baz string	
}

func NewBaz() *Baz {
	return &Baz{}
}

f.Bind(`baz`, NewBaz())
```

### 函数绑定
```go
func NewFoo() Foo{
	return Foo{Bar:"default string"}
}

f.Bind(`foo`, NewFoo)
```
> 使用`Firmeve`在绑定的时候暂时不支持参数注入
> 注意：如果是非函数类型绑定，则会默认为单实例类型

### 已绑定值覆盖
我们通过`WithBindCover(true)`方法可以轻松设定覆盖参数
```go
f.Bind(`bool`, false)
fmt.Printf("%t",f.Get(`bool`))
f.Bind(`bool`, true, WithBindCover(true))
fmt.Printf("%t",f.Get(`bool`))
```

## 对象获取

### 验证对象是否存在
在不确定容器中是否包含此对象时，请使用`Has`方法进行判断
```go
f.Has(`foo`)
```
### 对象获取
直接获取容器中的值
```go
f.Get(`foo`)
```
> 注意：如果指的定的`key`不存在，则会抛出`panic`错误

### 新对象获取
```go
func NewFoo() Foo{
	return Foo{Bar:"default string"}
}

f.Bind(`foo`, NewFoo)
fmt.Printf("%p\n",f.Get("foo"))
fmt.Printf("%p\n",f.Get("foo"))
fmt.Printf("%p\n",f.Get("foo"))
```
> 在`Firmeve`中，如果需要每次重新得到一个新的结构体对象
必须绑定一个函数值，否则得到的将是一个单实例

### 单例获取
```go
func NewFoo() Foo{
	return Foo{Bar:"default string"}
}

f.Bind(`foo`, NewFoo())
fmt.Printf("%p\n",f.Get("foo"))
fmt.Printf("%p\n",f.Get("foo"))
fmt.Printf("%p\n",f.Get("foo"))
```

## 参数解析
### 函数自动解析
让我们来看一个简单的示例
```go
type PersonName struct {
	Name string
}

func NewPersonName() PersonName {
	return PersonName{Name: "Firmeve"}
}

type PersonAge struct {
	Age int
}

func PersonAge() *PersonAge {
	return &PersonAge{Age: 1}
}

type Person struct {
	name PersonName
	age *PersonAge
}

func NewPerson(name PersonName,age *PersonAge) *NewPerson {
    return NewPerson{
    	name: name
    	age: age
    }
}

f.Bind("PersonName", NewPersonName)
f.Bind("PersonAge", PersonAge)
fmt.Printf("%#v", f.Resolve(NewPerson))
```

### 结构体tag自动解析
现在，让我们修改下上面的`Person`
```go
type Person struct {
	Name PersonName `inject:"PersonName"`
	Age *PersonAge `inject:"PersonAge"`
	age1 *PersonAge `inject:"PersonAge"`
}
```
然后我们使用`new`函数直接创建一个新的结构体指针
```go
fmt.Printf("%#v", new(NewPerson))
```
> 注意：此时 `Person`中的`Name`字段并不是指针类型，而`age1`不符合`struct`的`tag`规范，所以`Firmeve`都会自动忽略。

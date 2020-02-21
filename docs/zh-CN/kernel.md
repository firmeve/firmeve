## 内核模块 Kernel
`Firmeve` 是面向接口编程的框架，在内核模块中提供了所有的契约接口以及基础结构。
可以说`Kernel`是`Firmeve`的所有核心和灵魂。

`kernel/contract` 是整个Firmeve的运行约束，主要提供运行接口以及对应接口会需要的静态常量。
`kernel/*` 内核基础功能实现，或者提供一些常用基础结构体，如`BaseProvider`

严格来说`kernel`只会调用`support`中原生或第三方的扩展方法或函数，`kernel`是最底层的约束

## 扩展接口

### 接口命名规则

必须以模块名称打头，或直接使用模块名称。

如：`cache`模块 分别接口为：
- `Cache`
- `CacheStore`
- `CacheSerializable`
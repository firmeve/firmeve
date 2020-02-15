## 内核模块 Kernel
Firmeve 是面向接口编程的框架，在内核模块中提供了所有的契约接口以及基础结构。
可以说Kernel是Firmeve的所有核心和灵魂。

## 接口命名规则

以模块名称打头，或直接使用模块名称（单一接口情况下）

如：`cache`模块 分别接口为：
- `CacheStore`
- `CacheSerializable`
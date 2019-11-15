- 支持 Server Fields 选项
- 支持 Server Scene 选项
- 支持传递的 Fields 选项-> 优先级高于Server Field
- 支持传递的 Scene 选项-> 优先级高于Server Scene
- 支持scene和field组合？？

Transformer : 资源转换器 ，主要把strcut字段对应的转换成相应的格式化

Resource : 资源处理器

Serialize : 资源打包包装器

m := Manager.New()
m.item(resource,transformer)

func item(resource,transformer) {
    new Item(new transformer(resource))
    .Only(fields)
    .OnlyScene(scene)
    .Resolve()
}
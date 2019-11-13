只支持Array,Slice,Map,Struct,Prt,Func

F.Bind("")


binding struct {
    name,
    prototype,
}

instance struct {
    name
    instance
}

resolvedType struct {
    
}

Get() 只能取字符串的存在
Make() 是重新解析，但如果是instance则也会返回原单例

bindings
instances
make 函数可以重复解析，传入不同的参数

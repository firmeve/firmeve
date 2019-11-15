package serializer

type ResolveData map[string]interface{}

type Resolver interface {
	Resolve() interface{}
}

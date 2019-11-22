package serializer

type (
	Resolver interface {
		Resolve() interface{}
	}

	ResolveData map[string]interface{}
)

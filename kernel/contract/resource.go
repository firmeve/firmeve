package contract

type (
	ResourceMetaData map[string]interface{}

	ResourceData map[string]interface{}

	ResourceDataCollection []ResourceData

	ResourceLinkData map[string]string

	ResourceFields []string

	ResourceMeta interface {
		SetMeta(meta ResourceMetaData)
		Meta() ResourceMetaData
	}

	ResourceLink interface {
		SetLink(links ResourceLinkData)
		Link() ResourceLinkData
	}

	ResourceDatable interface {
		Data() ResourceData
	}

	ResourceCollectionData interface {
		CollectionData() ResourceDataCollection
	}

	ResourceResolver interface {
		Resolve() interface{}
	}

	ResourceResolveData map[string]interface{}

	ResourceTransformer interface {
		Resource() interface{}
		SetResource(resource interface{})
	}
)

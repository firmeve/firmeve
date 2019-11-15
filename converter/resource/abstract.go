package resource

type Meta map[string]interface{}
type Data map[string]interface{}
type DataCollection []Data
type Link map[string]string

type IMeta interface {
	SetMeta(meta Meta)
	Meta() Meta
}

type ILink interface {
	SetLink(links Link)
	Link() Link
}

type Datable interface {
	Data() Data
}
type CollectionData interface {
	CollectionData() DataCollection
}

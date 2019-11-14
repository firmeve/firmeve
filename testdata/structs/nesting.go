package structs

import "time"

type Nesting struct {
	NId        int
	NTitle     string
	NBoolean   bool
	NArray     [3]int
	NMapString map[string]string
	NSlice     []string
	nprivate   string
}

type Sample struct {
	Id        int
	Title     string
	Boolean   bool
	Array     [3]int
	MapString map[string]string
	Slice     []string
	private   string
}

type SampleNesting struct {
	Sample
	Nesting      *Nesting
	NestingValue Nesting
}

type SampleNestingInject struct {
	Sample
	Nesting  *Nesting `inject:"nesting"`
	Nesting2 *Nesting
}

type SampleNestingPtr struct {
	Sample
	Nesting *Nesting
}

type Main struct {
	id        int
	title     string
	show      bool
	hideSub   hideSub
	Sub       Sub
	PublicKey string
	PrtSub    *Sub
	PrtSub2   *Sub `inject:"sub"`
	//Array     [2]int
	Array [3]int
	Slice []string
	Map   map[string]string
}
type Sub struct {
	id           int
	title        string
	time         time.Time
	SubPublicKey string
}
type hideSub struct {
	id    int
	title string
}

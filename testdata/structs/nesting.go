package structs

import "time"

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

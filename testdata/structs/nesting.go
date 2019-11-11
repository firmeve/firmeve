package structs

import "time"

type Main struct {
	id        int
	title     string
	show      bool
	hideSub   hideSub
	sub       Sub
	PublicKey string
	PrtSub    *Sub
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

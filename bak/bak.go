package bak

type Bak struct {
	Title string
	At    int
	cast   string `di:"singleton"`
}
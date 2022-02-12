package tpl

type P struct {
	Name  string
	Files []F
	Child []P
}

type F struct {
	Name    string
	Content string
}

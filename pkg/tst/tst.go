package tst

type Name struct {
	name string
}

func (n *Name) New(name string) {
	n.name = name
}

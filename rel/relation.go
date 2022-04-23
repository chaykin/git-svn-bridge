package rel

type Relation struct {
	parent string
	child  string
}

func New(parent, child string) *Relation {

	return &Relation{parent: parent, child: child}
}

func (r *Relation) GetParent() string {
	return r.parent
}

func (r *Relation) GetChild() string {
	return r.child
}

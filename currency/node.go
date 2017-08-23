package currency

type CNode struct {
	root *node
	now  *node
	len  int
}

type node struct {
	value int
	next  *node
}

func (n *CNode) Size() int {
	return n.len
}

func (n *CNode) Add(value int) {
	n.len++

	part := new(node)
	part.value = value

	if n.root == nil {
		n.root = part
		n.now = part
	} else {
		n.now.next = part
		n.now = part
	}
}

func (n *CNode) RemoveFirst() {
	if n.root == nil {
		return
	}

	n.len--
	n.root = n.root.next
}

func (n *CNode) ToSlice(cap int) *[]int {
	ar := make([]int, 0, cap)
	for tmp := n.root; tmp != nil; tmp = tmp.next {
		ar = append(ar, tmp.value)
	}

	return &ar
}

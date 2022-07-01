package art

type Art struct {
	root *node
}

type node struct {
	key      []byte
	val      []byte
	children []*node
	size     int16
}

func New() *Art {
	return &Art{}
}

func (a *Art) Set(key, val []byte) (replaced bool) {
	//fmt.Println("Set", key)
	if a.root == nil {
		a.root = &node{
			key: key,
			val: val,
		}
		return
	}
	return a.root.set(key, val, 0)
}

func (a *Art) Get(key []byte) (val []byte) {
	//fmt.Println("Get", key)
	if a.root == nil {
		return nil
	}
	return a.root.get(key, 0)
}

func (a *Art) String() string {
	return a.root.String(0)
}

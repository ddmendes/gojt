package node

// Node is a generic JSON element
type Node interface {
	Get(token string) (Node, error)
	GetKeys() ([]string, error)
	GetInterface() interface{}
}

package projection

import (
	"sync"
	"strconv"
)

type Projection struct {
	contextId	string
	Version		uint64
	root		*node
}

func NewProjection(contextId string) *Projection {
	return &Projection{
		contextId: contextId,
		root: newNode(),
	}
}

type node struct {
	children	map[string]*node
	value		map[uint64]interface{}
	deleted		map[uint64]bool
	mutex		sync.Mutex
}

func newNode() *node {
	return &node{
		children:	make(map[string]*node),
		value: 		make(map[uint64]interface{}),
		deleted:	make(map[uint64]bool),
	}
}

func (prj *Projection) Upsert(key string, value interface{}) error {
	parts, err := ParseKey(key)
	if err != nil {
		return err
	}

	return upsertValue(parts, value, prj.root, prj.Version)
}

func (prj *Projection) Delete(key string) error {
	parts, err := ParseKey(key)
	if err != nil {
		return err
	}

	return deleteValue(parts, prj.root, prj.Version)
}

func (prj *Projection) Get(key string) (uint64, interface{}, error) {
	parts, err := ParseKey(key)
	if err != nil {
		return 0, nil, err
	}
	return getValue(parts, prj.root, prj.Version)
}

func (prj *Projection) Commit() error {
	prj.Version++
	return nil
}

func (prj *Projection) Rollback() error {
	return rollback(prj.root, prj.Version)
}

func upsertValue(parts []string, value interface{}, n *node, version uint64) error {
	if len(parts) == 0 {
		return nil
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	child, ok := n.children[parts[0]]
	if !ok {
		child = newNode()
		n.children[parts[0]] = child
	}

	child.deleted[version] = false

	if len(parts) > 1 {
		child.value[version] = nil
		return upsertValue(parts[1:], value, child, version)
	} else {
		if err := deleteChildren(child, version); err != nil {
			return err
		}
		child.value[version] = value
	}

	return nil
}

func deleteValue(parts []string, n *node, version uint64) error {
	if len(parts) == 0 {
		return nil
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if child, ok := n.children[parts[0]]; ok {
		if len(parts) > 1 {
			return deleteValue(parts[1:], child, version)
		} else {
			child.deleted[version] = true
			return deleteChildren(child, version)
		}
	}

	return nil
}

func deleteChildren(n *node, version uint64) error {
	for _, child := range n.children {
		child.deleted[version] = true
		if err := deleteChildren(child, version); err != nil {
			return err
		}
	}
	return nil
}

func getValue(parts []string, n *node, version uint64) (uint64, interface{}, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if isDeleted(n, version) {
		return 0, nil, nil
	}
	if len(parts) > 0 {
		child, ok := n.children[parts[0]]
		if ok {
			ver, res, err := getValue(parts[1:], child, version)
			return ver, res, err
		} else {
			return 0, nil, nil
		}
	}

	ver, res := getTreeValue(n, version)
	return ver, res, nil
}

func isDeleted(n *node, version uint64) bool {
	maxV := uint64(0)
	deleted := false
	for v,d := range n.deleted {
		if v <= version && v >= maxV {
			deleted = d
			maxV = v
		}
	}
	return deleted
}

func getLeafValue(n *node, version uint64) (uint64, interface{}) {
	maxV := uint64(0)
	var val interface{}
	for v,vv := range n.value {
		if v <= version && v >= maxV {
			val = vv
			maxV = v
		}
	}
	for v,del := range n.deleted {
		if v <= version && v >= maxV && del {
			val = nil
			maxV = v
		}
	}
	return maxV, val
}

func getTreeValue(n *node, version uint64) (uint64, interface{}) {
	if isDeleted(n, version) {
		return 0, nil
	}
	ver, val := getLeafValue(n, version)
	if val != nil {
		return ver, val
	}

	treeMaxVer := uint64(0)
	var tree map[string]interface{}
	for key, child := range n.children {
		child.mutex.Lock()
		ver, childVal := getTreeValue(child, version)
		child.mutex.Unlock()
		if childVal != nil {
			if tree == nil {
				tree = make(map[string]interface{})
			}
			tree[key] = childVal
			if ver > treeMaxVer {
				treeMaxVer = ver
			}
		}
	}
	if tree == nil {
		return 0, nil
	}

	// Check if tree is a slice
	isSlice := true
	for i:=0; i < len(tree); i++ {
		if _, ok := tree[strconv.Itoa(i)]; !ok {
			isSlice = false
			break
		}
	}

	if (isSlice) {
		slicedTree := make([]interface{}, 0, len(tree))
		for i:=0; i < len(tree); i++ {
			val, _ := tree[strconv.Itoa(i)]
			slicedTree = append(slicedTree, val)
		}
		return treeMaxVer, slicedTree
	} else {
		return treeMaxVer, tree
	}
}

func rollback(n *node, version uint64) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	delete(n.deleted, version)
	delete(n.value, version)
	for _, child := range n.children {
		if err := rollback(child, version); err != nil {
			return err
		}
	}
	return nil
}
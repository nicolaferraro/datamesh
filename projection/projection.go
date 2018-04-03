package projection

import (
	"strings"
	"regexp"
	"errors"
)

const (
	PathSeparator			= "."

	ErrorIllegalKey			= "EDB_ILLEGAL_KEY"
	ErrorNoRootSpecified	= "EDB_NO_ROOT_SPECIFIED"

	keyPartPattern			= "^[A-Za-z0-9]+[A-Za-z0-9_-]*$"
)

type Projection struct {
	Version		uint64
	root		*node
}

func NewProjection() *Projection {
	return &Projection{
		root: newNode(),
	}
}

type node struct {
	children	map[string]*node
	value		map[uint64]interface{}
	deleted		map[uint64]bool
}

func newNode() *node {
	return &node{
		children:	make(map[string]*node),
		value: 		make(map[uint64]interface{}),
		deleted:	make(map[uint64]bool),
	}
}

func (prj *Projection) Upsert(key string, value interface{}) error {
	parts, err := parseKey(key)
	if err != nil {
		return err
	}

	return upsertValue(parts, value, prj.root, prj.Version)
}

func (prj *Projection) Delete(key string) error {
	parts, err := parseKey(key)
	if err != nil {
		return err
	}

	return deleteValue(parts, prj.root, prj.Version)
}

func (prj *Projection) Get(key string) (interface{}, error) {
	parts, err := parseKey(key)
	if err != nil {
		return nil, err
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

func getValue(parts []string, n *node, version uint64) (interface{}, error) {
	if isDeleted(n, version) {
		return nil, nil
	}
	if len(parts) > 0 {
		child, ok := n.children[parts[0]]
		if ok {
			res, err := getValue(parts[1:], child, version)
			return res, err
		} else {
			return nil, nil
		}
	}

	res := getTreeValue(n, version)
	return res, nil
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

func getLeafValue(n *node, version uint64) interface{} {
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
	return val
}

func getTreeValue(n *node, version uint64) interface{} {
	if isDeleted(n, version) {
		return nil
	}
	val := getLeafValue(n, version)
	if val != nil {
		return val
	}

	var tree map[string]interface{}
	for key, child := range n.children {
		childVal := getTreeValue(child, version)
		if childVal != nil {
			if tree == nil {
				tree = make(map[string]interface{})
			}
			tree[key] = childVal
		}
	}
	if tree == nil {
		return nil
	}
	return tree
}

func rollback(n *node, version uint64) error {
	delete(n.deleted, version)
	delete(n.value, version)
	for _, child := range n.children {
		if err := rollback(child, version); err != nil {
			return err
		}
	}
	return nil
}

func parseKey(key string) ([]string, error) {
	parts := strings.Split(key, PathSeparator)

	partsPattern, err := regexp.Compile(keyPartPattern)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		if !partsPattern.MatchString(part) {
			return nil, errors.New(ErrorIllegalKey)
		}
	}

	if len(parts) == 0 {
		return nil, errors.New(ErrorNoRootSpecified)
	}

	return parts, nil
}

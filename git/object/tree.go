package object

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/0x414a/GitHack/utils"
)

type TreeItem interface {
	Name() string
	Hash() string
	Type() Type
	String() string
}

type treeItem struct {
	name string
	hash string
	t    Type
}

func (t treeItem) Name() string {
	return t.name
}

func (t treeItem) Hash() string {
	return t.hash
}

func (t treeItem) Type() Type {
	return t.t
}

func (t treeItem) String() string {
	return fmt.Sprintf("Name: %s\nHash: %s\nType: %s\n", t.Name(), t.Hash(), t.Type())
}

type TreeObject interface {
	Object
	Items() []TreeItem
}

type treeObject struct {
	object
	items []treeItem
}

func (o *treeObject) Items() []TreeItem {
	items := make([]TreeItem, len(o.items))
	for i := range o.items {
		items[i] = &o.items[i]
	}
	return items
}

func (o *treeObject) String() string {
	str := "[TreeObject]\n"
	for _, i := range o.items {
		str += i.String() + "\n"
	}
	return str
}

func NewTreeObject(o object) *treeObject {
	obj := &treeObject{
		object: o,
		items:  []treeItem{},
	}

	buffer := bytes.NewBuffer(o.data)

	for {
		t, err := utils.ReadStringWithoutDelimiter(buffer, byte(' '))
		if err != nil {
			break
		}

		i := treeItem{}
		switch t {
		case "100644":
			i.t = TypeSimpleBlob
		case "100755":
			i.t = TypeExecutableBlob
		case "120000":
			i.t = TypeSymbolLinkBlob
		case "40000":
			i.t = TypeTree
		}
		i.name, _ = utils.ReadStringWithoutDelimiter(buffer, byte(0))
		i.hash = hex.EncodeToString(buffer.Next(20))
		obj.items = append(obj.items, i)
	}
	return obj
}

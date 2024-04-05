package object

import (
	"bytes"
	"fmt"
	"github.com/0x414a/GitHack/utils"
)

type CommitObject interface {
	Object
	Tree() string
	Parent() string
	Author() string
	Committer() string
	GPGSignature() string
	Message() string
}

type commitObject struct {
	object
	tree      string
	parent    string
	author    string
	committer string
	signature string
	message   string
}

func (o *commitObject) Tree() string {
	return o.tree
}

func (o *commitObject) Parent() string {
	return o.parent
}

func (o *commitObject) Author() string {
	return o.author
}

func (o *commitObject) Committer() string {
	return o.committer
}

func (o *commitObject) GPGSignature() string {
	return o.signature
}

func (o *commitObject) Message() string {
	return o.message
}

func (o *commitObject) String() string {
	return fmt.Sprintf("Tree: %s\nParent: %s\nAuthor: %s\nCommitter: %s\nGPGSignature: %s\nMessage: \n%s\n",
		o.Tree(),
		o.Parent(),
		o.Author(),
		o.Committer(),
		o.GPGSignature(),
		o.Message(),
	)
}

func NewCommitObject(o object) *commitObject {
	obj := &commitObject{
		object:    o,
		tree:      "",
		parent:    "",
		author:    "",
		committer: "",
		signature: "NOT_IMPLEMENTED",
		message:   "",
	}

	buffer := bytes.NewBuffer(o.data)

	for {
		if obj.tree != "" && obj.parent != "" && obj.author != "" && obj.committer != "" {
			obj.message = string(buffer.Bytes())
			break
		}

		t, _ := utils.ReadStringWithoutDelimiter(buffer, byte(' '))
		switch t {
		case "tree":
			obj.tree, _ = utils.ReadStringWithoutDelimiter(buffer, byte(10))
		case "parent":
			obj.parent, _ = utils.ReadStringWithoutDelimiter(buffer, byte(10))
		case "author":
			obj.author, _ = utils.ReadStringWithoutDelimiter(buffer, byte(10))
		case "committer":
			obj.committer, _ = utils.ReadStringWithoutDelimiter(buffer, byte(10))
			//case "gpgsig":
			//	obj.signature, _ = utils.ReadStringWithoutDelimiter(buffer, byte(10))
		}
	}
	return obj
}

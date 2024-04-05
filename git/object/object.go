package object

import (
	"bytes"
	"github.com/0x414a/GitHack/utils"

	"strconv"
)

type Object interface {
	IsBlob() bool
	IsTree() bool
	IsCommit() bool

	Length() int
	String() string
}

type Type int

const (
	TypeSimpleBlob     Type = 100644
	TypeExecutableBlob Type = 100755
	TypeSymbolLinkBlob Type = 120000
	TypeBlob           Type = 1
	TypeTree           Type = 040000
	TypeCommit         Type = 2
)

func (t Type) String() string {
	switch t {
	case 100644:
		return "blob(text)"
	case 100775:
		return "blob(executable)"
	case 120000:
		return "blob(symbol)"
	case 1:
		return "blob"
	case 040000:
		return "tree"
	case 2:
		return "commit"
	}
	return "unknown"
}

type object struct {
	o      Type
	length int
	data   []byte
}

func (o object) IsBlob() bool {
	return o.o == TypeSimpleBlob || o.o == TypeExecutableBlob || o.o == TypeSymbolLinkBlob || o.o == TypeBlob
}

func (o object) IsTree() bool {
	return o.o == TypeTree
}

func (o object) IsCommit() bool {
	return o.o == TypeCommit
}

func (o object) Length() int {
	return o.length
}

func (o object) String() string {
	return string(o.data) + "\n"
}

func NewObject(obj []byte) Object {
	buffer := bytes.NewBuffer(obj)
	objType, _ := utils.ReadStringWithoutDelimiter(buffer, byte(' '))
	size, _ := utils.ReadStringWithoutDelimiter(buffer, byte(0))
	length, _ := strconv.Atoi(size)
	data := buffer.Bytes()

	o := object{
		length: length,
		data:   data,
	}

	switch objType {
	case "commit":
		o.o = TypeCommit
		return NewCommitObject(o)
	case "tree":
		o.o = TypeTree
		return NewTreeObject(o)
	case "blob":
		o.o = TypeBlob
	}
	return &o
}

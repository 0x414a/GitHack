package git

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/0x414a/GitHack/utils"
	"io"
)

func (f *fetcher) GenerateObjectPath(hash string) string {
	if len(hash) >= 2 {
		return fmt.Sprintf("%sobjects/%s/%s", f.targetUrl, hash[0:2], hash[2:])
	} else {
		return ""
	}

}

func (f *fetcher) FetchObject(hash string) []byte {
	bin := utils.GetBinary(f.GenerateObjectPath(hash))
	r, err := zlib.NewReader(bytes.NewReader(bin))
	if err != nil {
		return nil
	}

	defer r.Close()
	buffer := bytes.Buffer{}
	_, _ = io.Copy(&buffer, r)
	return buffer.Bytes()
}

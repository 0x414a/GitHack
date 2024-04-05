package git

import (
	"github.com/0x414a/GitHack/utils"
)

func (f *fetcher) FetchRef(path string) (hash string) {
	hash = utils.GetText(f.targetUrl + path)

	if len(hash) != 40 {
		hash = ""
	}
	return
}

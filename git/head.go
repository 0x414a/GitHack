package git

import (
	"github.com/0x414a/GitHack/utils"
	"regexp"
)

// FetchHead returns ref name of current head
func (f *fetcher) FetchHead() (head, raw string) {
	raw = utils.GetText(f.targetUrl + HeadPath)
	r := regexp.MustCompile(HeadRegexp)

	match := r.FindStringSubmatch(raw)
	if len(match) != 2 {
		raw = ""
	} else {
		head = match[1]
	}
	return
}

func (f *fetcher) FetchOriginHead() (hash string) {
	return f.FetchRef(OriginHeadPath)
}

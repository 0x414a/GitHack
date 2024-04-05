package git

type Fetcher interface {
	FetchHead() (head, raw string)
	FetchOriginHead() (hash string)
	FetchRef(path string) (hash string)
	FetchObject(hash string) (out []byte)
}

type fetcher struct {
	targetUrl string
}

func NewFetcher(targetUrl string) Fetcher {
	return &fetcher{
		targetUrl: targetUrl,
	}
}

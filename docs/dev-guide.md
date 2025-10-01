### Prerequisites

1.  Ubuntu (Server or Desktop) operating system -- other similar systems and OS X might work, but aren't guaranteed to...
2.  GoLang installed

## How to Publish to Go Modules

https://go.dev/doc/modules/publishing

```
git tag v0.1.1 
git push origin v0.1.1
GOPROXY=proxy.golang.org go list -m github.com/exlinc/golang-utils@v0.1.1`
```

### Get the project

`go get github.com/exlinc/golang-utils/...`
`go get -u github.com/exlinc/golang-utils@v0.1.1`

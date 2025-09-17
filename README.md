# go-serper

[![Go Reference](https://pkg.go.dev/badge/github.com/iamwavecut/go-serper.svg)](https://pkg.go.dev/github.com/iamwavecut/go-serper) [![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/go-serper)](https://goreportcard.com/report/github.com/iamwavecut/go-serper) [![CI](https://github.com/iamwavecut/go-serper/actions/workflows/ci.yml/badge.svg)](https://github.com/iamwavecut/go-serper/actions/workflows/ci.yml) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Minimal serper.dev client with no external dependencies.

## Install

```bash
go get github.com/iamwavecut/go-serper
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	serper "github.com/iamwavecut/go-serper"
)

func main() {
	client := serper.NewClient("SERPER_API_KEY")
	resp, err := client.Search(context.Background(), serper.SearchRequest{Query: "golang"})
	if err != nil {
		panic(err)
	}
	for _, r := range resp.Results {
		fmt.Println(r.Title, r.URL)
	}
}
```

Optional configuration:

- `WithBaseURL(url)`
- `WithTimeouts(request, total)`
- `WithRetryConfig(count, baseDelay)`
- `WithHTTPClient(*http.Client)`
- `WithLogger(Logger)`

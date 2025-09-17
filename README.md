# go-serper

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

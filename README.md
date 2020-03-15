# logger

[![Build Status](https://travis-ci.org/opentogo/logger.svg?branch=master)](https://travis-ci.org/opentogo/logger)
[![GoDoc](https://godoc.org/github.com/opentogo/logger?status.png)](https://godoc.org/github.com/opentogo/logger)
[![codecov](https://codecov.io/gh/opentogo/logger/branch/master/graph/badge.svg)](https://codecov.io/gh/opentogo/logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/opentogo/logger)](https://goreportcard.com/report/github.com/opentogo/logger)
[![Open Source Helpers](https://www.codetriage.com/opentogo/logger/badges/users.svg)](https://www.codetriage.com/opentogo/logger)

A `log` wrapper that implements a `http.HandlerFunc` based in [Apache combined log format](https://httpd.apache.org/docs/2.2/logs.html#combined).

## Installation

```bash
go get github.com/opentogo/logger
```

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/opentogo/logger"
)

func main() {
	var (
		log    = logger.NewLogger(os.Stdout, "[app] ", 0)
		mux    = http.NewServeMux()
		server = &http.Server{
			Addr:    ":3000",
			Handler: log.Handler(mux),
		}
	)

	mux.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
```

And the output is:

```bash
[app] ::1 - - [12/Oct/2019:11:57:02 +0000] "GET / HTTP/1.1" 404 19 "-" "curl/7.54.0" 0.0000
[app] ::1 - - [12/Oct/2019:11:57:07 +0000] "GET /hello-world HTTP/1.1" 200 31 "-" "curl/7.54.0" 0.0000
```

## Contributors

- [rogeriozambon](https://github.com/rogeriozambon) Rogério Zambon - creator, maintainer

Filter Transformer
====================================

A package that transforms filtering conditions from requests to Elasticsearch, SQL, and other backends.

[![Build Status](https://www.travis-ci.com/wernerdweight/filter-transformer-go.svg?branch=master)](https://www.travis-ci.com/wernerdweight/filter-transformer-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/wernerdweight/filter-transformer-go)](https://goreportcard.com/report/github.com/wernerdweight/filter-transformer-go)
[![GoDoc](https://godoc.org/github.com/wernerdweight/filter-transformer-go?status.svg)](https://godoc.org/github.com/wernerdweight/filter-transformer-go)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/wernerdweight/filter-transformer-go)


Installation
------------

### 1. Installation

```bash
go get github.com/wernerdweight/filter-transformer-go
```

Configuration
------------

TODO:

Usage
------------

TODO:

### Errors

The following errors can occur (you can check for specific code since different errors have different severity):

```go
var FilterTransformerErrorCodes = map[FilterTransformerErrorCode]string{
    Unknown:                   "unknown error",
    // TODO: add other errors
}
```

License
-------
This package is under the MIT license. See the complete license in the root directory of the package.

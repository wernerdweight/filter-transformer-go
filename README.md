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

This package doesn't require any configuration.

Usage
------------

### Filter structure

The filter structure needs to be in the following format:

```js
{
  filter: {
    logic: "and|or",
    conditions: [
      {
        // regular filter
        field: "fieldName",
        operator: "eq|neq|...",
        value: "filtering value"
      },
      {
        // nested filter
        logic: "and|or",
        conditions: [ /*...*/ ]
      },
      ...
    ]
  }
}
```

Filters can be nested and support `AND` and `OR` logic.

The following operators are supported:

* **eq** - is equal to (equivalent of `=` in SQL),
* **neq** - is not equal to (equivalent of `!=` in SQL),
* **gt** - is greater than (equivalent of `>` in SQL),
* **gte** - is greater than or equal (equivalent of `>=` in SQL),
* **gten** - is greater than or equal or NULL (equivalent of `>= OR IS NULL` in SQL),
* **lt** - is lower than (equivalent of `<` in SQL),
* **lte** - is lower than or equal (equivalent of `<=` in SQL),
* **begins** - begins with (equivalent of `LIKE '...%'` in SQL),
* **contains** - contains (equivalent of `LIKE '%...%'` in SQL),
* **not-contains** - does not contain (equivalent of `NOT LIKE '%...%'` in SQL),
* **ends** - ends with (equivalent of `LIKE '%...'` in SQL),
* **null** - is null (equivalent of `IS NULL` in SQL),
* **not-null** - is not null (equivalent of `IS NOT NULL` in SQL),
* **empty** - is empty (equivalent of `IS NULL OR = ''` in SQL),
* **not**-empty - is not empty (equivalent of `IS NOT NULL AND != ''` in SQL),
* **in** - is contained in (equivalent of `IN` in SQL).

**Currently, only JSON input is supported.** FormData input will be supported in the future, with the same structure. The input could then look like this:

```http request
POST /some/path HTTP/1.1
Content-Type: multipart/form-data; boundary=---some-boundary
Host: your-api-host.com
Content-Length: 123

-----some-boundary
Content-Disposition: form-data; name="filter[logic]"
and
-----some-boundary
Content-Disposition: form-data; name="filter[conditions][0][field]"
key
-----some-boundary
Content-Disposition: form-data; name="filter[conditions][0][operator]"
eq
-----some-boundary
Content-Disposition: form-data; name="filter[conditions][0][value]"
val
-----some-boundary--
```

### Basic usage

```go
import (
    "github.com/wernerdweight/filter-transformer-go"
)

// this would normally come from a request
var jsonInput, _ = contract.NewInputOutputType(
    []byte(`{"logic": "and", "conditions": [{"field": "key", "operator": "eq", "value": "val"}]}`),
    &input.JsonInput{},
)

func main() {
    // transform from JSON to Elasticsearch
    ft := NewJsonToElasticFilterTransformer()
    output, err := ft.Transform(jsonInput)
    if err != nil {
        // handle error (see below)
    }
    // output = {"query": {"bool": {"must": [{"term": {"key": "val"}}]}}}
	
    // transform from JSON to SQL
    ft = NewJsonToSQLFilterTransformer()
    output, err = ft.Transform(jsonInput)
    if err != nil {
        // handle error (see below)
    }
    // output = Query: "key" = $1, Params: ["val"]
	
    // set up transformer with custom input and output
    it := input.JsonInputTransformer{} // this can be a custom input transformer
    ot := output.ElasticOutputTransformer{} // this can be a custom output transformer
    ft = NewFilterTransformer[[]byte, map[string]any, *input.JsonInput, *output.ElasticOutput](&it, &ot)
    output, err = ft.Transform(jsonInput)
    if err != nil {
        // handle error (see below)
    }
    // output = ...
}
```

### Supported input and output types

**The following input types are supported:**

* **JSON** - `input.JsonInput` - input is `[]byte`,
* ~~**FormData** - `input.FormDataInput` - input is `map[string]any`.~~ (not yet supported)

**The following output types are supported:**

* **Elasticsearch** - `output.ElasticOutput` - output is `map[string]any`,
* **SQL** - `output.SQLOutput` - output is `struct { Query string; Params []any }`.

For SQL, the output is a struct with a query and parameters. The query is a string with placeholders for parameters (e.g. `$1`, `$2`, ...). The parameters are an array of values that are used to replace the placeholders in the query.

### Errors

The following errors can occur (you can check for specific code since different errors have different severity):

```go
var ErrorCodes = map[ErrorCode]string{
    Unknown:                   "unknown error",
    UnreadableInputData:       "unreadable input data",
    InvalidInputDataStructure: "invalid input data structure",
    InvalidFiltersStructure:   "invalid filters structure",
    NonWriteableOutputData:    "can't write output data",
}
```

### Custom transformers

You can create custom input and output transformers by implementing the `InputTransformerInterface` and `OutputTransformerInterface` interfaces.

```go
type InputTransformerInterface[T any, IOT InputOutputInterface[T]] interface {
    Transform(input IOT) (Filters, *Error)
}

type OutputTransformerInterface[T any, IOT InputOutputInterface[T]] interface {
    Transform(input Filters) (IOT, *Error)
}
```

### Known issues, limitations and missing features

* **FormData input** - not yet supported.
* **Validation** - input validation is not yet supported - therefore, the produced output might not be usable and you should handle such cases in your application. This package doesn't have any information about your fields, their types and permissions. This will be supported in the future:
  * **Field validation** - input fields validation is not yet supported - once done, you'll be able to validate the fields (e.g. whether they exist, can be filtered, can be used with a specific operator, etc.) based on your own validation logic.
  * **Value validation** - input validation is not yet supported - once done, you'll be able to validate the input values (condition values) before transforming based on your own validation logic.
* **SQL output** - the GetDataString method is not not safe for use in production, it's intended for debugging purposes only (there's a log line printed in the method to make this explicit).

License
-------
This package is under the MIT license. See the complete license in the root directory of the package.

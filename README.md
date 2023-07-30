# My very own JSON parser

This package attempts to clean-room something close to the Go Standard Library's `encoding/json`
package. It is inspired by [this coding
challenge](https://codingchallenges.substack.com/p/coding-challenge-2).
Currently, it supports a minimal subset of functionality.

* Decodes simple objects (i.e., no nested objects) with string values.
* Decodes into interface types (i.e., no structs with tags).

## Current benchmarks vs Stdlib

```
✗ go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/yardbirdsax/json
BenchmarkDecode-8                 433648              2634 ns/op
BenchmarkDecodeNative-8           566842              1805 ns/op
```

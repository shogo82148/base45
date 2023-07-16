# base45

[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/base45.svg)](https://pkg.go.dev/github.com/shogo82148/base45)
[![test](https://github.com/shogo82148/base45/actions/workflows/test.yml/badge.svg)](https://github.com/shogo82148/base45/actions/workflows/test.yml)

Package base45 implements base45 encoding as specified by [RFC 9285](https://www.rfc-editor.org/rfc/rfc9285.html).

## SYNOPSIS

```go
msg := "Hello, 世界"
encoded := base45.EncodeToString([]byte(msg))
fmt.Println(encoded)
decoded, err := base45.DecodeString(encoded)
if err != nil {
    fmt.Println("decode error:", err)
    return
}
fmt.Println(string(decoded))
// Output:
// %69 VDK2E5744FNKCT53
// Hello, 世界
```

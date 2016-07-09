package httputil2

import (
    "net/http"
    "strings"
)

func IsMobile(header http.Header) bool {
    a := header["User-Agent"]
    if len(a) == 0 {
        return false
    }
    return (strings.Index(a[0], "iPhone") >= 0)
}

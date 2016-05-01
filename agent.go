package httputil2

import "net/http"

func SetupAgentChromeMac(request *http.Request) {
    request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.75 Safari/537.36")
}

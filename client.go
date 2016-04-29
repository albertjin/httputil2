package httputil2

import (
    "bytes"
    "errors"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "strings"

    "github.com/albertjin/net2"
    "github.com/albertjin/goquery"
    "golang.org/x/text/transform"
    "golang.org/x/text/encoding"
)

// Make http request. If non-nil form is passed, POST method is used; otherwise, GET is used.
func Request(client *http.Client, link, cookie string, form url.Values, prepare func(request *http.Request)) (response *http.Response, err error) {
    var request *http.Request
    if form == nil {
        request, err = http.NewRequest("GET", link, nil)
    } else {
        request, err = http.NewRequest("POST", link, strings.NewReader(form.Encode()))
        if request != nil {
            request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        }
    }
    if err != nil {
        return
    }
    if len(cookie) > 0 {
        request.Header.Set("Cookie", cookie)
    }
    if prepare != nil {
        prepare(request)
    }

    response, err = client.Do(request)
    return
}

// The stream returned is in UTF-8 encoding, which is considered as text by Go. It fails early when the text type checking is not expected.
func TextStreamFromResponse(response *http.Response, check func(textType string) error) (stream io.Reader, textType string, err error) {
    charset, textType, err := CharsetFromContentType(response.Header.Get("Content-Type"))
    if err != nil {
        return
    }

    var data []byte
    var encoding encoding.Encoding
    if (len(charset) == 0) && (textType == "html") {
        data, err = ioutil.ReadAll(response.Body)
        if err != nil {
            return
        }
        charset = DetectCharset(data)
        if encoding, err = GetEncoding(charset); err != nil {
            return
        }
        if encoding == nil {
            stream = bytes.NewReader(data)
            return
        }

        stream = transform.NewReader(bytes.NewReader(data), encoding.NewDecoder())
    } else {
        if encoding, err = GetEncoding(charset); err != nil {
            return
        }
        if encoding != nil {
            stream = transform.NewReader(response.Body, encoding.NewDecoder())
        } else {
            stream = response.Body
        }
    }
    return
}

// There might be a shortcut when stream is not required, comparing to StreamFromResponse().
func TextFromResponse(response *http.Response) (text, textType string, err error) {
    charset, textType, err := CharsetFromContentType(response.Header.Get("Content-Type"))
    if err != nil {
        return
    }

    var data []byte
    var stream io.Reader
    var encoding encoding.Encoding
    if (len(charset) == 0) && (textType == "html") {
        data, err = ioutil.ReadAll(response.Body)
        if err != nil {
            return
        }
        charset = DetectCharset(data)
        if encoding, err = GetEncoding(charset); err != nil {
            return
        }
        // No encoding, it assumed as UTF-8.
        if encoding == nil {
            text = string(data)
            return
        }

        stream = transform.NewReader(bytes.NewReader(data), encoding.NewDecoder())
    } else {
        if encoding, err = GetEncoding(charset); err != nil {
            return
        }
        if encoding != nil {
            stream = transform.NewReader(response.Body, encoding.NewDecoder())
        } else {
            stream = response.Body
        }
    }

    if data, err = ioutil.ReadAll(stream); err == nil {
        text = string(data)
    }
    return
}

// Get binary data.
func GetData(client *http.Client, link, cookie string, form url.Values, prepare func(request *http.Request)) (data []byte, response *http.Response, err error) {
    response, err = Request(client, link, cookie, form, prepare)
    if err != nil {
        return
    }

    data, err = ioutil.ReadAll(response.Body)
    return
}

// Get text.
func GetText(client *http.Client, link, cookie string, form url.Values, parpare func(request *http.Request)) (text, textType string, response *http.Response, err error) {
    response, err = Request(client, link, cookie, form, parpare)
    if err != nil {
        return
    }
    defer response.Body.Close()

    text, textType, err = TextFromResponse(response)
    return
}

// Get goquery document.
func GetDocument(client *http.Client, link, cookie string, form url.Values, prepare func(request *http.Request)) (document *goquery.Document, response *http.Response, err error) {
    response, err = Request(client, link, cookie, form, prepare)
    if err != nil {
        return
    }
    defer response.Body.Close()

    stream, _, err := TextStreamFromResponse(response, func(textType string) (err error) {
        if textType != "html" {
            err = errors.New("The type of response is [" + textType + "] but not html.")
        }
        return
    })
    if err != nil {
        return
    }

    document, err = goquery.NewDocumentFromReader(stream)
    return
}

// New http.Client with timeouts.
func NewClient(timeout *net2.Timeout) *http.Client {
    if timeout == nil {
        timeout = net2.DefaultTimeout
    }
    return &http.Client {
        Transport: &http.Transport {
            Dial: func(network, addr string) (net.Conn, error) {
                return net2.Dial(network, addr, timeout)
            },
        },
    }
}

var DefaultClient = NewClient(nil)

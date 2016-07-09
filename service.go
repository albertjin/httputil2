package httputil2

import (
    "compress/gzip"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "sync"
    "time"

    "github.com/albertjin/ajax2"
)

type Service struct {
    listener net.Listener
}

type Control interface {
    Use()
    Unuse()
}

func Serve(server *http.Server, control Control) (service *Service, err error) {
    service = &Service{nil}
    addr := server.Addr
    if addr == "" {
       addr = ":http"
    }

    service.listener, err = net.Listen("tcp", addr)
    if err == nil {
        control.Use(); go func() {
            defer control.Unuse()
            server.Serve(service.listener)
        }()
    }

    return
}

func (t *Service) Close() {
    if t.listener != nil {
        t.listener.Close()
        t.listener = nil
    }
}

func GetHeadText(r *http.Request, name string) (text string) {
    ss := r.Header[name]
    if len(ss) > 0 {
        text = ss[0]
    }
    return
}

func IfModified(r *http.Request, latest time.Time) (modified bool, err error) {
    a := GetHeadText(r, "If-Modified-Since")
    if len(a) > 0 {
        var b time.Time
        b, err = time.Parse(time.RFC1123, a)
        if err == nil {
            modified = latest.After(b)
            return
        }
    }
    modified = true
    return
}

func SetContentEncodingGzip(w http.ResponseWriter) {
    w.Header().Set("Content-Encoding", "gzip")
}

func SetContentTypeJson(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func SetContentLength(w http.ResponseWriter, length int) {
    w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
}

func WriteBytes(w http.ResponseWriter, r *http.Request, data []byte, compressed, autoCompress bool, code int) (err error) {
    head := (r.Method == "HEAD")
    if length := len(data); compressed || !autoCompress || (length < 1024) {
        if compressed {
            SetContentEncodingGzip(w)
        }
        SetContentLength(w, length)
        w.WriteHeader(code)
        if !head {
            _, err = w.Write(data)
        }
    } else {
        SetContentEncodingGzip(w)
        w.WriteHeader(code)
        if !head {
            v := gzip.NewWriter(w); defer v.Close()
            _, err = v.Write(data)
        }
    }
    return
}

func WriteJson(w http.ResponseWriter, r *http.Request, object interface{}) (err error) {
    var data []byte
    switch t := object.(type) {
    case []byte:
        data = t
    case nil:
        data = ajax2.Null
    default:
        data, err = json.Marshal(object)
        if err != nil {
            data = ajax2.StatusInternalServerError
        }
    }
    SetContentTypeJson(w)
    err = WriteBytes(w, r, data, false, true, http.StatusOK)
    return
}

const httpHead = "HTTP/1.1 %s\r\n"+
    "Content-Length: %d\r\n"+
    "Content-Type: %s\r\n"+
    "Connection: close\r\n\r\n"

func WriteResponse(w io.Writer, data []byte, status, contentType string) {
    w.Write([]byte(fmt.Sprintf(httpHead, status, len(data), contentType)))
    w.Write(data)
}

type HttpService struct {
    listener net.Listener
}

func NewHttpService(context interface { Acquire(); Release() }, server *http.Server) (s *HttpService, err error) {
    s = &HttpService{}
    addr := server.Addr
    if addr == "" {
        addr = ":http"
    }

    s.listener, err = net.Listen("tcp", addr)
    if err != nil {
        return
    }

    var w sync.WaitGroup
    w.Add(1)
    context.Acquire(); go func() {
        defer context.Release()
        l := s.listener
        w.Done()

        server.Serve(l)
    }()

    return
}

func (s *HttpService) End() {
    if s.listener != nil {
        s.listener.Close()
        s.listener = nil
    }
}

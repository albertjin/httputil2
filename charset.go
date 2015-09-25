package httputil2

import (
    "bytes"
    "errors"
    "regexp"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "golang.org/x/text/encoding"
    "golang.org/x/text/encoding/charmap"
    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/encoding/korean"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/encoding/traditionalchinese"
)

func CharsetFromContentType(contentType string) (charset string, textType string, err error) {
    contentType = strings.ToLower(contentType)

    re := regexp.MustCompile(` *([a-z0-9\-.]+)/([a-z0-9\-.]+) *;?`)
    if ss := re.FindStringSubmatch(contentType); len(ss) == 3 {
        switch ss[1] {
        case "application":
            switch ss[2] {
            case "json", "javascript":
                textType = ss[2]
            }
        case "text":
            textType = ss[2]
        }
    }

    if len(textType) == 0 {
        err = errors.New("It's not text.\n" + contentType)
        return
    }
    re = regexp.MustCompile(`; *charset=([a-z0-9\-_]+)`)
    if ss := re.FindStringSubmatch(contentType); len(ss) == 2 {
        charset = ss[1]
    }
    return
}

// Detect charset from byte array when charset base is based on ascii codes and compatible with utf-8. Otherwise it does not work.
func DetectCharset(source []byte) (charset string) {
    data := source
    const peek = 256
    if len(source) > peek {
        data = source[:peek]
    }

    for i := 0; ; i++ {
        if doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data)); err == nil {
            if ss := doc.Find("meta[charset]"); ss.Length() > 0 {
                charset, _ = ss.Attr("charset")
                charset = strings.ToLower(charset)
                return
            }

            ss := doc.Find("meta[http-equiv]")
            for i, l := 0, ss.Length(); i < l; i++ {
                s := ss.Eq(i)
                if h, _ := s.Attr("http-equiv"); strings.ToLower(h) == "content-type" {
                    contentType, _ := s.Attr("content")
                    charset, _, _ = CharsetFromContentType(contentType)
                    return
                }
            }
        }
        if (i == 1) || (len(source) <= peek) {
            break
        }
        re := regexp.MustCompile(`</ *[hH][eE][aA][dD] *>`)
        if loc := re.FindIndex(source); (len(loc) == 2) && (loc[0] > peek) {
            data = source[:loc[0]]
        } else {
            break
        }
    }
    return
}

func GetEncoding(charset string) (encoding encoding.Encoding, err error) {
    // http://www.iana.org/assignments/character-sets/character-sets.xhtml
    switch charset {
    case "", "ascii", "utf-8":
    case "gbk":
        encoding = simplifiedchinese.GBK
    case "gb18030", "gb-18030":
        encoding = simplifiedchinese.GB18030
    case "gb2312", "gb-2312":
        encoding = simplifiedchinese.HZGB2312
    case "big5":
        encoding = traditionalchinese.Big5
    case "euckr":
        encoding = korean.EUCKR
    case "shiftjis", "shift_jis":
        encoding = japanese.ShiftJIS
    case "iso-8859-1", "windows-1252": // windows-1252 is the superset of iso-8859-1.
        encoding = charmap.Windows1252
    default:
        err = errors.New("The charset [" + charset + "] is not supported.")
        log(critical, err.Error(), stack)
    }
    return
}
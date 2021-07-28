package arc

import (
  "io"
  "bytes"
  "compress/gzip"
)

func gzdeflate(str string) string {
    var b bytes.Buffer

    w, _ := gzip.NewWriterLevel(&b, 9)
    w.Write([]byte(str))
    w.Close()
    return b.String()
}

func gzinflate(str string) string {
    b := bytes.NewReader([]byte(str))
    r, _ := gzip.NewReader(b)
    bb2 := new(bytes.Buffer)
    _, _ = io.Copy(bb2, r)
    r.Close()
    byts := bb2.Bytes()
    return string(byts)
}

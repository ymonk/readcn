package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"
)

func decodeBody(r *http.Request, v interface{}) error {
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(v)
}


func encodeBody(w http.ResponseWriter, v interface{}) error {
    return json.NewEncoder(w).Encode(v)
}


func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}, cache bool) {
    if cache && data != nil {
        buf := new(bytes.Buffer)
        json.NewEncoder(buf).Encode(data)
        webcache.store(r.URL.String(), buf.Bytes())
        respondRaw(w, status, buf.Bytes())
    } else {
        w.WriteHeader(status)
        if data != nil {
            encodeBody(w, data)
        }
    }
}

func respondRaw(w http.ResponseWriter, status int, bytes []byte) {
    w.WriteHeader(status)
    w.Write(bytes)
}

func respondErr(w http.ResponseWriter, r *http.Request, status int, arg ...interface{}) {
    respond(w, r, status, map[string]interface{}{
        "error": map[string]interface{}{
            "message": fmt.Sprint(arg...),
        },
    }, false)
}

func respondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
    respondErr(w, r, status, http.StatusText(status))
}


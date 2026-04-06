package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%s %s\n", r.Method, r.URL.Path)
		n := buf.Len()
		io.Copy(buf, r.Body)
		if !isUTF8(buf.Bytes()[n:]) {
			m := buf.Len() - n
			buf.Truncate(n)
			fmt.Fprintf(buf, "<binary stream of length %d>", m)
		}
		log.Printf("%s", buf.Bytes())
		w.WriteHeader(http.StatusOK)
	})
	var addr string
	if addr = os.Getenv("LISTEN_ADDR"); addr == "" {
		addr = ":8080"
	}
	log.Printf("Listening on %s\n", addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}

func isUTF8(s []byte) bool {
	if len(s) == 0 {
		return true
	}
	if !utf8.Valid(s) {
		return false
	}
	for _, r := range string(s) {
		switch r {
		case '\n', '\r', '\t':
			continue
		}
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

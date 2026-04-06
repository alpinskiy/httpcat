package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%s %s\n", r.Method, r.URL.Path)
		for k, v := range r.Header {
			fmt.Fprintf(buf, "%s: %s\n", k, strings.Join(v, ", "))
		}
		fmt.Fprint(buf, "\n")
		bodyStart := buf.Len()
		bodyLen, _ := io.Copy(buf, r.Body)
		if bodyLen > 0 {
			if !isUTF8(buf.Bytes()[bodyStart:]) {
				buf.Truncate(bodyStart)
				fmt.Fprintf(buf, "<binary stream of length %d>", bodyLen)
			}
			fmt.Fprint(buf, "\n\n")
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

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%s %s\n", r.Method, r.URL.Path)
		io.Copy(buf, r.Body)
		buf.WriteByte('\n')
		buf.WriteTo(os.Stdout)
		w.WriteHeader(http.StatusOK)
	})
	var addr string
	if addr = os.Getenv("LISTEN_ADDR"); addr == "" {
		addr = ":8080"
	}
	fmt.Printf("Listening on %s\n", addr)
	fmt.Println(http.ListenAndServe(addr, nil))
}

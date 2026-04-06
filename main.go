package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		fmt.Fprintf(buf, "%s %s\n", r.Method, r.URL.Path)
		io.Copy(buf, r.Body)
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

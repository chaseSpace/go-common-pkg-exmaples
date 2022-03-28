package main

import "net/http"

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(202)
		writer.Write([]byte(`xxx`))
	})
	http.ListenAndServe(":1111", nil)
}

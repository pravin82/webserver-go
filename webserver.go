package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", servePage)
	http.ListenAndServe(":8080", nil)

}

func servePage(writer http.ResponseWriter, request *http.Request) {
	urlPath := request.URL.Path
	if urlPath == "" {
		urlPath = "index.html"
	}
	filePath := "www/" + urlPath
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	defer file.Close()
	writer.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(writer, file)
	if err != nil {
		http.Error(writer, "Error copying file contents", http.StatusInternalServerError)
		return
	}

}

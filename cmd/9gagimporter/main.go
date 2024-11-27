package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type importFile struct {
	ServerURL string `json:"server-url"`
	Items     []item `json:"items"`
}

type item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func main() {
	filePath := os.Args[1]
	if filePath == "" {
		//nolint: forbidigo // At the moment we authorize the logs for this tool
		fmt.Println("9gagimporter [jsonPath] [cookie]")
		os.Exit(1)
	}

	cookie := os.Args[2]
	if cookie == "" {
		//nolint: forbidigo // At the moment we authorize the logs for this tool
		fmt.Println("9gagimporter [jsonPath] [cookie]")
		os.Exit(1)
	}

	rawFile, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var file importFile

	err = json.Unmarshal(rawFile, &file)
	if err != nil {
		panic(err)
	}

	for _, item := range file.Items {
		if strings.HasSuffix(item.URL, ".mp4") {
			continue
		}

		//nolint: forbidigo // At the moment we authorize the logs for this tool
		fmt.Printf("%q -> %v\n", item.Title, item.URL)

		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		fw, err := w.CreateFormField("title")
		if err != nil {
			panic("failed to create the title form field")
		}

		_, err = fw.Write([]byte(item.Title))
		if err != nil {
			panic("failed to write the title form field value")
		}

		res, err := http.Get(item.URL)
		if err != nil {
			panic(err)
		}

		rawImg, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		res.Body.Close()

		fw, err = w.CreateFormFile("file", "file.unknown")
		if err != nil {
			panic(err)
		}

		fw.Write(rawImg)

		w.Close()

		req, err := http.NewRequest("POST", file.ServerURL+"/submit", &b)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", w.FormDataContentType())
		req.AddCookie(&http.Cookie{
			Name:  "session_token",
			Value: cookie,
		})

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		res.Body.Close()

		if res.StatusCode != 200 {
			panic(fmt.Sprintf("should return 302, have %d", res.StatusCode))
		}
	}
}

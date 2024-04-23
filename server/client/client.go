package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func NewClient() {
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	resp, err := client.Get("https://google.com/info")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
    token := ""
    url := "https://graph.facebook.com/v19.0/me/permissions?access_token=" + token
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}

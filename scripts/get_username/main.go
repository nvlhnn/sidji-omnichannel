package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// The ID and Token provided by the user
	igID := ""
	token := ""

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=username&access_token=%s", igID, token)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fmt.Printf("Instagram ID: %s\n", igID)
	fmt.Printf("Username: %v\n", result["username"])
	fmt.Printf("Full Response: %s\n", string(body))
}

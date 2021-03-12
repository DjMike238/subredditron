package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/NicoNex/echotron/v2"
)

func checkMsg(msg string) bool {
	return strings.Index(msg, "r/") != -1 && strings.Index(msg, "reddit.com") == -1
}

func extractMsg(message *echotron.Message) string {
	if message.Text != "" {
		return message.Text
	} else if message.Caption != "" {
		return message.Caption
	}

	return ""
}

func getStatus(url string) int {
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	return response.StatusCode
}

func getPreview(sub string) (string, string, string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://reddit.com/%s/about.json", sub), nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "Subredditron/2.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		log.Println("too many requests to reddit servers: try again later")
		return sub, "", ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var about About
	err = json.Unmarshal(body, &about)
	if err != nil {
		log.Println(err)
	}

	title := fmt.Sprintf("%s · %s", about.Data.Title, about.Data.Name)
	var thumb string

	if about.Data.Icon != "" {
		thumb = about.Data.Icon
	} else if about.Data.Banner != "" {
		thumb = about.Data.Banner
	}

	return title, about.Data.Description, thumb
}

func getSub(msg string) string {
	re := regexp.MustCompile(`(^|[ /])r\/[a-zA-Z_0-9]*`)
	sub := re.FindString(msg)
	var url string

	// Check if the matched string is longer than the minimum length for a subreddit
	// name (which is 3) and shorter than the maximum length for a subreddit name
	// (which is 21), both also counting "r/" or "*r/", where * is a character
	// that can be a space (" ") or a slash ("/").
	if len(sub) >= 5 && len(sub) <= 23 && sub[:2] == "r/" {
		url = fmt.Sprintf("https://www.reddit.com/%s", sub)
	} else if len(sub) >= 6 && len(sub) <= 24 && sub[1:3] == "r/" {
		url = fmt.Sprintf("https://www.reddit.com/%s", sub[1:])
	}

	return url
}

package main

import (
	"fmt"
)

// About is a custom struct for the about.json returned by the Reddit API.
type About struct {
	Data *Data `json:"data"`
}

// Data is a custom struct for the data object contained in the about.json returned by the Reddit API.
type Data struct {
	Title         string `json:"title"`
	DisplayName   string `json:"display_name_prefixed"`
	Description   string `json:"public_description,omitempty"`
	HeaderTitle   string `json:"header_title,omitempty"`
	Icon          string `json:"icon_img,omitempty"`
	Banner        string `json:"header_img,omitempty"`
}

func getThumb(data *Data) string {
	if data.Icon != "" {
		return data.Icon
	} else if data.Banner != "" {
		return data.Banner
	}

	return ""
}

func getTitle(data *Data) string {
	if data.Title != "" && data.Title != data.DisplayName {
		return fmt.Sprintf("%s · %s", data.Title, data.DisplayName)
	}
	
	return data.DisplayName
}

func getDesc(data *Data) string {
	if data.Description != "" {
		return data.Description
	} else if data.HeaderTitle != "" {
		return data.HeaderTitle
	}

	return ""
}

func getName(data *Data) string {
	return data.DisplayName
}

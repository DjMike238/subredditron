package main

// About is a custom struct for the about.json returned by the Reddit API.
type About struct {
	Data *Data `json:"data"`
}

// Data is a custom struct for the data object contained in the about.json returned by the Reddit API.
type Data struct {
	Title         string `json:"title"`
	Name          string `json:"display_name_prefixed"`
	Description   string `json:"public_description,omitempty"`
	Icon          string `json:"icon_img,omitempty"`
	CommunityIcon string `json:"community_icon,omitempty"`
	Banner        string `json:"header_img,omitempty"`
}

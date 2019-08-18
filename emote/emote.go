package emote

// Emote contains metadata about an emote
type Emote struct {
	Name             string `json:"name"`
	Filename         string `json:"filename"`
	OriginalFilename string `json:"originalFilename"`
	Credit           string `json:"credit"`
}

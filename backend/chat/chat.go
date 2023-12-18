package chat

type Chat interface {
	GetType() string
}

type TextChat struct {
	Type     string `json:"type"`
	Role     string `json:"role"`
	Text     string `json:"data"`
	DateTime string `json:"dateTime"`
}

func (tc *TextChat) GetType() string {
	return "text"
}

type ImageChat struct {
	Type     string   `json:"type"`
	Role     string   `json:"role"`
	ImageId  string   `json:"data"`
	DateTime string   `json:"dateTime"`
	Tags     []string `json:"tags"`
}

func (imageChat *ImageChat) GetType() string {
	return "image"
}

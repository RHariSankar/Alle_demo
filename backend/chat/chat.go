package chat

type Chat interface {
	GetData() string
}

type TextChat struct {
	Type     string `json:"type"`
	Role     string `json:"role"`
	Data     string `json:"data"`
	DateTime string `json:"dateTime"`
}

func (tc *TextChat) GetData() string {
	return tc.Data
}

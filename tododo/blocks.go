package tododo

type Response struct {
	Blocks []*Block `json:"blocks"`
}

type Block struct {
	Type    string        `json:"type"`
	BText   *BlockText    `json:"text,omitempty"`
	BFields []*BlockField `json:"fields,omitempty"`
}

type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type BlockField struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewSectionTextBlock(textType string, text string) *Block {
	block := Block{}
	block.Type = "section"
	block.BText = &BlockText{Type: textType, Text: text}
	return &block
}

func NewSectionFieldsBlock(fields ...*BlockField) *Block {
	block := Block{}
	block.Type = "section"
	arr := make([]*BlockField, 0)
	for _, f := range fields {
		arr = append(arr, f)
	}
	block.BFields = arr
	return &block
}

func NewHeaderBlock(text string) *Block {
	block := Block{}
	block.Type = "header"
	block.BText = &BlockText{Type: "plain_text", Text: text}
	return &block
}

func NewDividerBlock() *Block {
	block := Block{}
	block.Type = "divider"
	return &block
}

func NewField(fieldType string, text string) *BlockField {
	field := BlockField{}
	field.Type = fieldType
	field.Text = text
	return &field
}

func NewResponse(blocks ...*Block) *Response {
	resp := Response{}
	arr := make([]*Block, 0)
	for _, b := range blocks {
		arr = append(arr, b)
	}
	resp.Blocks = arr
	return &resp
}

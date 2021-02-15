package tododo

// Block kit response. Can visualize the response with https://app.slack.com/block-kit-builder.
type Response struct {
	Blocks []*Block `json:"blocks"`
}

type Block struct {
	Type    string        `json:"type"`
	BText   *BlockText    `json:"text,omitempty"`
	BFields []*BlockField `json:"fields,omitempty"`
}

// Text element in block. Example types "plain_text", "mrkdwn".
type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Field element that can be a part of a block. Divides the block into arranged fields.
type BlockField struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Construct block of type "section" with one text element.
// Pass text type - "plain_text" or "markdown" and text.
func NewSectionTextBlock(textType string, text string) *Block {
	block := Block{}
	block.Type = "section"
	block.BText = &BlockText{Type: textType, Text: text}
	return &block
}

// Constructs block of type "section" with many field elements.
// Pass any number BlockField objects.
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

// Constructs block of type "header" with text.
// Pass the text of the header.
func NewHeaderBlock(text string) *Block {
	block := Block{}
	block.Type = "header"
	block.BText = &BlockText{Type: "plain_text", Text: text}
	return &block
}

// Constructs block of type "divider".
// It is a simple divider line.
func NewDividerBlock() *Block {
	block := Block{}
	block.Type = "divider"
	return &block
}

// Constructs a field to be put inside of block.
// Pass field type and text.
func NewField(fieldType string, text string) *BlockField {
	field := BlockField{}
	field.Type = fieldType
	field.Text = text
	return &field
}

// Constructs the final response to be returned to slack client.
// Pass any number of Block objects
func NewResponse(blocks ...*Block) *Response {
	resp := Response{}
	arr := make([]*Block, 0)
	for _, b := range blocks {
		arr = append(arr, b)
	}
	resp.Blocks = arr
	return &resp
}

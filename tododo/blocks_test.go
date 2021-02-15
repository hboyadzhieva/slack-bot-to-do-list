package tododo

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNewSectionTextBlock(t *testing.T) {
	expectedBlockText := &BlockText{
		Type: "mrkdwn",
		Text: "hello",
	}
	expected := &Block{
		Type:  "section",
		BText: expectedBlockText,
	}
	real := NewSectionTextBlock("mrkdwn", "hello")
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

func ExampleNewSectionTextBlock() {
	// Example
	block := NewSectionTextBlock("mrkdwn", "hello")
	bytes, _ := json.Marshal(block)
	fmt.Println(string(bytes))
	// Output: {"type":"section","text":{"type":"mrkdwn","text":"hello"}}
}

func TestNewSectionFieldsBlock(t *testing.T) {
	field1 := &BlockField{Type: "mrkdwn", Text: "field1"}
	field2 := &BlockField{Type: "mrkdwn", Text: "field2"}
	expected := &Block{
		Type:    "section",
		BFields: []*BlockField{field1, field2},
	}
	real := NewSectionFieldsBlock(field1, field2)
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

func TestNewHeaderBlock(t *testing.T) {
	expectedBlockText := &BlockText{
		Type: "plain_text",
		Text: "hello",
	}
	expected := &Block{
		Type:  "header",
		BText: expectedBlockText,
	}
	real := NewHeaderBlock("hello")
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

func TestNewDividerBlock(t *testing.T) {
	expected := &Block{
		Type: "divider",
	}
	real := NewDividerBlock()
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

func TestNewField(t *testing.T) {
	expected := &BlockField{
		Type: "mrkdwn",
		Text: "hello",
	}
	real := NewField("mrkdwn", "hello")
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

func TestNewResponse(t *testing.T) {
	block1 := &Block{
		Type: "divider",
	}
	block2 := &Block{
		Type:  "header",
		BText: &BlockText{Type: "plain_text", Text: "hello"},
	}
	expected := &Response{
		Blocks: []*Block{block1, block2},
	}
	real := NewResponse(block1, block2)
	if !reflect.DeepEqual(expected, real) {
		t.Errorf("Expected equal but not equal, expected: %v , real: %v", expected, real)
	}
}

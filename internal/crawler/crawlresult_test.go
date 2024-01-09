package crawler

import (
	"testing"
)

func TestMarshalBinary(t *testing.T) {
	p := PageData{
		URL:           "https://example.com",
		Links:         []string{"https://example.com/link1", "https://example.com/link2"},
		SearchTerms:   []string{"test", "search"},
		MatchingTerms: []string{"test", "match"},
		Error:         "",
	}

	_, err := p.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary() error = %v", err)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	p := &PageData{}
	data := []byte(`{
		"url": "https://example.com",
		"crawl_time": "2022-01-01T00:00:00Z",
		"status_code": 200,
		"metadata": {
			"description": "Test Description",
			"keywords": ["test", "keywords"]
		},
		"content": {
			"title": "Test Title",
			"body": "Test Body"
		},
		"links": ["https://example.com/link1", "https://example.com/link2"],
		"search_terms": ["test", "search"],
		"matching_terms": ["test", "match"],
		"error": ""
	}`)

	err := p.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("UnmarshalBinary() error = %v", err)
	}
}

func TestMarshalBinary_Error(t *testing.T) {
	p := PageData{
		URL: string([]byte{0x80, 0x81, 0x82}), // invalid string
	}

	_, err := p.MarshalBinary()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUnmarshalBinary_Error(t *testing.T) {
	p := &PageData{}
	data := []byte(`{
		"url": "https://example.com",
		"crawl_time": "invalid time", // invalid time format
		"status_code": 200
	}`)

	err := p.UnmarshalBinary(data)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

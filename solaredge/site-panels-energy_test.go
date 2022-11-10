package solaredge

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/alistairpialek/solaredge-panels-go/mocks"
)

func TestErrorStatusCode(t *testing.T) {
	c := NewClient(&mocks.MockClient{}, "abc", "123")

	// Build a response JSON.
	json := `{"error":"message"}`

	// Create a new reader with that JSON.
	reader := io.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 500,
			Body:       reader,
		}, nil
	}

	data, err := c.Site.PanelsEnergy("123")
	if err == nil {
		t.Errorf("expected an error, got nil")
	}

	if data != nil {
		t.Errorf("expected no data, got %+v", data)
	}
}

func TestUnexpectedField(t *testing.T) {
	c := NewClient(&mocks.MockClient{}, "abc", "123")

	// Build a response JSON.
	json := `{"error":"message"}`

	// Create a new reader with that JSON.
	reader := io.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       reader,
		}, nil
	}

	data, err := c.Site.PanelsEnergy("123")
	if err != nil {
		// Should not get an error if source and target json fields do not correspond.
		t.Errorf("got %s", err)
	}

	if len(data) != 0 {
		t.Errorf("expected zero length data, got %d", len(data))
	}
}

func TestSuccessfulUnmarshal(t *testing.T) {
	c := NewClient(&mocks.MockClient{}, "abc", "123")

	// Build a response JSON.
	jsonFile, err := os.Open("../mocks/expected.json")
	if err != nil {
		t.Errorf("got %s", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Create a new reader with that JSON.
	reader := io.NopCloser(bytes.NewReader(byteValue))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       reader,
		}, nil
	}

	data, err := c.Site.PanelsEnergy("123")
	if err != nil {
		t.Errorf("got %s", err)
	}

	if len(data) == 0 {
		t.Errorf("expected non-zero length data")
	}

	// E.g. 1.0.13
	want := `\d+.\d+.\d+`
	got := data[0].DisplayName
	match, _ := regexp.MatchString(want, got)
	if !match {
		t.Errorf("expected regex match for %s, got %s", want, got)
	}

	// E.g. 1562.500000
	got2 := data[0].Energy
	if got2 < 0 {
		t.Errorf("expected positive energy value, got %f", got2)
	}
}

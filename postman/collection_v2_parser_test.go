package postman

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestV2CanParse(t *testing.T) {
	parser := CollectionV2Parser{}

	file := "tests_data/collection-02.json"
	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed read data json %+v", err)
		return
	}

	collection, err := parser.Parse(b, BuilderOptions{})
	if err != nil {
		t.Errorf("Error when parse json")
		return
	}

	json, _ := json.Marshal(collection)
	t.Log(string(json))
}

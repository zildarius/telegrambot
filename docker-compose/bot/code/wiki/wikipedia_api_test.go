package wiki

import (
	"testing"
)

func TestWikipedia_api(t *testing.T) {

	if s := WikipediaAPI(`test`); s == nil {
		t.Fail()
	}

	q, _ := URLEncoded("325f547b68987098234c46:::::@#$%!@#@^%&&^*")

	request := "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + q + "&limit=3&origin=*&format=json"

	if s := WikipediaAPI(request); s[0] != "Something going wrong, try to change your question" {
		t.Fail()
	}

	q1, _ := URLEncoded("test")

	request1 := "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + q1 + "&limit=3&origin=*&format=json"

	if s := WikipediaAPI(request1); len(s) != 3 {

	}
}

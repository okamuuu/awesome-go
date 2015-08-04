package main

/*
http://qiita.com/umisama/items/3d44560e1c06fa531069
*/

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error("unexpected")
		return
	}
	if res.StatusCode != 200 {
		t.Error("Status code error")
		return
	}
}

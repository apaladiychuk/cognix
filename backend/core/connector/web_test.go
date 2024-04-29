package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"testing"
)

func TestWeb_Execute(t *testing.T) {
	web, err := NewWeb(
		&model.Connector{
			ID:        1,
			Name:      "web test",
			Source:    "web",
			InputType: "src",
			ConnectorSpecificConfig: model.JSONMap{
				"url":  "https://help.collaboard.app/",
				"url2": "https://apidog.com/blog/openapi-specification/",
				"url1": "https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML",
			},
			DocsMap: make(map[string]*model.Document),
		})
	if err != nil {
		t.Log(err.Error())
		t.Fatal(err)
	}
	conn, err := web.Execute(context.Background(), nil)
	if err != nil {
		t.Log(err.Error())
		t.Fatal(err)
	}
	for url, history := range (web.(*Web)).history {
		if len(history) > 30 {
			history = history[:30]
		}
		t.Logf("%s => %s ", url, history)
	}
	for _, doc := range conn.Docs {
		t.Log(doc.DocumentID)
	}
}

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
				"url": "https://en.wikipedia.org/",
			},
		})
	if err != nil {
		t.Log(err.Error())
		t.Fatal(err)
	}
	if err = web.Execute(context.Background(), nil); err != nil {
		t.Log(err.Error())
		t.Fatal(err)
	}
}

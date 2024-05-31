package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/gabriel-vasile/mimetype"
	_ "github.com/gabriel-vasile/mimetype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

var param = model.JSONMap{
	"token": model.JSONMap{
		"access_token":  "",
		"token_type":    "Bearer",
		"refresh_token": "",
		"expiry":        "2024-05-24T12:24:33.327841375+03:00",
	},
}

func TestOneDrive_Execute(t *testing.T) {
	od, err := NewOneDrive(&model.Connector{
		ID:                      decimal.Decimal{},
		CredentialID:            decimal.NullDecimal{},
		Name:                    "",
		Type:                    "",
		ConnectorSpecificConfig: param,
	})
	assert.NoError(t, err)
	chResult := od.Execute(context.TODO(), nil)
	for data := range chResult {
		mtype := mimetype.Detect(data.Content)
		t.Logf("%s -- %s ", string(data.MimeType), mtype)

	}
}

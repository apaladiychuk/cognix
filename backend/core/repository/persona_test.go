package repository

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestPersonaRepo(t *testing.T) {
	cfg := Config{
		URL:       "postgres://root:123@localhost:26257/defaultdb?sslmode=disable",
		DebugMode: "true",
	}
	gdb, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		t.Error(err)
		return
	}
	persona := model.Persona{
		Name:             "my persona",
		Llm_id:           2,
		DefaultPersona:   false,
		Description:      "dest",
		Tenant_id:        uuid.New(),
		Search_type:      "KEYWORD",
		Is_visible:       false,
		Display_priority: 0,
		Starter_messages: nil,
	}

	err = gdb.Table("personas").Create(&persona).Error
	if err != nil {
		t.Error(err)
	}
	t.Log(persona.Id)

	persona.Id = 0

	db, err := NewDatabase(&cfg)
	db = pg.Connect(&pg.Options{
		Network:  "",
		Addr:     "localhost:26257",
		User:     "root",
		Password: "123",
		Database: "defaultdb",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err = db.Ping(context.Background()); err != nil {
		t.Error(err)
		return
	}
	_, err = db.Model(&persona).Insert()
	if err != nil {
		t.Error(err)
	}
	t.Log(persona.Id)
}

package connector

import (
	_ "github.com/gabriel-vasile/mimetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestOneDrive_Folder(t *testing.T) {

	c := OneDrive{
		param: &OneDriveParameters{
			Folder:    "",
			Recursive: false,
			Token:     oauth2.Token{},
		},
	}
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), false)
	assert.Equal(t, c.isFilesAnalysing(""), true)
	assert.Equal(t, c.isFilesAnalysing("folder"), false)

	c.param.Recursive = true
	assert.Equal(t, c.isFolderAnalysing(""), true)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFilesAnalysing(""), true)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)

	c.param.Folder = "folder"

	assert.Equal(t, c.isFolderAnalysing(""), false)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)
	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)

	assert.Equal(t, c.isFilesAnalysing("folder/chapter1"), true)
	assert.Equal(t, c.isFilesAnalysing("folder2/chapter1"), false)

	c.param.Recursive = false
	assert.Equal(t, c.isFolderAnalysing(""), false)
	assert.Equal(t, c.isFolderAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), false)
	assert.Equal(t, c.isFilesAnalysing(""), false)
	assert.Equal(t, c.isFilesAnalysing("folder"), true)
	assert.Equal(t, c.isFolderAnalysing("folder/chapter1"), false)
	assert.Equal(t, c.isFolderAnalysing("folder2/chapter1"), false)

}

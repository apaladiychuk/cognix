package bll

import "cognix.ch/api/v2/core/repository"

type (
	DocumentBL interface {
	}
	documentBL struct {
		documentRepo repository.DocumentRepository
	}
)

func NewDocumentBL(documentRepo repository.DocumentRepository) DocumentBL {
	return documentBL{documentRepo: documentRepo}
}

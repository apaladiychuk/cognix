package parameters

type DocumentSetParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DocumentSetConnectorPairsParam struct {
	DocumentSetID int64   `json:"document_set_id"`
	ConnectorIDs  []int64 `json:"connector_ids"`
}

type DocumentUploadResponse struct {
	FileName string      `json:"file_name"`
	Error    string      `json:"error"`
	Document interface{} `json:"document"`
}

package parameters

type DocumentSetParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DocumentSetConnectorPairsParam struct {
	DocumentSetID int64   `json:"documentSetID"`
	ConnectorIDs  []int64 `json:"connectorIDs"`
}

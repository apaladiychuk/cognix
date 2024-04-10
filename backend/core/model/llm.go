package model

type LLM struct {
	tableName struct{} `pg:"llm"`
	ID        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	ModelID   string   `json:"model_id,omitempty"`
	Url       string   `json:"url,omitempty"`
	ApiKey    string   `json:"-"`
	Endpoint  string   `json:"endpoint,omitempty"`
}

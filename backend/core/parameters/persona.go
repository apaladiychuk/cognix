package parameters

import (
	"cognix.ch/api/v2/core/ai"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type PersonaParam struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ModelID         string            `json:"model_id"`
	URL             string            `json:"url"`
	APIKey          string            `json:"api_key"`
	Endpoint        string            `json:"endpoint"`
	SystemPrompt    string            `json:"system_prompt"`
	TaskPrompt      string            `json:"task_prompt"`
	StarterMessages []*StarterMessage `json:"starter_messages,omitempty"`
}

type StarterMessage struct {
	Name        string `json:"name"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (v PersonaParam) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name, validation.Required),
		validation.Field(&v.ModelID, validation.Required,
			validation.By(func(value interface{}) error {
				if _, ok := ai.SupportedModels[v.ModelID]; !ok {
					return fmt.Errorf("model %s not supported", v.ModelID)
				}
				return nil
			})),
		validation.Field(&v.APIKey, validation.Required))
}

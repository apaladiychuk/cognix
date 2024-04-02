package parameters

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateChatSession struct {
	Description string `json:"description"`
	PersonaID   int64  `json:"persona_id"`
	OneShot     bool   `json:"one_shot"`
}

func (v CreateChatSession) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.PersonaID, validation.Required),
	)
}

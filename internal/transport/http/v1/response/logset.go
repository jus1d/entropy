package response

import "github.com/google/uuid"

type CreateLogset struct {
	LogsetID uuid.UUID `json:"logset_id"`
	Count    int       `json:"count"`
}

type GetLogset struct {
	UUID  uuid.UUID        `json:"uuid"`
	Count int              `json:"count"`
	Logs  []map[string]any `json:"logs"`
}

package egdm

type Entity struct {
	ID         string                 `json:"id,omitempty"`
	InternalID uint64                 `json:"internalId,omitempty"`
	Recorded   uint64                 `json:"recorded,omitempty"`
	IsDeleted  bool                   `json:"deleted,omitempty"`
	References map[string]interface{} `json:"refs"`
	Properties map[string]interface{} `json:"props"`
}

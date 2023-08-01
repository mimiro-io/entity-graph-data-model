package egdm

type Continuation struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func NewContinuation() *Continuation {
	c := &Continuation{}
	c.ID = "@continuation"
	return c
}

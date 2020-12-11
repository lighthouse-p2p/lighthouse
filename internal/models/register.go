package models

// RegisterRequest is used to serialize the register request
type RegisterRequest struct {
	PubKey   string `json:"pub_key"`
	NickName string `json:"nickname"`
}

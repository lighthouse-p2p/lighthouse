package models

// Metadata is the model for metadata.json file
type Metadata struct {
	PubKey   string `json:"pub_key"`
	PrivKey  string `json:"priv_key"`
	NickName string `json:"nickname"`
	Host     string `json:"host"`
}

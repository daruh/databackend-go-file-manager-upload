package main

type SaveFile struct {
	FileName   string            `json:"fileName"`
	TenantId   string            `json:"tenantId"`
	Mime       string            `json:"mime"`
	Category   string            `json:"category"`
	FileType   string            `json:"fileType"`
	User       string            `json:"user"`
	DeleteMark bool              `bson:"deleteMark"`
	ExpiryTs   int64             `bson:"expiryTs"`
	App        string            `json:"app"`
	Size       int64             `json:"size"`
	Tags       map[string]string `json:"tags"`
}

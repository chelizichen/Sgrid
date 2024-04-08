package dto

type PageBasicReq struct {
	Size    int    `json:"size,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Keyword string `json:"keyword,omitempty"`
	Id      int    `json:"id,omitempty"`
}

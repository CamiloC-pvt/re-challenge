package order

import "time"

type Order struct {
	Created  time.Time   `json:"created"`
	Deleted  bool        `json:"deleted"`
	ID       int32       `json:"id"`
	Modified time.Time   `json:"modified"`
	Packs    []OrderPack `json:"packs"`
	Size     int32       `json:"size"`
}

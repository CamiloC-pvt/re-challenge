package pack

import "time"

type Pack struct {
	Created time.Time `json:"created"`
	ID      int32     `json:"id"`
	Size    int32     `json:"size"`
}

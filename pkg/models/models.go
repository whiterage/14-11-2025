package models

import "time"

type LinkRequest struct {
	Links []string `json:"links"`
}

type LinkStatus struct {
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CheckTime time.Time `json:"check_time,omitempty"`
}

const (
	StatusPending      = "pending"
	StatusProcessing   = "processing"
	StatusDone         = "done"
	StatusAvailable    = "available"
	StatusNotAvailable = "not_available"
)

type Task struct {
	ID        int          `json:"links_num"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"`
	Results   []LinkStatus `json:"results"`
}

type ReportRequest struct {
	LinksList []int `json:"links_list"`
}

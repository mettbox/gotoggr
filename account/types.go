package account

import "time"

type Account struct {
	APIToken        string      `json:"api_token"`
	Timezone        string      `json:"timezone"`
	ID              int         `json:"id"`
	Workspaces      []Workspace `json:"workspaces"`
	Clients         []Client    `json:"clients"`
	Projects        []Project   `json:"projects"`
	TimeEntries     []TimeEntry `json:"time_entries"`
	BeginningOfWeek int         `json:"beginning_of_week"`
}

type Client struct {
	Wid  int    `json:"wid"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	Wid             int        `json:"wid"`
	ID              int        `json:"id"`
	Cid             int        `json:"cid"`
	Name            string     `json:"name"`
	Active          bool       `json:"active"`
	ServerDeletedAt *time.Time `json:"server_deleted_at,omitempty"`
}

type TimeEntry struct {
	Wid         int        `json:"wid,omitempty"`
	ID          int        `json:"id,omitempty"`
	Pid         int        `json:"pid"`
	Tid         int        `json:"tid"`
	Description string     `json:"description,omitempty"`
	Stop        *time.Time `json:"stop,omitempty"`
	Start       *time.Time `json:"start,omitempty"`
	Tags        []string   `json:"tags"`
	Duration    int64      `json:"duration,omitempty"`
	DurOnly     bool       `json:"duronly"`
}

type DetailedEntry struct {
	Workspace   string
	Pid         int
	Project     string
	Cid         int
	Client      string
	Description string
	Duration    int64
}

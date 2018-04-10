package types

type City struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Borders []int64 `json:"borders,omitempty"`
}

type Cities []City

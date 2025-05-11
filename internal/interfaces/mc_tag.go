package interfaces

type McTag struct {
	Name   string   `json:"-"`
	Values []string `json:"values"`
}

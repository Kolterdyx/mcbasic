package interfaces

type DatapackHeader struct {
	Namespace   string `json:"namespace"`
	Definitions struct {
		Functions []struct {
			Name string `json:"name"`
			Args []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"args"`
			ReturnType string `json:"returnType"`
		} `json:"functions"`
	} `json:"definitions"`
}

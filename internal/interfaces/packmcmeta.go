package interfaces

import "github.com/Kolterdyx/mcbasic/internal/packformat"

type PackMcMetaPack struct {
	Description string                `json:"description"`
	PackFormat  packformat.PackFormat `json:"pack_format"`
}

type PackMcMetaMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackMcMeta struct {
	Pack PackMcMetaPack `json:"pack"`
	Meta PackMcMetaMeta `json:"meta"`
}

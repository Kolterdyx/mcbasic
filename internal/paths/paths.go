package paths

import "path"

const (
	Functions = "function"
	Internal  = "internal"
	Branches  = "branches"
	Data      = "data"
	Tags      = "tags"
)

var (
	FunctionBranches = path.Join(Internal, Branches)
	MinecraftData    = path.Join(Data, "minecraft")
	McbData          = path.Join(Data, "mcb")
	MathData         = path.Join(Data, "math")
	McbFunctions     = path.Join(McbData, Functions)
	MinecraftTags    = path.Join(MinecraftData, Tags)
	MathFunctions    = path.Join(MathData, Functions)
)

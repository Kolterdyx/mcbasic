package ops

import "fmt"

type TextFormat struct {
	Id     string
	Format string
}

// Colors
var (
	Black         = TextFormat{Id: "black", Format: "§0"}
	DarkBlue      = TextFormat{Id: "dark_blue", Format: "§1"}
	DarkGreen     = TextFormat{Id: "dark_green", Format: "§2"}
	DarkAqua      = TextFormat{Id: "dark_aqua", Format: "§3"}
	DarkRed       = TextFormat{Id: "dark_red", Format: "§4"}
	DarkPurple    = TextFormat{Id: "dark_purple", Format: "§5"}
	Gold          = TextFormat{Id: "gold", Format: "§6"}
	Gray          = TextFormat{Id: "gray", Format: "§7"}
	DarkGray      = TextFormat{Id: "dark_gray", Format: "§8"}
	Blue          = TextFormat{Id: "blue", Format: "§9"}
	Green         = TextFormat{Id: "green", Format: "§a"}
	Aqua          = TextFormat{Id: "aqua", Format: "§b"}
	Red           = TextFormat{Id: "red", Format: "§c"}
	LightPurple   = TextFormat{Id: "light_purple", Format: "§d"}
	Yellow        = TextFormat{Id: "yellow", Format: "§e"}
	White         = TextFormat{Id: "white", Format: "§f"}
	Reset         = TextFormat{Id: "reset", Format: "§r"}
	Bold          = TextFormat{Id: "bold", Format: "§l"}
	Italic        = TextFormat{Id: "italic", Format: "§o"}
	Underline     = TextFormat{Id: "underline", Format: "§n"}
	Strikethrough = TextFormat{Id: "strikethrough", Format: "§m"}
	Obfuscated    = TextFormat{Id: "obfuscated", Format: "§k"}
)

func (o *Op) Exception(message string) string {
	cmd := ""
	cmd += o.LoadArgConst("error", "text", fmt.Sprintf("Exception: %s", message), true)
	cmd += o.Call("mcb:error", "")
	// TODO: Implement schedule clear
	cmd += o.Return()
	return cmd
}

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

func (o *Op) Trace(desc string, content string) string {
	cmd := ""

	if o.EnableTraces {
		// Add color to the end of the content
		// All but the last character
		content = fmt.Sprintf("%s, \"color\": \"%s\"}", content[:len(content)-1], Green.Id)
		cmd += "#@ BEGIN TRACE\n"
		cmd += fmt.Sprintf("tellraw @a [{\"text\": \"%sTrace %s: \"},%s]\n", Red.Format, desc, content)
		cmd += "#@ END TRACE\n"
	}
	return cmd
}

func (o *Op) TraceStorage(storage string, path string) string {
	if path == "" {
		return o.Trace(fmt.Sprintf("(%s)", storage), fmt.Sprintf("{\"storage\":\"%s\"}", storage))
	} else {
		return o.Trace(fmt.Sprintf("(%s: %s)", storage, path), fmt.Sprintf("{\"storage\":\"%s\",\"nbt\":\"%s\"}", storage, path))
	}
}

func (o *Op) TraceScore(varName string, ns string) string {
	return o.Trace(fmt.Sprintf("(%s:%s)", ns, varName), fmt.Sprintf("{\"score\":{\"name\":\"%s\",\"objective\":\"%s\"}}]", varName, ns))
}

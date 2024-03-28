package ops

// Print prints a message to the console.
func (o *Op) Print(message string) string {
	cmd := ""
	cmd += o.Set("mcb/message", message)
	cmd += o.ArgLoad("builtin/print", "text", "mcb/message")
	cmd += o.Call("builtin/print")
	return cmd
}

package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (o *Op) PathGet(obj, path, to string) string {
	cmd := "### BEGIN List index operation  ###\n"
	cmd += o.LoadArgConst("internal/path/get", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArgConst("internal/path/get", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/path/get", "obj", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, obj)))
	cmd += o.LoadArg("internal/path/get", "path", path)
	cmd += o.Call("mcb:internal/path/get", to)
	cmd += "### END   List index operation ###\n"
	return cmd
}

func (o *Op) PathSet(obj, path, valuePath string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/path/set", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/path/set", "obj", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, obj)))
	cmd += o.LoadArgConst("internal/path/set", "value_path", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, valuePath)))
	cmd += o.LoadArg("internal/path/set", "path", path)
	cmd += o.Call("mcb:internal/path/set", "")
	return cmd
}

func (o *Op) PathDelete(obj, path string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/path/delete", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/path/delete", "obj", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, obj)))
	cmd += o.LoadArg("internal/path/delete", "path", path)
	cmd += o.Call("mcb:internal/path/delete", "")
	return cmd
}

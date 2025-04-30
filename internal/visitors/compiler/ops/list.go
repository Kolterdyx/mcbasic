package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (o *Op) MakeList(to string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s set value []\n", VarPath, to)
}

func (o *Op) AppendList(to, from string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s append from storage %s:data %s.%s\n", VarPath, to, o.Namespace, VarPath, from)
}

func (o *Op) GetListIndex(from, index, to string) string {
	cmd := "### BEGIN List index operation  ###\n"
	cmd += o.LoadArgConst("internal/list_index/get", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArgConst("internal/list_index/get", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/list_index/get", "from", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, from)))
	cmd += o.LoadArg("internal/list_index/get", "index", index)
	cmd += o.Call("mcb:internal/list_index/get", to)
	cmd += "### END   List index operation ###\n"
	return cmd
}

func (o *Op) SetListIndex(list, index, valuePath string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/set", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/list_index/set", "list", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, list)))
	cmd += o.LoadArg("internal/list_index/set", "index", index)
	cmd += o.LoadArgConst("internal/list_index/set", "value_path", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, valuePath)))
	cmd += o.Call("mcb:internal/list_index/set", "")
	return cmd
}

func (o *Op) DeleteListIndex(list, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/delete", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/list_index/delete", "list", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, list)))
	cmd += o.LoadArg("internal/list_index/delete", "index", index)
	cmd += o.Call("mcb:internal/list_index/delete", "")
	return cmd
}

func (o *Op) MakeIndex(res, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/make_index", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/list_index/make_index", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArg("internal/list_index/make_index", "index", index)
	cmd += o.Call("mcb:internal/list_index/make_index", res)
	return cmd
}

package ops

import (
	"fmt"
)

func (o *Op) MakeList(to string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s set value []\n", VarPath, to)
}

func (o *Op) AppendList(to, from string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s append from storage %s:data %s.%s\n", VarPath, to, o.Namespace, VarPath, from)
}

func (o *Op) GetListIndex(from, index, to string) string {
	cmd := "### BEGIN List index operation  ###\n"
	cmd += o.LoadArgConst("internal/list_index/get", "res", fmt.Sprintf("%s.%s", VarPath, RET))
	cmd += o.LoadArgConst("internal/list_index/get", "storage", fmt.Sprintf("%s:data", o.Namespace))
	cmd += o.LoadArgConst("internal/list_index/get", "from", fmt.Sprintf("%s.%s", VarPath, from))
	cmd += o.LoadArg("internal/list_index/get", "index", index)
	cmd += o.Call("mcb:internal/list_index/get", to)
	cmd += "### END   List index operation ###\n"
	return cmd
}

func (o *Op) SetListIndex(list, index, valuePath string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/set", "storage", fmt.Sprintf("%s:data", o.Namespace))
	cmd += o.LoadArgConst("internal/list_index/set", "list", fmt.Sprintf("%s.%s", VarPath, list))
	cmd += o.LoadArg("internal/list_index/set", "index", index)
	cmd += o.LoadArgConst("internal/list_index/set", "value_path", fmt.Sprintf("%s.%s", VarPath, valuePath))
	cmd += o.Call("mcb:internal/list_index/set", "")
	return cmd
}

func (o *Op) DeleteListIndex(list, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/delete", "storage", fmt.Sprintf("%s:data", o.Namespace))
	cmd += o.LoadArgConst("internal/list_index/delete", "list", fmt.Sprintf("%s.%s", VarPath, list))
	cmd += o.LoadArg("internal/list_index/delete", "index", index)
	cmd += o.Call("mcb:internal/list_index/delete", "")
	return cmd
}

func (o *Op) MakeIndex(res, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/make_index", "storage", fmt.Sprintf("%s:data", o.Namespace))
	cmd += o.LoadArgConst("internal/list_index/make_index", "res", fmt.Sprintf("%s.%s", VarPath, RET))
	cmd += o.LoadArg("internal/list_index/make_index", "index", index)
	cmd += o.Call("mcb:internal/list_index/make_index", res)
	return cmd
}

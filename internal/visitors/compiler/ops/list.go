package ops

import (
	"fmt"
)

func (o *Op) MakeList(to string) string {
	return fmt.Sprintf("data modify storage example:vars %s set value []\n", to)
}

func (o *Op) AppendList(to, from string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s append from storage %s:%s %s\n", VarPath, to, o.Namespace, VarPath, from)
}

func (o *Op) GetListIndex(from, index, to string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/get", "res", RET)
	cmd += o.LoadArgConst("internal/list_index/get", "storage", fmt.Sprintf("%s:data %s", o.Namespace, VarPath))
	cmd += o.LoadArgConst("internal/list_index/get", "from", from)
	cmd += o.LoadArg("internal/list_index/get", "index", index)
	cmd += o.Call("mcb:internal/list_index/get", to)
	return cmd
}

func (o *Op) SetListIndex(list, index, valuePath string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/set", "storage", fmt.Sprintf("%s:data %s", o.Namespace, VarPath))
	cmd += o.LoadArgConst("internal/list_index/set", "list", list)
	cmd += o.LoadArg("internal/list_index/set", "index", index)
	cmd += o.LoadArgConst("internal/list_index/set", "value_path", valuePath)
	cmd += o.Call("mcb:internal/list_index/set", "")
	return cmd
}

func (o *Op) DeleteListIndex(list, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index/delete", "storage", fmt.Sprintf("%s:data %s", o.Namespace, VarPath))
	cmd += o.LoadArgConst("internal/list_index/delete", "list", list)
	cmd += o.LoadArg("internal/list_index/delete", "index", index)
	cmd += o.Call("mcb:internal/list_index/delete", "")
	return cmd
}

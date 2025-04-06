package ops

import "fmt"

func (o *Op) MakeList(to string) string {
	return fmt.Sprintf("data modify storage example:vars %s set value []\n", to)
}

func (o *Op) AppendList(to, from string) string {
	return fmt.Sprintf("data modify storage example:vars %s append from storage %s:%s %s\n", to, o.Namespace, VarPath, from)
}

func (o *Op) GetListIndex(from, index, to string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/list_index", "res", RET)
	cmd += o.LoadArgConst("internal/list_index", "storage", fmt.Sprintf("%s:%s", o.Namespace, VarPath))
	cmd += o.LoadArgConst("internal/list_index", "from", from)
	cmd += o.LoadArg("internal/list_index", "index", index)
	cmd += o.Call("mcb:internal/list_index", to)
	return cmd
}

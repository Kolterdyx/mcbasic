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

func (o *Op) MakeIndex(res, index string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/path/make_index", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/path/make_index", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArg("internal/path/make_index", "index", index)
	cmd += o.Call("mcb:internal/path/make_index", res)
	return cmd
}

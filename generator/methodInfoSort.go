package generator

import "github.com/alvaroloes/ocgen/parser"

type MethodsInfoByPosStart []*parser.MethodInfo

func (mi MethodsInfoByPosStart) Len() int {
	return len(mi)
}

func (mi MethodsInfoByPosStart) Swap(i, j int) {
	mi[i], mi[j] = mi[j], mi[i]
}

func (mi MethodsInfoByPosStart) Less(i, j int) bool {
	return mi[i].PosStart < mi[j].PosStart
}

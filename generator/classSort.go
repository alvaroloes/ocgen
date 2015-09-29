package generator

import "github.com/alvaroloes/ocgen/parser"

type ClassesByAppearanceInMFile []parser.ObjCClass

func (oc ClassesByAppearanceInMFile) Len() int {
	return len(oc)
}

func (oc ClassesByAppearanceInMFile) Swap(i, j int) {
	oc[i], oc[j] = oc[j], oc[i]
}

func (oc ClassesByAppearanceInMFile) Less(i, j int) bool {
	return oc[i].NSCodingInfo.InitWithCoder.PosStart < oc[j].NSCodingInfo.InitWithCoder.PosStart
}


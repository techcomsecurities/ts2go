package gol

import (
	"fmt"

	"github.com/techcomsecurities/ts2go"
)

const StructDefFmt = "// %s generated code by ts2go\ntype %s struct {\n"
const StructPropDefFmt = "\t%s %s `json:\"%s\"`\n"
const StructPropGetterFmt = "func (s %s) %s() %s {\n\treturn s.%s\n}\n"

// Struct Golang struct
type Struct struct {
	Name     string
	IsPublic bool
	Props    []StructProp
}

// StructProp Golang struct property
type StructProp struct {
	Name       string
	Type       string // in case union type, eg: bid3Price: number | string
	IsReadOnly bool   // eg: readyonly bid3Price: number
}

// ======= StructProp ========
func (s StructProp) toPrintName() string {
	if s.IsReadOnly {
		return ts2go.LowerFirstChar(s.Name)
	}
	return ts2go.UpperFirstChar(s.Name)
}

func (s StructProp) toFmtDefPlaceData() []string {
	a := [...]string{s.toPrintName(), s.Type, s.Name}
	return a[:]
}

// readonly func
func (s StructProp) toFmtGetterPlaceData(structName string) []string {
	a := [...]string{structName, ts2go.UpperFirstChar(s.Name), s.Type, s.toPrintName()}
	return a[:]
}

// ======= Struct ========
func (s Struct) toPrintName() string {
	if !s.IsPublic {
		return ts2go.LowerFirstChar(s.Name)
	}
	return ts2go.UpperFirstChar(s.Name)
}

func (s Struct) toFmtDefPlaceData() []string {
	a := [...]string{s.toPrintName(), s.toPrintName()}
	return a[:]
}

func (s Struct) toFmtString() string {
	sFmt := StructDefFmt
	for range s.Props {
		sFmt = sFmt + StructPropDefFmt
	}
	sFmt = sFmt + CloseBlockFmt
	for _, p := range s.Props {
		if p.IsReadOnly {
			sFmt = sFmt + StructPropGetterFmt
		}
	}
	return sFmt
}

func (s Struct) toFmtPlaceData() []string {
	dt := make([]string, 0)
	dt = append(dt, s.toFmtDefPlaceData()...)
	for _, p := range s.Props {
		dt = append(dt, p.toFmtDefPlaceData()...)
	}
	for _, p := range s.Props {
		if p.IsReadOnly {
			dt = append(dt, p.toFmtGetterPlaceData(s.toPrintName())...)
		}
	}
	return dt
}

// ToPrintString generate printed string
func (s Struct) ToPrintString() string {
	strData := s.toFmtPlaceData()

	infData := make([]interface{}, 0, len(strData))
	for _, s := range strData {
		infData = append(infData, s)
	}

	return fmt.Sprintf(s.toFmtString(), infData...)
}

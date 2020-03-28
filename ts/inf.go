package ts

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/techcomsecurities/ts2go"
)

const InfDefRegExp = `(|export)( ?)interface (.*) {`
const InfPropDefRegExp = `(|readonly)( ?)([^?]*)(\??): (.*);`

// Inf ts interface
type Inf struct {
	Name     string
	IsExport bool
	Props    []InfProp
}

// InfProp ts interface property
type InfProp struct {
	Name       string
	Type       []string // in case union type, eg: bid3Price: number | string
	IsOpt      bool     // eg: bid3Price?: number
	IsReadOnly bool     // eg: readyonly bid3Price: number
}

// ExtractAllInf get all ts interface from a ts file
func ExtractAllInf(f ts2go.FFile) ([]Inf, error) {
	infDef := f.FindLinesWithPattern(InfDefRegExp, true)

	keys := make([]int, 0, len(infDef))
	for k := range infDef {
		keys = append(keys, k)
	}

	infs := make([]Inf, 0, len(infDef))
	for _, l := range keys {
		defWords := infDef[l]
		log.Infof("inf def %d %s", l, defWords[3])
		cls, _ := f.FindNearestLineWithPattern(l, CloseBlockRegExp, true)
		inf := toInf(defWords)

		propLines := f.GetLinesBetween(l+1, cls, true)
		log.Infof("no props %d", len(propLines))

		props := make([]InfProp, 0, len(propLines))
		for _, p := range propLines {
			pWords := ts2go.FindWordsWithPattern(p, InfPropDefRegExp)
			props = append(props, toInfProp(pWords))
		}

		inf.Props = props
		infs = append(infs, inf)
	}

	return infs, nil
}

func toInf(defWords []string) Inf {
	return Inf{
		Name:     strings.TrimSpace(defWords[3]),
		IsExport: strings.TrimSpace(defWords[1]) == "export",
	}
}

func toInfProp(pWords []string) InfProp {
	return InfProp{
		Name:       strings.TrimSpace(pWords[3]),
		Type:       ts2go.SplitAndTrim(pWords[5], "|", true),
		IsOpt:      strings.TrimSpace(pWords[4]) == "?",
		IsReadOnly: strings.TrimSpace(pWords[4]) == "readonly",
	}
}

package convert

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/techcomsecurities/ts2go"
	"github.com/techcomsecurities/ts2go/gol"
	"github.com/techcomsecurities/ts2go/ts"
)

// currently, just some basic types
var ts2GoTypeMap = map[string]string{
	"boolean": "bool",
	"number":  "float64",
	"string":  "string",
}

// I2SOpt config of i2s command
type I2SOpt struct {
	FilePath          string
	UnionTypePriority string
	OutputPath        string
}

// I2SConverter like its name
type I2SConverter struct {
	Opt I2SOpt
}

// NewI2SConverter create new I2SConverter
func NewI2SConverter(opt I2SOpt) I2SConverter {
	return I2SConverter{
		Opt: opt,
	}
}

// Run convert TS interface to Golang struct
func (s I2SConverter) Run() (string, error) {
	tsFile, err := ts2go.ToFFile(s.Opt.FilePath)
	if err != nil {
		log.WithError(err).Errorf("fail to create ts file %s", s.Opt.FilePath)
		return "", err
	}

	tsInfs, err := ts.ExtractAllInf(tsFile)
	if err != nil {
		log.WithError(err).Errorf("fail to extract interface from ts file %s", s.Opt.FilePath)
		return "", err
	}

	if s.Opt.OutputPath == "" {
		s.Opt.OutputPath = tsFile.Desc.NameWoExt + ".go"
	}

	log.Infof("To Golang struct")
	gos := make([]gol.Struct, 0, len(tsInfs))
	for _, tsInf := range tsInfs {
		goStruct := s.TSInf2GoStruct(tsInf)
		log.Infof(goStruct.Name)
		gos = append(gos, goStruct)
	}
	err = s.toOutputFile(gos)
	if err != nil {
		log.WithError(err).Errorf("fail to write Golang struct to file %s", s.Opt.OutputPath)
		return "", err
	}

	return s.Opt.OutputPath, nil
}

func (s I2SConverter) toOutputFile(gs []gol.Struct) error {
	f, err := os.Create(s.Opt.OutputPath)
	if err != nil {
		log.WithError(err).Errorf("failed opening file: %s", s.Opt.OutputPath)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.WriteString("package main\n\n"); err != nil {
		log.WithError(err).Error("fail to write file")
		return err
	}
	for _, s := range gs {
		if _, err := w.WriteString(s.ToPrintString() + "\n"); err != nil {
			log.WithError(err).Error("fail to write file")
			return err
		}
	}
	w.Flush()

	return nil
}

// TSInfProp2GoStructProp convert TS inf prop to Golang struct prop
func (s I2SConverter) TSInfProp2GoStructProp(tsP ts.InfProp) gol.StructProp {
	return gol.StructProp{
		Name:       tsP.Name,
		Type:       s.getGoType(tsP),
		IsReadOnly: tsP.IsReadOnly,
	}
}

// AllTSInfProp2GoStructProp convert arrays of TS inf prop
func (s I2SConverter) AllTSInfProp2GoStructProp(tsPs []ts.InfProp) []gol.StructProp {
	goSPs := make([]gol.StructProp, 0, len(tsPs))
	for _, p := range tsPs {
		goSPs = append(goSPs, s.TSInfProp2GoStructProp(p))
	}
	return goSPs
}

// TSInf2GoStruct convert TS inf to Golang struct
func (s I2SConverter) TSInf2GoStruct(tsI ts.Inf) gol.Struct {
	return gol.Struct{
		Name:     tsI.Name,
		IsPublic: tsI.IsExport,
		Props:    s.AllTSInfProp2GoStructProp(tsI.Props),
	}
}

func (s I2SConverter) getGoType(tsP ts.InfProp) string {
	tsType := tsP.Type[0]
	if s.Opt.UnionTypePriority == "l" {
		tsType = tsP.Type[len(tsP.Type)-1]
	}
	goType, ok := ts2GoTypeMap[tsType]
	if !ok {
		log.Fatalf("unsupported TS type %s", tsType)
	}
	return goType
}

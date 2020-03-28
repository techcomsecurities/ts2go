package ts2go

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FileDesc struct contain all descriptor of a .ts file
type FileDesc struct {
	FullPath  string
	NameWExt  string
	NameWoExt string
}

// FFile struct contain all information of a .ts file
type FFile struct {
	Desc  FileDesc
	Lines []string
	Loc   int // line of codes
}

// ToFFile create FFile from file path
func ToFFile(fp string) (FFile, error) {
	file, err := os.Open(fp)
	if err != nil {
		log.Errorf("failed opening file: %s", err)
		return FFile{}, err
	}
	defer file.Close()

	tsF := FFile{
		Desc: FileDesc{
			FullPath:  fp,
			NameWExt:  filepath.Base(fp),
			NameWoExt: strings.TrimSuffix(filepath.Base(fp), filepath.Ext(filepath.Base(fp))),
		},
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	tsF.Loc = len(txtlines)
	tsF.Lines = txtlines

	return tsF, nil
}

// FindLinesWithPattern return pairs <line number, line text> that match give pattern
func (f *FFile) FindLinesWithPattern(pattern string, isTrim bool) map[int]([]string) {
	m := make(map[int]([]string), 0)
	for no, line := range f.Lines {
		tmpLine := line
		if isTrim {
			tmpLine = strings.TrimSpace(tmpLine)
		}
		matches := FindWordsWithPattern(tmpLine, pattern)
		if len(matches) != 0 {
			m[no] = matches
		}
	}
	return m
}

// FindNearestLineWithPattern like its name
func (f *FFile) FindNearestLineWithPattern(from int, pattern string, isTrim bool) (int, string) {
	for i := from + 1; i < f.Loc; i++ {
		tmpLine := f.Lines[i]
		if isTrim {
			tmpLine = strings.TrimSpace(tmpLine)
		}
		if ExactlyMathPattern(tmpLine, pattern) {
			return i, tmpLine
		}
	}
	return -1, "" // if not found
}

// GetLinesBetween get all line between from (include) and to (exclude),
func (f *FFile) GetLinesBetween(from, to int, isTrim bool) []string {
	a := make([]string, 0)

	for i := from; i < to; i++ {
		tmpLine := f.Lines[i]
		if isTrim {
			tmpLine = strings.TrimSpace(tmpLine)
		}
		a = append(a, tmpLine)
	}

	return a
}

package i18ngo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	stateAwaitKey     = 1
	stateContentKey   = 2
	stateContentValue = 3
)

var (
	languageCode string //en
	objects      map[string]*PoFile
)

func setup() {
	if objects == nil {
		objects = make(map[string]*PoFile)
	}
}

func SetLanguageCode(lc string) {
	languageCode = lc
}

func GetLanguageCodes() []string {
	if objects == nil {
		return nil
	}
	output := make([]string, len(objects))
	i := 0
	for k, _ := range objects {
		output[i] = k
		i++
	}
	return output
}

func GetDefaultLanguageCode() string {
	if len(languageCode) < 1 {
		languageCode = "en"
	}
	return languageCode
}

func LoadPoAll(path string) error {
	var fil *PoFile
	var err error
	setup()
	vPath := func(path string, f os.FileInfo, err error) error {
		if strings.LastIndex(path, ".po") == len(path)-3 {
			fil, err = ParsePoFile(path)
			if fil != nil {
				if len(fil.LangCode) > 0 {
					objects[fil.LangCode] = fil
				}
			}
		}
		return nil
	}
	err = filepath.Walk(path, vPath)

	return err
}

func LoadPoFile(filename string) error {
	setup()
	pofile, err := ParsePoFile(filename)
	if err != nil {
		return err
	}
	objects[pofile.LangCode] = pofile
	return nil
}

func LoadPoStr(contents string) error {
	setup()
	pofile, err := ParsePoStr(contents)
	if err != nil {
		return err
	}
	objects[pofile.LangCode] = pofile
	return nil
}

func ParsePoFile(filename string) (*PoFile, error) {
	var bbytes []byte
	var err error
	bbytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parsePo(bbytes)
}

func ParsePoStr(contents string) (*PoFile, error) {
	return parsePo([]byte(contents))
}

func parsePo(bbytes []byte) (*PoFile, error) {
	buffer := bytes.NewBuffer(bbytes)
	var err error

	target := &PoFile{}
	target.Entries = make(map[string]string, 0)
	var currentKey, currentValue bytes.Buffer
	var line string
	state := stateAwaitKey
	for {
		line, err = readLine(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("readLine Error: " + err.Error())
			}
			if state == stateContentValue && currentValue.Len() > 0 {
				appendKeyVal(target, &currentKey, &currentValue)
			}
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			if state == stateContentValue && currentValue.Len() > 0 {
				appendKeyVal(target, &currentKey, &currentValue)
			}
			state = stateAwaitKey
			continue
		} else if line[0] == '#' { // comment
			continue
		}
		msgIdIndex := strings.Index(line, "msgid ")
		msgValIndex := strings.Index(line, "msgstr ")
		if msgIdIndex == 0 || msgValIndex == 0 {
			if state == stateContentValue && msgIdIndex == 0 {
				// append previous key/value
				if currentValue.Len() > 0 {
					appendKeyVal(target, &currentKey, &currentValue)
				}
			}
			if msgIdIndex == 0 {
				state = stateContentKey
			} else {
				state = stateContentValue
			}
			quote0 := strings.Index(line, "\"")
			if quote0 == -1 {
				continue
			}

			quote1 := strings.LastIndex(line, "\"")
			if quote0 == quote1 {
				continue
			}

			piece := line[quote0+1 : quote1]
			if msgIdIndex == 0 { // write new key
				currentKey.Reset()
				currentKey.WriteString(piece)
			} else { // write new value
				currentValue.Reset()
				currentValue.WriteString(piece)
			}
		} else if line[0] == '"' {
			lastIndex := strings.LastIndex(line, "\"")
			if lastIndex > 0 {
				line = line[1:lastIndex]
			} else {
				line = line[1:]
			}
			if state == stateContentKey {
				currentKey.WriteString(line)
			} else if state == stateContentValue {
				currentValue.WriteString(line)
			}
		}
	}
	return target, nil
}

func T(format string, a ...interface{}) string {
	setup()
	if len(languageCode) < 1 {
		return fmt.Sprintf(format, a...)
	}
	lang, ok := objects[languageCode]
	if !ok {
		return fmt.Sprintf(format, a...)
	}
	str, ok2 := lang.Entries[format]
	if !ok2 {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf(str, a...)
}

func TL(langcode string, format string, a ...interface{}) string {
	setup()
	lang, ok := objects[langcode]
	if !ok {
		return fmt.Sprintf(format, a...)
	}
	str, ok2 := lang.Entries[format]
	if !ok2 {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf(str, a...)
}

func readLine(reader *bytes.Buffer) (string, error) {
	return reader.ReadString('\n')
}

func appendKeyVal(file *PoFile, key *bytes.Buffer, val *bytes.Buffer) {
	if key.Len() == 0 {
		// HANDLE HEADER
		parseHeader(file, decodeKV(val.String()))
		key.Reset()
		val.Reset()
		return
	}
	file.Entries[decodeKV(key.String())] = decodeKV(val.String())
	key.Reset()
	val.Reset()
}

func decodeKV(input string) string {
	outp := strings.Replace(input, "\\n", "\n", -1)
	return strings.Replace(outp, "\\\"", "\"", -1)
}

func parseHeader(file *PoFile, header string) {
	lines := strings.Split(header, "\n")
	for _, v := range lines {
		if len(v) == 0 {
			continue
		}
		v = strings.TrimSpace(v)
		if strings.Index(v, "Language-Code:") == 0 {
			file.LangCode = strings.TrimSpace(v[14:])
		} else if strings.Index(v, "Language-Name:") == 0 {
			file.LangName = strings.TrimSpace(v[14:])
		}
		//TODO: add more header stuff
		//log.Println("PARSE HEADER:'" + v + "'")
	}
}

type PoFile struct {
	LangCode string
	LangName string
	Entries  map[string]string
}

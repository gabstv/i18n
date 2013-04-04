package i18ngo

import (
	"fmt"
	//"os"
	"testing"
)

func TestMultiple(t *testing.T) {
	user := "i18n"
	// English
	LoadPoStr("msgid \"\"\nmsgstr \"\"\n\"Content-Type: text/plain; charset=utf-8\\n\"\n\"Content-Transfer-Encoding: 8bit\\n\"\n\"Language-Code: en\\n\"\n\"Language-Name: English\\n\"\n\n# THIS COMMENT SHOULD BE IGNORED.\nmsgid \"Hello, %s!\"\nmsgstr \"Hello, %s!\"\n")
	// Portuguese
	LoadPoStr("msgid \"\"\nmsgstr \"\"\n\"Content-Type: text/plain; charset=utf-8\\n\"\n\"Content-Transfer-Encoding: 8bit\\n\"\n\"Language-Code: pt\\n\"\n\"Language-Name: Português\\n\"\n\n# THIS COMMENT SHOULD BE IGNORED.\nmsgid \"Hello, %s!\"\nmsgstr \"Olá, %s!\"\n")
	// German
	LoadPoStr("msgid \"\"\nmsgstr \"\"\n\"Content-Type: text/plain; charset=utf-8\\n\"\n\"Content-Transfer-Encoding: 8bit\\n\"\n\"Language-Code: de\\n\"\n\"Language-Name: Deutsch\\n\"\n\n# THIS COMMENT SHOULD BE IGNORED.\nmsgid \"Hello, %s!\"\nmsgstr \"Hallo, %s!\"\n")
	// set english as default language
	fmt.Println("S")
	fmt.Println(objects["pt"])
	fmt.Println("S")

	SetLanguageCode("en")
	t1 := T("Hello, %s!", user)
	fmt.Println(t1)
	if t1 != fmt.Sprintf("Hello, %s!", user) {
		t.Errorf("Failed English translation!")
	}
	// display portuguese translation
	t2 := TL("pt", "Hello, %s!", user)
	fmt.Println(t2)
	if t2 != fmt.Sprintf("Olá, %s!", user) {
		t.Errorf("Failed Portuguese translation!")
	}
	// display german translation
	t3 := TL("de", "Hello, %s!", user)
	fmt.Println(t3)
	if t3 != fmt.Sprintf("Hallo, %s!", user) {
		t.Errorf("Failed German translation!")
	}
}

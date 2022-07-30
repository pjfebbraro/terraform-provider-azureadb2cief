package util

import (
	"bytes"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"strings"
)

func XmlDiff(_, old, new string, _ *schema.ResourceData) bool {
	oldTokens, err := canonicalXML(old)

	if err != nil {
		return false
	}

	newTokens, err := canonicalXML(new)

	if err != nil {
		return false
	}

	return *oldTokens == *newTokens
}

func canonicalXML(s string) (*string, error) {
	reader := strings.NewReader(s)
	decoder := xml.NewDecoder(reader)
	var tokens []xml.Token
	var err error
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		if _, ok := token.(xml.ProcInst); ok {
			continue
		}
		if charData, ok := token.(xml.CharData); ok {
			if len(strings.TrimSpace(string(charData))) == 0 {
				continue
			}
		}
		if _, ok := token.(xml.Comment); ok {
			continue
		}
		tokens = append(tokens, xml.CopyToken(token))
	}

	var outBuffer bytes.Buffer
	encoder := xml.NewEncoder(&outBuffer)
	encoder.Indent("", " ")
	for _, v := range tokens {
		err = encoder.EncodeToken(v)
		if err != nil {
			return nil, err
		}
	}

	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	rawString := string(outBuffer.Bytes())
	return &rawString, nil
}

func NotCaseSensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

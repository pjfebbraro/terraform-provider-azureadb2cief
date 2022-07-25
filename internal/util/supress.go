package util

import (
	"bytes"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"regexp"
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

		tokens = append(tokens, xml.CopyToken(token))
	}

	var outBuffer bytes.Buffer
	encoder := xml.NewEncoder(&outBuffer)

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
	re := regexp.MustCompile(`\s`)
	results := re.ReplaceAllString(rawString, "")
	return &results, nil
}

func NotCaseSensitive(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

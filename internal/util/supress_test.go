package util_test

import (
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/util"
	"os"
	"testing"
)

func TestCanonicalXml(t *testing.T) {
	xmlbytes1, err := os.ReadFile("./testdata/TestXml.xml")
	if err != nil {
		t.Error(err)
	}
	xml := string(xmlbytes1)

	xmlBytes2, err := os.ReadFile("./testdata/TestXml2.xml")
	if err != nil {
		t.Error(err)
	}
	xml2 := string(xmlBytes2)

	isDifferent := util.XmlDiff("", xml, xml2, nil)
	if isDifferent == false {
		t.Fatalf("Diff returned false")
	}
}

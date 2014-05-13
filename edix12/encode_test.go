package edix12

import (
	"fmt"
	"testing"

	"github.com/kdar/health/edix12/node"
)

func TestAll(t *testing.T) {
	n := &node.Segment{
		Tag: "ISA",
		Data: []interface{}{
			"00", "No Auth   ", "00", "No Securit", "ZZ", "HIPAASpace     ",
			"ZZ", "Test Receiver  ", "140510", "1137", ">", "00501", "163636685",
			"0", "P", ":",
		},
	}

	data, err := Marshal(n)
	fmt.Println(data)
	fmt.Println(err)
}

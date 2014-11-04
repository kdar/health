package hl7v2_3

import "github.com/kdar/health/hl7"

var (
	_ = Obx{}
)

func stringIndex(data hl7.Data, n int) string {
	if v, ok := data.Index(n); ok {
		if v2, ok2 := v.(hl7.Field); ok2 {
			return string(v2)
		}
	}

	return ""
}

// func TestSomething(t *testing.T) {
// 	msh := "MSH|^~\\&|some company|some company|External|External|20140922091808||ORU^R01|20140922091808|P|2.3\r"
// 	obx := "OBX|2|CE|001719^HIV-1 ABS, SEMI-QN^L||HTN|||||N|F|19910123|| 19980729155700|BN"

// 	segments, err := hl7.Unmarshal([]byte(msh + obx))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, segment := range segments {
// 		if stringIndex(segment, 0) == "MSH" {
// 			var msh Msh
// 			err = hl7x.Unmarshal(segment, &msh)
// 			goon.Dump(msh)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 		}
// 		// if stringIndex(segment, 0) == "OBX" {
// 		// 	var obx Obx
// 		// 	err = unmarshal(&obx, segment)
// 		// 	goon.Dump(obx)
// 		// 	if err != nil {
// 		// 		t.Fatal(err)
// 		// 	}
// 		// }
// 	}
// }

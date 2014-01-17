package util

import (
	"fmt"
)

// IdentifyNDCFormat attempts to identify what format
// the given NDC is in.
func IdentifyNDCFormat(ndc string) string {
	var ret [3]int
	dashCount := 0

	for _, v := range ndc {
		if v == '-' {
			dashCount++
		} else {
			ret[dashCount]++
		}
	}

	if dashCount < 2 {
		return ""
	}

	return fmt.Sprintf("%d-%d-%d", ret[0], ret[1], ret[2])
}

// NormalizeNDC transforms a given NDC code to HIPAA 11-digit format
// This function is slightly optimized.
// www.nlm.nih.gov/research/umls/rxnorm/NDC_Normalization_Code.rtf‎
func NormalizeNDC(ndc string) (string, error) {
	var parts [3][]rune
	//buf := &bytes.Buffer{}
	//buf.Grow(11)
	dashCount := 0

	for _, v := range ndc {
		if v == '-' {
			dashCount++
			if dashCount > 2 {
				return "", fmt.Errorf("NDC is invalid.")
			}
		} else if !('0' <= v && v <= '9') && v != '*' {
			return "", fmt.Errorf("NDC is invalid.")
		} else if v == '*' {
			parts[dashCount] = append(parts[dashCount], '0')
		} else {
			parts[dashCount] = append(parts[dashCount], v)
		}
	}

	if len(parts[0]) < 5 {
		zeros := []rune{'0', '0', '0', '0', '0'}
		parts[0] = append(zeros[:5-len(parts[0])], parts[0]...)
	}
	if len(parts[1]) < 4 {
		zeros := []rune{'0', '0', '0', '0', '0'}
		parts[1] = append(zeros[:4-len(parts[1])], parts[1]...)
	}
	if len(parts[2]) < 2 {
		zeros := []rune{'0', '0', '0', '0', '0'}
		parts[2] = append(zeros[:2-len(parts[2])], parts[2]...)
	}

	if dashCount == 2 {
		ndc = string(parts[0][len(parts[0])-5:]) +
			string(parts[1][len(parts[1])-4:]) +
			string(parts[2][len(parts[2])-2:])
	} else if dashCount == 0 && len(ndc) == 12 && ndc[0] == '0' {
		ndc = ndc[1:]
	}

	return ndc, nil
}

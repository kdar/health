package util

import (
	"fmt"
	"strings"
	"unicode"
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
// This functions is unoptimized.
// www.nlm.nih.gov/research/umls/rxnorm/NDC_Normalization_Code.rtf‎
func NormalizeNDC(ndc string) (string, error) {
	var parts [3][]rune
	dashCount := 0

	for _, v := range ndc {
		if v == '-' {
			dashCount++
		} else if !unicode.IsDigit(v) && v != '*' {
			return "", fmt.Errorf("NDC is invalid.")
		} else {
			parts[dashCount] = append(parts[dashCount], v)
		}
	}

	for len(parts[0]) < 5 {
		parts[0] = append([]rune{'0'}, parts[0]...)
	}
	for len(parts[1]) < 4 {
		parts[1] = append([]rune{'0'}, parts[1]...)
	}
	for len(parts[2]) < 2 {
		parts[2] = append([]rune{'0'}, parts[2]...)
	}

	if dashCount == 2 {
		ndc = fmt.Sprintf("%s%s%s",
			string(parts[0][len(parts[0])-5:]),
			string(parts[1][len(parts[1])-4:]),
			string(parts[2][len(parts[2])-2:]))
	} else if dashCount == 0 && len(ndc) == 12 && ndc[0] == '0' {
		ndc = ndc[1:]
	}

	ndc = strings.Replace(ndc, "*", "0", -1)

	return ndc, nil
}

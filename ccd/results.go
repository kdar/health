package ccd

import (
	"errors"
	"github.com/jteeuwen/go-pkg-xmlx"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ResultsTid = []string{"2.16.840.1.113883.10.20.1.14", "2.16.840.1.113883.10.20.22.2.3.1"}

	ResultsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ResultsTid,
		Priority: 0,
		Func:     parseResults,
	}

	RANGE_SPLIT_RE = regexp.MustCompile(`\s*(;|,|\|)\s*`)
	RANGE_RE       = regexp.MustCompile(`(?P<text>[a-zA-Z\s]*?)\s*\(?(?P<low>[\d.]+)\s*[-–]\s*(?P<high>[\d.]+).*?\)?`)
	RANGE_MATH_RE  = regexp.MustCompile(`(?P<text>[a-zA-Z\s]*?)\s*\(?(?P<symbol>[<>=]+)\s*(?P<value>[\d.]+).*?\)?`)
)

type ResultValue struct {
	Type  string
	Value string
	Unit  string
}

type ResultRange struct {
	Gender  *string // M or F
	AgeLow  *float64
	AgeHigh *float64
	Low     *float64
	High    *float64
	Text    *string
}

func (r ResultRange) IsZero() bool {
	return r.Gender == nil && r.AgeLow == nil &&
		r.AgeHigh == nil && r.Low == nil &&
		r.High == nil && r.Text == nil
}

type ResultRanges []ResultRange

func (r *ResultRanges) Parse(s string) {
	for _, part := range RANGE_SPLIT_RE.Split(s, -1) {
		if part == "NA" || part == "No data" {
			continue
		}

		rr := ResultRange{}

		if strings.HasPrefix(part, "M ") {
			gender := "M"
			rr.Gender = &gender
			part = part[2:]
		} else if strings.HasPrefix(part, "F ") {
			gender := "F"
			rr.Gender = &gender
			part = part[2:]
		}

		text := ""
		colonsplit := strings.Split(part, ":")
		if len(colonsplit) == 2 {
			text = colonsplit[0]
			part = colonsplit[1]
		}

		part = strings.Replace(part, "less than", "<", -1)
		part = strings.Replace(part, "below", "<", -1)
		part = strings.Replace(part, "greater than", ">", -1)
		part = strings.Replace(part, "above", ">", -1)
		part = strings.Replace(part, "equal to", "=", -1)

		err := parseRange(part, &text, &rr.Low, &rr.High)
		if err != nil {
			parseRangeMath(part, &text, &rr.Low, &rr.High)
		}

		// Handle when years are specified
		if strings.Contains(text, "years") {
			var empty string
			err := parseRange(text, &empty, &rr.AgeLow, &rr.AgeHigh)
			if err != nil {
				parseRangeMath(text, &empty, &rr.AgeLow, &rr.AgeHigh)
			}

			text = ""
		}

		if rr.IsZero() {
			part = strings.Trim(part, "()")
			if len(part) == 0 || strings.Contains(part, "/") {
				continue
			}
			text = part
		}

		if len(text) > 0 {
			rr.Text = &text
		}

		*r = append(*r, rr)
	}
}

// Find [numbers] - [numbers]
func parseRange(s string, text *string, low **float64, high **float64) error {
	data := RANGE_RE.FindStringSubmatch(s)
	if len(data) == 4 {
		if *text == "" {
			*text = data[1]
		}

		lowf, err := strconv.ParseFloat(data[2], 64)
		if err == nil {
			*low = &lowf
		}

		highf, err := strconv.ParseFloat(data[3], 64)
		if err == nil {
			*high = &highf
		}

		return nil
	}

	return errors.New("Not a range")
}

// Find <>=[numbers]
func parseRangeMath(s string, text *string, low **float64, high **float64) error {
	data := RANGE_MATH_RE.FindStringSubmatch(s)
	if len(data) == 4 {
		if *text == "" {
			*text = data[1]
		}

		value, err := strconv.ParseFloat(data[3], 64)
		if err == nil {
			switch data[2] {
			case "<", "<=":
				*high = &value
			case ">", ">=":
				*low = &value
			case "=", "==":
				*low = &value
				*high = &value
			}
		}
	}

	return errors.New("Not a math range")
}

type ResultObservation struct {
	Date                time.Time
	DisplayName         string
	Loinc               string
	Value               ResultValue
	InterpretationCodes []string
	Ranges              []ResultRange
}

type Result struct {
	Date         time.Time
	Observations []ResultObservation
}

func parseResults(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		organizerNode := Nget(entryNode, "organizer")
		if organizerNode == nil {
			return nil
		}

		componentNodes := organizerNode.SelectNodes("*", "component")
		if componentNodes == nil {
			return nil
		}

		result := Result{}

		effectiveTimeNode := Nget(organizerNode, "effectiveTime")
		t := ParseTimeNode(effectiveTimeNode)
		result.Date = t.Value

		for _, componentNode := range componentNodes {
			obNode := Nget(componentNode, "observation")
			if obNode == nil {
				continue
			}

			observation := ResultObservation{}

			effectiveTimeNode := Nget(obNode, "effectiveTime")
			t = ParseTimeNode(effectiveTimeNode)
			observation.Date = t.Value

			codeNode := Nget(obNode, "code")
			if codeNode != nil {
				observation.DisplayName = codeNode.As("*", "displayName")
				observation.Loinc = codeNode.As("*", "code")
			}

			valueNode := Nget(obNode, "value")
			if valueNode != nil {
				observation.Value.Type = valueNode.As("*", "type")
				observation.Value.Value = valueNode.As("*", "value")
				observation.Value.Unit = valueNode.As("*", "unit")
			}

			icodeNodes := obNode.SelectNodes("*", "interpretationCode")
			if icodeNodes != nil {
				for _, icodeNode := range icodeNodes {
					observation.InterpretationCodes = append(observation.InterpretationCodes, icodeNode.As("*", "code"))
				}
			}

			// <referenceRange>
			//      <observationRange>
			//        <text>
			//          <reference value="#TESTRANGE_3"/>
			//        </text>
			//        <value xsi:type="IVL_PQ">
			//          <low value="150" unit="10+3/ul"/>
			//          <high value="350" unit="10+3/ul"/>
			//        </value>
			//      </observationRange>
			//    </referenceRange>

			//    <referenceRange>
			//     <observationRange>
			//       <text>M 13-18 g/dl; F 12-16 g/dl<reference value="#TESTRANGE_1"/>
			//       </text>
			//     </observationRange>
			//   </referenceRange>

			obvRangeNode := Nget(obNode, "referenceRange", "observationRange")
			if obvRangeNode != nil {
				var resultRanges ResultRanges

				valueNode := Nget(obvRangeNode, "value")
				if valueNode != nil {
					lowNode := Nget(valueNode, "low")
					highNode := Nget(valueNode, "high")

					lowf, _ := strconv.ParseFloat(lowNode.As("*", "value"), 64)
					highf, _ := strconv.ParseFloat(highNode.As("*", "value"), 64)

					resultRanges = append(resultRanges, ResultRange{
						Low:  &lowf,
						High: &highf,
					})
				} else {
					rangeText := obvRangeNode.S("*", "text")
					if len(rangeText) > 0 {
						resultRanges.Parse(rangeText)
					}
				}

				observation.Ranges = resultRanges
			}

			result.Observations = append(result.Observations, observation)
		}

		ccd.Results = append(ccd.Results, result)
	}

	return nil
}

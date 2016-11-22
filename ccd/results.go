package ccd

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jteeuwen/go-pkg-xmlx"
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

//decodeResultValue decodes the <value> in an Observation Result
func decodeResultValue(n *xmlx.Node) ResultValue {
	//reference: http://www.cdapro.com/know/25047
	var rv ResultValue
	if n == nil {
		return rv
	}

	rv.Type = n.As("*", "type")
	switch rv.Type {

	case "PQ": //PhysicalQuantity(PQ) types contain a "unit" and a "value" attribute
		rv.Value = n.As("*", "value")
		rv.Unit = n.As("*", "unit")
	case "ST": //Character String (ST) types contain their value as the content of their node
		rv.Value = n.GetValue()

		//some EMRs use ST for data that should be a PQ, so we check for that.
		if strings.Count(rv.Value, " ") == 1 {
			parts := strings.Split(rv.Value, " ")
			value, unit := parts[0], parts[1]

			//values must be real numbers
			if _, err := strconv.ParseFloat(value, 64); err != nil {
				//probably a normal ST segment, don't do anything special.
				break
			}
			rv.Unit = unit
			rv.Value = value
			//TODO:  since we're essentially turning this ST into a PQ, should we change type to PQ?
			//as it is, if something is checking our result's Type and handling it specially, it wouldn't expect a ST result to have a .Unit
		}

		/*NYI: "ED"(Encapsulated Data) which is kind of a 'catch all' for any arbitrary data
		it can refer to elements anywhere in the document, and it could have any format, so we cant really support it
		*/
	case "CV", "CD": //Coded Value(CV) and ConceptDiscriptor(CD) are both similar.
		//this is a rather simplistic decoding, ignoring the codesystem and translations, but it works with our samples.
		//two of our sample ccdas use CV, but they both are nullflavor so we can't test them.
		//A lot more use "CD", when they aren't null our samples all use them for "positive" "negative" "normal" etc.
		rv.Value = n.As("*", "displayName")
	}

	return rv
}

type ResultRange struct {
	Gender       *string // M or F
	AgeLow       *float64
	AgeHigh      *float64
	Low          *float64
	High         *float64
	Text         *string
	OriginalText string
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

		rr.OriginalText = s

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

// Find <>=[numbers]. e.g. <5, >=6.5
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
	Code                Code
	Value               ResultValue
	InterpretationCodes []string
	Ranges              []ResultRange
}

type Result struct {
	Date         time.Time
	Code         Code
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

		codeNode := Nget(organizerNode, "code")
		if codeNode != nil {
			result.Code.decode(codeNode)
		}

		if len(result.Code.DisplayName) == 0 {
			continue
		}

		effectiveTimeNode := Nget(organizerNode, "effectiveTime")
		t := decodeTime(effectiveTimeNode)
		result.Date = t.Value

		for _, componentNode := range componentNodes {
			obNode := Nget(componentNode, "observation")
			if obNode == nil {
				continue
			}

			observation := ResultObservation{}

			effectiveTimeNode := Nget(obNode, "effectiveTime")
			t = decodeTime(effectiveTimeNode)
			observation.Date = t.Value

			codeNode := Nget(obNode, "code")
			if codeNode != nil {
				observation.Code.decode(codeNode)
			}

			if len(observation.Code.DisplayName) == 0 {
				continue
			}

			observation.Value = decodeResultValue(Nget(obNode, "value"))

			icodeNodes := obNode.SelectNodes("*", "interpretationCode")
			if icodeNodes != nil {
				for _, icodeNode := range icodeNodes {
					observation.InterpretationCodes = append(observation.InterpretationCodes, icodeNode.As("*", "code"))
				}
			}

			obvRangeNode := Nget(obNode, "referenceRange", "observationRange")
			if obvRangeNode != nil {
				var resultRanges ResultRanges

				valueNode := Nget(obvRangeNode, "value")
				if valueNode != nil {
					lowNode := Nget(valueNode, "low")
					highNode := Nget(valueNode, "high")
					if lowNode == nil || highNode == nil {
						continue
					}
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

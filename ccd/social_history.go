package ccd

import (
	"github.com/jteeuwen/go-pkg-xmlx"
	"time"
)

var (
	SocialHistoryTid = []string{"2.16.840.1.113883.10.20.1.15", "2.16.840.1.113883.10.20.22.2.17"}

	SocialHistoryParser = Parser{
		Type:     PARSE_SECTION,
		Values:   SocialHistoryTid,
		Priority: 0,
		Func:     parseSocialHistory,
	}
)

type SocialHistory struct {
	Code      Code
	Value     Code
	StartDate time.Time
	StopDate  time.Time
}

func parseSocialHistory(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		obvNode := Nget(entryNode, "observation")
		codeNode := Nget(obvNode, "code")
		valueNode := Nget(obvNode, "value")
		if codeNode == nil || valueNode == nil {
			continue
		}

		socialHistory := SocialHistory{}

		socialHistory.Code.decode(codeNode)
		if socialHistory.Code.DisplayName == "" {
			orgText := Nget(codeNode, "originalText")
			if orgText != nil {
				socialHistory.Code.DisplayName = orgText.S("*", "originalText")
			}
		}

		socialHistory.Value.decode(valueNode)

		effectiveTimeNode := Nget(obvNode, "effectiveTime")
		t := decodeTime(effectiveTimeNode)
		socialHistory.StartDate = t.Low
		socialHistory.StopDate = t.High

		ccd.SocialHistory = append(ccd.SocialHistory, socialHistory)
	}

	return nil
}

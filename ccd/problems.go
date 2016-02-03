package ccd

import "github.com/jteeuwen/go-pkg-xmlx"

var (
	ProblemsTid = []string{"2.16.840.1.113883.10.20.1.11", "2.16.840.1.113883.10.20.22.2.5.1"}

	ProblemsParser = Parser{
		Type:     PARSE_SECTION,
		Values:   ProblemsTid,
		Priority: 0,
		Func:     parseProblems,
	}
)

//Problem represents an Observation Problem  (templateId: 2.16.840.1.113883.10.20.22.4.4)
type Problem struct {
	Name string
	Time Time
	// Duration time.Duration
	Status      string
	ProblemType string
	Code        Code
}

func parseProblems(node *xmlx.Node, ccd *CCD) []error {
	entryNodes := node.SelectNodes("*", "entry")
	for _, entryNode := range entryNodes {
		observationNode := Nget(entryNode, "act", "entryRelationship", "observation")
		problem := decodeProblem(observationNode)

		ccd.Problems = append(ccd.Problems, problem)
	}

	return nil
}

func decodeProblem(node *xmlx.Node) Problem {
	problem := Problem{}

	valueNode := Nget(node, "value")
	if valueNode == nil {
		//the spec says there must be a value, but better to be safe than to panic.
		return Problem{}
	}
	problem.Name = valueNode.As("*", "displayName")

	//The Value node is a ConceptDescriptor, so we decode it as a Code.
	problem.Code.decode(valueNode)

	//get the problem type from the highest level code node
	if topCode := Nget(node, "code"); topCode != nil {
		name := topCode.As("*", "displayName")
		if name == "" {
			name = topCode.As("*", "code")
		}
		problem.ProblemType = name
	}

	if effectiveTimeNode := Nget(node, "effectiveTime"); effectiveTimeNode != nil {
		problem.Time = decodeTime(effectiveTimeNode)
	}

	// observationNode2 := Nget(observationNode, "entryRelationship", "observation")
	// if observationNode2 != nil {
	//   problem.Status = Nget(observationNode2, "value").As("*", "displayName")
	// }

	erNodes := node.SelectNodes("*", "entryRelationship")
	for _, erNode := range erNodes {
		oNode := Nget(erNode, "observation")
		codeNode := Nget(oNode, "code")
		valueNode := Nget(oNode, "value")

		if codeNode == nil {
			continue
		}

		if codeNode.As("*", "code") == "33999-4" { // Status
			problem.Status = valueNode.As("*", "displayName")
		}
	}
	return problem
}

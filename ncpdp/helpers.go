package ncpdp

import (
	"github.com/kdar/health/edifact"
)

func getString(values edifact.Values, index int) string {
	if len(values) > index {
		if str, ok := values[index].(string); ok {
			return str
		}
	}

	return ""
}

func getValues(values edifact.Values, index int) edifact.Values {
	if len(values) > index {
		if vals, ok := values[index].(edifact.Values); ok {
			return vals
		}
	}

	return edifact.Values{}
}

func getName(values edifact.Values, index int) *Name {
	subValues := getValues(values, index)
	n := &Name{}
	n.Last = getString(subValues, 0)
	n.First = getString(subValues, 1)
	n.Middle = getString(subValues, 2)
	n.Suffix = getString(subValues, 3)
	n.Prefix = getString(subValues, 4)

	if len(n.Last) > 0 || len(n.First) > 0 || len(n.Middle) > 0 || len(n.Suffix) > 0 || len(n.Prefix) > 0 {
		return n
	}

	return nil
}

func getAddress(values edifact.Values, index int) *Address {
	subValues := getValues(values, index)
	a := &Address{}
	a.Line1 = getString(subValues, 0)
	a.City = getString(subValues, 1)
	a.State = getString(subValues, 2)
	a.Postal = getString(subValues, 3)
	a.LocationQualifier = getString(subValues, 4)
	a.Location = getString(subValues, 5)

	return a
}

func getPhones(values edifact.Values, index int) []*Phone {
	var phones []*Phone
	subValues := getValues(values, index)
	for n := range subValues {
		subVals := getValues(subValues, n)
		phones = append(phones, &Phone{getString(subVals, 0), getString(subVals, 1)})
	}

	return phones
}

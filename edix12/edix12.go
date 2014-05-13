package edix12

// type Interchange struct {
// 	FunctionalGroups []FunctionalGroup
// }

// func (i *Interchange) AddFunctionalGroup() FunctionalGroup {
// 	return FunctionalGroup{}
// }

// type FunctionalGroup struct {
// 	Transactions []Transaction
// }

// func (f *FunctionalGroup) AddTransaction() Transaction {
// 	return Transaction{}
// }

// type Transaction struct {
// }

// func (t *Transaction) AddSegment() {}
// func (t *Transaction) AddHLoop()   {}

// type HierarchicalLoop struct {

// }

type Interchange struct {
	ElementSeparator    byte
	ComponentSeparator  byte
	SegmentTerminator   byte
	RepetitionSeparator byte
}

func NewInterchange() Interchange {
	return Interchange{
		ElementSeparator:    '*', // 4rd character
		ComponentSeparator:  ':', // 16th data element, first char
		SegmentTerminator:   '~', // 16th data element, second char
		RepetitionSeparator: '^', // 11th data element
	}
}

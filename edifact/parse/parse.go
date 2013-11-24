// Package parse builds parse trees of nodes to be
// executed elsewhere.
// Uses shift-reduce based parsing:
//   http://en.wikipedia.org/wiki/Shift-reduce_parser

package parse

import (
	"errors"
	"fmt"
	"github.com/kdar/health/edifact/token"
	"regexp"
)

// A function used to take the latest node off the stack,
// the lookahead token, and return a node that replaces
// the latest node of the stack.
type reduceFunc func(node Node, tok token.Token) (Node, error)

// Tree is the representation of a single parsed template.
type Tree struct {
	Root *ListNode // top-level root of the tree.

	// Parsing only; cleared after parse.
	stack     NodeStack
	lex       *lexer
	token     [1]token.Token // one-token lookahead for parser.
	peekCount int
}

var (
	// our reduce table
	REDUCETABLE = map[NodeType]map[token.TokenType]reduceFunc{
		NodeText: map[token.TokenType]reduceFunc{
			token.COMPONENT_DELIMITER:  reduceTextComponent,
			token.REPETITION_DELIMITER: reduceTextRepetition,
			token.DATA_DELIMITER:       reduceTextData,
		},

		NodeComponent: map[token.TokenType]reduceFunc{
			token.TEXT:                 reduceComponentText,
			token.DATA_DELIMITER:       reduceComponentData,
			token.REPETITION_DELIMITER: reduceComponentRepetition,
		},

		NodeRepetition: map[token.TokenType]reduceFunc{
			token.TEXT:                 reduceRepetitionText,
			token.DATA_DELIMITER:       reduceRepetitionData,
			token.COMPONENT_DELIMITER:  reduceRepetitionComponent,
			token.REPETITION_DELIMITER: reduceRepetitionRepetition,
		},
	}
)

// Parses the text and returns a parse tree.
func Parse(text string) (listnode *ListNode, err error) {
	tree := &Tree{
		Root: newList(),
		lex:  lex("", text),
	}

	// these are two regexps to help us in removing the
	// release indicator from text and replacing it if
	// necessary and appropriate
	// --
	// this first regex's job is to take any release indicator
	// that is not paired with a delimiter, and replace it with
	// a space. now this is not in the spec at all but I have seen
	// this in the wild (relayhealth) where they will use the release
	// indicator as a space. Such as: "CVS?PHARMACY"
	var releaseRegex1 *regexp.Regexp
	// this regex will simply just remove the release indicator
	// wherever it is paired with a delimiter. Such as:
	// "??" -> "?" and "?^_?^" -> "^_^"
	var releaseRegex2 *regexp.Regexp

LOOP:
	for {
		tok := tree.next()
		switch tok.Typ {
		case token.EOF:
			break LOOP
		case token.ERROR:
			return nil, errors.New(tok.Val)
		case token.SEGMENT_TERMINATOR:
			// If we get a segment terminator, then append it
			// to our root and clear the stack.
			seg := newSegment()
			seg.List.Nodes = append(seg.List.Nodes, tree.stack...)
			tree.Root.append(seg)
			tree.stack.clear()
		case token.UNA_SEGMENT:
			tree.stack.push(newText(tok.Val))
		case token.UNA_TEXT:
			hdr := newHeader()
			hdr.SegmentName = tree.stack.last()
			hdr.Text = newText(tok.Val)
			tree.Root.append(hdr)
			tree.stack.clear()

			// at this point our lex parsed all the delimiters.
			// so we can create our release regexps.

			// the %%s will get replaced later for the regex's
			// specific purpose.
			// i use QuoteMeta here just in case our delimiters
			// conflict with the regexp.
			baseRegStr := fmt.Sprintf(`%s([%%s%s])`,
				regexp.QuoteMeta(string(tree.lex.releaseIndicator)),
				regexp.QuoteMeta(fmt.Sprintf("%c%c%c%c%c", tree.lex.componentDelimiter, tree.lex.dataDelimiter,
					tree.lex.releaseIndicator, tree.lex.repetitionDelimiter,
					tree.lex.segmentTerminator)))

			releaseRegex1, err = regexp.Compile(fmt.Sprintf(baseRegStr, "^"))
			if err != nil {
				return nil, err
			}
			releaseRegex2, err = regexp.Compile(fmt.Sprintf(baseRegStr, ""))
			if err != nil {
				return nil, err
			}
		case token.TEXT:
			// explanation of this is commented by the regexp declarations
			tok.Val = releaseRegex1.ReplaceAllString(tok.Val, " $1")
			tok.Val = releaseRegex2.ReplaceAllString(tok.Val, "$1")
			fallthrough
		default:
			// if addToStack is true, then we push the text onto
			// the stack.
			addToStack := true
			if tree.stack.len() > 0 {
				lastnode := tree.stack.last()

				// Try to find a reduce function in our table for
				// the given last node on the stack and the lookahead
				// token.
				if tokMap, ok := REDUCETABLE[lastnode.Type()]; ok {
					if reducefn, ok := tokMap[tok.Typ]; ok {
						// we don't add this to stack since we found a
						// reduce function to handle it
						addToStack = false
						reducedNode, err := reducefn(lastnode, tok)
						if err != nil {
							return nil, err
						}

						// replace the last node of the stack with our
						// reduced node.
						tree.stack.setLast(reducedNode)
					}
				}
			}

			// add the text to the stack if we didn't find a reduce
			// function and this is a segment or text.
			if addToStack && (tok.Typ == token.SEGMENT || tok.Typ == token.TEXT) {
				tree.stack.push(newText(tok.Val))
			}
		}
	}

	return tree.Root, nil
}

// reduce a text-data with a new data as text as its child
func reduceTextData(node Node, tok token.Token) (Node, error) {
	return newData(node), nil
}

// reduce a text-component to a new component with text appended
func reduceTextComponent(node Node, tok token.Token) (Node, error) {
	c := newComponent()
	c.List.append(node)
	return c, nil
}

// reduce a text-repetition to a new repetition with text appended
func reduceTextRepetition(node Node, tok token.Token) (Node, error) {
	r := newRepetition()
	r.List.append(node)
	return r, nil
}

// reduce a component-text to the component with text appended
func reduceComponentText(node Node, tok token.Token) (Node, error) {
	c, _ := node.(*ComponentNode)
	c.List.append(newText(tok.Val))
	return c, nil
}

// reduce a component-data with a new data as component as its child
func reduceComponentData(node Node, tok token.Token) (Node, error) {
	return newData(node), nil
}

// reduce a component-repetition to a new repetition with component appended
func reduceComponentRepetition(node Node, tok token.Token) (Node, error) {
	r := newRepetition()
	r.List.append(node)
	return r, nil
}

// reduce a repetition-data with a new data as repetition as its child
func reduceRepetitionData(node Node, tok token.Token) (Node, error) {
	return newData(node), nil
}

// reduce a repetition-text to the repetition with the text appended.
// if the last node in the list of the repetition node is a component,
// then append the text to that component node instead.
// this makes it possible to parse things like:
//   hey~there^to~you
// when we get to "you", we will have this parse tree:
//   [repetition]
//    /         \
//  [component]  [component]
// and the "you" belongs in the second component
func reduceRepetitionText(node Node, tok token.Token) (Node, error) {
	r, _ := node.(*RepetitionNode)

	// must have at least two nodes in our repetition for
	// these conditions to happen
	if len(r.List.Nodes) > 1 {
		lastnode := r.List.Nodes[len(r.List.Nodes)-1]
		switch t := lastnode.(type) {
		// our last node in the repetition is a component.
		// so just append our text node to it.
		case *ComponentNode:
			t.List.append(newText(tok.Val))
			return r, nil
		// the last node in the repetition is nil, which
		// means we had just encountered another repetition delimiter
		// when we already had a repetition. this is done in
		// reduceRepetitionRepetition. we simply just replace
		// the last element with our text node
		case nil:
			r.List.Nodes[len(r.List.Nodes)-1] = newText(tok.Val)
			return r, nil
		}
	}

	r.List.append(newText(tok.Val))
	return r, nil
}

// reduce a repetition-component to the repetition with the component appended.
// if the last element of the repetition is a text, then we just convert it to
// a component and append the text to it.
func reduceRepetitionComponent(node Node, tok token.Token) (Node, error) {
	r, _ := node.(*RepetitionNode)
	lastnode := r.List.Nodes[len(r.List.Nodes)-1]
	switch lastnode.(type) {
	case *TextNode:
		c := newComponent()
		c.List.append(lastnode)
		r.List.Nodes[len(r.List.Nodes)-1] = c
	}

	return r, nil
}

// reudce a repetition-repetition to the repetition with a nil appended.
// this effectively just expands the list by one, and consecutive tokens
// will replace it. this case is handled in reduceRepetitionText()
func reduceRepetitionRepetition(node Node, tok token.Token) (Node, error) {
	r, _ := node.(*RepetitionNode)
	r.List.append(nil)
	return r, nil
}

// next returns the next token.
func (t *Tree) next() token.Token {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextToken()
	}
	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// // backup2 backs the input stream up two tokens.
// // The zeroth token is already there.
// func (t *Tree) backup2(t1 item) {
//   t.token[1] = t1
//   t.peekCount = 2
// }

// // backup3 backs the input stream up three tokens
// // The zeroth token is already there.
// func (t *Tree) backup3(t2, t1 item) { // Reverse order: we're pushing back.
//   t.token[1] = t1
//   t.token[2] = t2
//   t.peekCount = 3
// }

// peek returns but does not consume the next token.
func (t *Tree) peek() token.Token {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}
	t.peekCount = 1
	t.token[0] = t.lex.nextToken()
	return t.token[0]
}

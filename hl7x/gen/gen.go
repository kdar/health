package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/pkg/inflect"

	"github.com/jteeuwen/go-pkg-xmlx"
)

func main() {
	fp, err := os.Open("vendor/2.3/datatypes.xsd")
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	doc := xmlx.New()
	err = doc.LoadStream(fp, nil)
	if err != nil {
		log.Fatal(err)
	}

	schema := doc.SelectNode("*", "schema")

	// build map of types
	typs := make(map[string]*Node)
	nodes1 := schema.SelectNodes("*", "complexType")
	nodes2 := schema.SelectNodes("*", "simpleType")
	for _, node := range append(nodes1, nodes2...) {
		typs[node.As("*", "name")] = NewNode(node)
	}

	// build map of elements
	elements := make(map[string]*Node)
	nodes1 = schema.SelectNodes("*", "element")
	for _, node := range nodes1 {
		elements[node.As("*", "name")] = NewNode(node)
	}

	rs := inflect.NewRuleset()
	for k, v := range typs {
		//FIXME: REMOVE
		//if k != "XPN" {
		//	continue
		//}

		if !strings.HasSuffix(k, "CONTENT") && v.Name.Local == "complexType" {
			typeName := NormalizeTypeName(rs, k)
			fileName := NormalizeFileName(rs, k)

			fmt.Println(fileName + ".go:")
			fmt.Printf("type %s struct {\n", typeName)

			sequence := v.SelectNode("*", "sequence")
			if sequence == nil {
				continue
			}

			els := sequence.SelectNodes("*", "element")
			if els == nil {
				continue
			}

			//Namings.normalize_name(Meta.type_desc(tp) || Meta.ref(el_ref) || Meta.name(el_ref), Meta.ref(el_ref) || Meta.name(el_ref)),
			//Namings.mk_class_name(Meta.base_type(tp) || Meta.name(tp).split('.').first),

			for _, el := range els {
				refStr := el.As("*", "ref")
				typStr := elements[refStr].As("*", "type")
				typ := typs[typStr]

				doc := typ.FindNode("annotation", "documentation")
				if doc != nil {
					fmt.Printf("  // %s\n", doc.S("*", "documentation"))
				}
			}

			fmt.Printf("}\n\n")
		}
	}
}

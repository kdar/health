package edifact

// import (
//   "fmt"
//   //"github.com/davecgh/go-spew/spew"
//   "regexp"
// )

// //`UNA:+./*'UIB+UNOA:0++:0+++Per-Se:ZZZ+Per-Se:ZZZ+20130222:065442'UIH+SCRIPT:008:001:ERROR'STS+900+007+Missing 'Request=''UIT++3'UIZ++1'`),

// // just doing this for fun to parse it in a easier way
// func Unmarshal2(d []byte) (Values, error) {
//   var ret Values
//   //componentD := d[3]
//   dataD := d[4]
//   //decimal := d[5]
//   //release := d[6]
//   //repetitionD := d[7]
//   segmentT := d[8]

//   d = d[8:]
//   if d[len(d)-1] == segmentT {
//     d = d[:len(d)-1]
//   }

//   regex := fmt.Sprintf("%s([A-Z]{3})%s", regexp.QuoteMeta(string(segmentT)), regexp.QuoteMeta(string(dataD)))
//   segmentr := regexp.MustCompile(regex)
//   regex = fmt.Sprintf("(%s)", regexp.QuoteMeta(string(dataD)))
//   datar := regexp.MustCompile(regex)
//   //regex = fmt.Sprintf("(%s)", regexp.QuoteMeta(string(componentD)))
//   //componentr := regexp.MustCompile(regex)

//   segments := split(segmentr, d, 0)
//   for sn, segmentData := range segments {
//     ret = append(ret, Values{})

//     datas := split(datar, segmentData, 1)
//     for dn, data := range datas {
//       ret[sn] = append(ret[sn].(Values), Values{})
//       ret[sn].(Values)[dn] = append(ret[sn].(Values)[dn].(Values), string(data))
//       //components := split(componentr, data, 1)

//     }
//   }

//   // datas := split(r, []byte("UIH+SCRIPT:008:001:ERROR"))
//   // for _, d := range datas {
//   //   fmt.Println(string(d))
//   // }

//   return append(Values{Header{"UNA", ":+./*'"}}, ret...), nil
// }

// func split(r *regexp.Regexp, data []byte, skip int) [][]byte {
//   var ret [][]byte
//   indexes := r.FindAllSubmatchIndex(data, -1)
//   if indexes != nil {
//     ret = append(ret, data[0:indexes[0][2]])
//   }
//   //fmt.Println(ret)
//   for n, x := range indexes {
//     end := len(data)
//     if n+1 < len(indexes) {
//       end = indexes[n+1][0]
//     }

//     ret = append(ret, data[x[2]+skip:end])
//     //fmt.Println(string(data[x[2]+skip : end]))
//   }
//   return ret
// }

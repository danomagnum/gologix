package l5x

import (
	"fmt"
	"strconv"
)

func LoadTags(l5xData RSLogix5000Content) (map[string]any, error) {

	var err error
	tags := make(map[string]any)

	for _, tag := range l5xData.Controller.Tags.Tag {
		if tag.Data != nil {
			if len(tag.Data) > 1 {
				if tag.Data[1].DataValue != nil {
					tags[tag.NameAttr], err = L5xTypeToGoType(tag.DataTypeAttr, tag.Data[1].DataValue.ValueAttr)
				}
			}
		}
		if err != nil {
			return nil, fmt.Errorf("error converting %s: %s", tag.NameAttr, err)
		}
	}

	for _, program := range l5xData.Controller.Programs.Program {
		progprefix := fmt.Sprintf("program:%s", program.NameAttr)
		progtags := make(map[string]any)
		for _, tag := range program.Tags.Tag {
			tagname := tag.NameAttr
			tagtype := tag.DataTypeAttr
			if len(tag.Data) > 1 {
				// could be an atomic value
				if tag.Data[1].DataValue != nil {
					progtags[tagname], err = L5xTypeToGoType(tagtype, tag.Data[1].DataValue.ValueAttr)
					if err != nil {
						return nil, fmt.Errorf("error converting %s.%s: %s", progprefix, tagname, err)
					}
					continue
				}
				// or a structure
				if tag.Data[1].Structure != nil {
					mapValue := make(map[string]any)
					for _, member := range tag.Data[1].Structure.DataValueMember {
						mapValue[member.NameAttr], err = L5xTypeToGoType(member.DataTypeAttr, member.ValueAttr)
						if err != nil {
							return nil, fmt.Errorf("error converting %s.%s.%s: %s", progprefix, tagname, member.NameAttr, err)
						}
					}
					progtags[tagname] = mapValue
					continue
				}
				// or an array of atomic values or structures
				if tag.Data[1].Array != nil {

					dims, err := strconv.ParseInt(tag.DimensionsAttr, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("%s.%s invalid dimensions on %s: %s", progprefix, tagname, tag.DimensionsAttr, err)
					}
					arrValue := make([]any, 0, dims)
					for _, element := range tag.Data[1].Array.Element {
						if element.Structure != nil {
							// this is an array of structures
							myValue := make(map[string]any)
							for _, member := range element.Structure[0].DataValueMember {
								if member.DataTypeAttr == "STRING" {
									myValue[member.NameAttr], err = L5xTypeToGoType(member.DataTypeAttr, member.CData())
									if err != nil {
										return nil, fmt.Errorf("error converting %s: %s", member.NameAttr, err)
									}
									continue
								}
								myValue[member.NameAttr], err = L5xTypeToGoType(member.DataTypeAttr, member.ValueAttr)
								if err != nil {
									return nil, fmt.Errorf("error converting %s: %s", member.NameAttr, err)
								}
							}
							arrValue = append(arrValue, myValue)
							continue
						}

						val, err := L5xTypeToGoType(tagtype, element.ValueAttr)
						if err != nil {
							return nil, fmt.Errorf("error converting %s: %s", tagname, err)
						}
						arrValue = append(arrValue, val)
					}
					progtags[tagname] = arrValue
				}
			}
		}
		tags[progprefix] = progtags

	}
	return tags, nil

}

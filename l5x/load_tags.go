package l5x

import (
	"fmt"
	"strconv"
	"strings"
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

func LoadTagComments(l5xData RSLogix5000Content) (map[string]string, error) {
	dtm := GetDataTypeMap(l5xData)

	var err error
	tags := make(map[string]string)

	for _, module := range l5xData.Controller.Modules.Module {
		if module.Communications == nil {
			continue
		}
		if module.Communications.Connections == nil {
			continue
		}
		if module.Communications.Connections.Connection == nil {
			continue
		}
		for _, connection := range module.Communications.Connections.Connection {
			in := connection.InputTag
			if in != nil {
				if in.Comments != nil {
					for _, commentGroup := range in.Comments {
						for _, comment := range commentGroup.Comment {
							name := fmt.Sprintf("%s:I%s", module.NameAttr, comment.OperandAttr)
							c := comment.CData()
							if c == "" {
								continue
							}
							tags[name] = c
						}
					}
				}
			}
			out := connection.OutputTag
			if out != nil {
				if out.Comments != nil {
					for _, commentGroup := range out.Comments {
						for _, comment := range commentGroup.Comment {
							name := fmt.Sprintf("%s:O%s", module.NameAttr, comment.OperandAttr)
							c := comment.CData()
							if c == "" {
								continue
							}
							tags[name] = c
						}
					}
				}
			}

		}
	}

	for _, tag := range l5xData.Controller.Tags.Tag {

		mydesc := ""
		if tag.Data != nil {
			if len(tag.Description) > 0 {
				mydesc = tag.Description[0].CData()
				tags[tag.NameAttr] = mydesc
			}
		}
		typeComments := GetTypeComments(dtm, tag.DataTypeAttr)
		for typeName, typeComment := range typeComments {
			if typeComment != "" {
				if mydesc != "" {
					typeComment = mydesc + " " + typeComment
				}
				tags[fmt.Sprintf("%s.%s", tag.NameAttr, typeName)] = typeComment
			}
		}
	}

	for _, program := range l5xData.Controller.Programs.Program {
		progprefix := fmt.Sprintf("program:%s", program.NameAttr)
		progtags := make(map[string]any)
		for _, tag := range program.Tags.Tag {
			tagname := tag.NameAttr
			if len(tag.Description) > 1 {
				comment := ""
				for d := range tag.Description {
					comment = comment + tag.Description[d].CData()
				}
				progtags[tagname] = comment
				if err != nil {
					return nil, fmt.Errorf("error converting %s.%s: %s", progprefix, tagname, err)
				}
				continue
			}
		}
	}
	return tags, nil

}

func LoadRungComments(l5xData RSLogix5000Content) (map[string]string, error) {
	comments := make(map[string]string)

	taskLookup := make(map[string]string)
	tasks := l5xData.Controller.Tasks.Task
	for _, task := range tasks {
		if task.ScheduledPrograms == nil {
			continue
		}
		for _, program := range task.ScheduledPrograms.ScheduledProgram {
			taskLookup[program.NameAttr] = task.NameAttr
		}
	}

	for _, program := range l5xData.Controller.Programs.Program {
		for _, routine := range program.Routines.Routine {
			for _, rungContent := range routine.RLLContent {
				for _, rung := range rungContent.Rung {
					for _, comment := range rung.Comment {
						c := comment.CData()
						if c == "" {
							continue
						}
						task := taskLookup[program.NameAttr]
						location := fmt.Sprintf("%s/%s/%s[%s]", task, program.NameAttr, routine.NameAttr, rung.NumberAttr)
						comments[location] = c
					}
				}
			}
		}
	}
	return comments, nil
}

// convert list of data types to a map[string]any
func GetDataTypeMap(l5xData RSLogix5000Content) map[string]*DataTypeType {
	dataTypes := make(map[string]*DataTypeType)

	for _, dataType := range l5xData.Controller.DataTypes.DataType {
		dataTypes[dataType.NameAttr] = dataType
	}

	return dataTypes
}

func GetTypeComments(types map[string]*DataTypeType, typeName string) map[string]string {
	dt, ok := types[typeName]
	if !ok {
		return nil
	}
	out := make(map[string]string)
	myDescription := ""
	if dt.Description != nil {
		myDescription = dt.Description.CData()
	}
	out[typeName] = myDescription
	members := dt.Members
	if members == nil {
	}
	for _, member := range dt.Members.Member {
		//typeName := fmt.Sprintf("%s.%s", typeName, member.NameAttr)
		memberDesc := myDescription
		if member.Description != nil {
			if myDescription != "" {
				memberDesc = myDescription + "$" + member.Description.CData()
			} else {
				memberDesc = member.Description.CData()
			}
		}
		name := member.NameAttr
		//out[member.NameAttr] = memberDesc
		submembers := GetTypeComments(types, member.DataTypeAttr)
		if member.DataTypeAttr == "BIT" {
			// bit number is BitNumberAttr + TargetAttr[-2:] - 1 ?
			offset_str := member.TargetAttr[len(member.TargetAttr)-2:]
			offset_str = strings.TrimLeft(offset_str, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
			num, err := strconv.Atoi(offset_str)
			if err != nil {
				continue // skip this member
			}
			pos := member.BitNumberAttr + num
			// This is assinine.
			switch num {
			case 9:
				pos--
			case 18:
				pos -= 2
			case 27:
				pos -= 3
			}
			name = fmt.Sprintf("%s<%d>", member.NameAttr, pos)
		}
		out[name] = memberDesc
		for subName, subDesc := range submembers {
			fullname := fmt.Sprintf("%s.%s", member.NameAttr, subName)
			if memberDesc != "" {
				out[fullname] = memberDesc + "." + subDesc
			} else {
				out[fullname] = subDesc
			}
		}
	}
	return out
}

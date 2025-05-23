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

func LoadTagComments(l5xData RSLogix5000Content) (map[string]any, error) {

	var err error
	tags := make(map[string]any)

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
		if tag.Data != nil {
			if len(tag.Description) > 0 {
				tags[tag.NameAttr] = tag.Description[0].CData()
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

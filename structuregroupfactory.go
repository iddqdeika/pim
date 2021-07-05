package pim

import "fmt"

var StructureGroupUpdates StructureGroupUpdateFactory = &structureGroupUpdateFactory{}

type structureGroupUpdateFactory struct {
}

func (s *structureGroupUpdateFactory) NewUpdateOrder(columns []string, groups ...StructureGroupUpdate) (*PimUpdateOrder, error) {
	ub, err := newStructureGroupUpdate(columns, groups)
	if err != nil {
		return nil, err
	}
	return &PimUpdateOrder{
		UrlPath:    StructureGroupPath,
		UpdateBody: ub,
	}, nil
}

func newStructureGroupUpdate(columns []string, groups []StructureGroupUpdate) (*PimUpdateBody, error) {
	ub := &PimUpdateBody{}
	cm := make(map[string]int)
	for i, column := range columns {
		cm[column] = i
		ub.Columns = append(ub.Columns, PimUpdateColumn{Identifier: column})
	}
	for _, group := range groups {
		row := PimUpdateRow{
			Object: PimUpdateObject{
				ID: "'" + group.Identifier + "'@" + group.StructureID,
			},
			Values: make([]string, len(cm)),
		}
		if len(group.Fields) != len(cm) {
			return nil, fmt.Errorf("group %v contains %v fields, but there must be %v",
				group.Identifier, len(group.Fields), len(cm))
		}
		for field, val := range group.Fields {
			i, ok := cm[field]
			if !ok {
				return nil, fmt.Errorf("group %v contains field %v that isn't declared in fields slice",
					group.Identifier, field)
			}
			row.Values[i] = val
		}
		ub.Rows = append(ub.Rows, row)
	}
	return ub, nil
}

type StructureGroupUpdate struct {
	Identifier  string
	StructureID string
	Fields      map[string]string
}

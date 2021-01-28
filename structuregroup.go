package pim

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	InfomodelPath = "StructureGroup/StructureGroupAttribute"
)

var (
	//fields for requests
	InfomodelFields = []string{"StructureGroupAttributeLang.Name(Russian)", "StructureGroupAttribute.Datatype",
		"StructureGroupAttributeLang.DomainValue(Russian)", "StructureGroupAttribute.IsMandatory",
		"StructureGroupAttribute.MultiValue"}
)

type structureGroupProvider struct {
	c *Client
}

func (i *structureGroupProvider) GetInfomodelByIdentifier(identifier string, structureID int) (*Infomodel, error) {
	url := i.c.baseListUrl() + InfomodelPath + "/byItems?" +
		"items=" + "'" + identifier + "'@" + strconv.Itoa(structureID) +
		"&fields=" + strings.Join(InfomodelFields, ",") +
		"&pageSize=-1"
	res, err := i.c.get(url)
	if err != nil {
		return nil, err
	}
	fs := make(map[string]Feature)
	for _, row := range res.Rows {
		if len(row.Values) != len(InfomodelFields) {
			return nil, fmt.Errorf("cant parse infomodel, wrong num of values in a row")
		}
		name, ok := row.Values[0].(string)
		if !ok {
			return nil, TypeCastErr
		}
		dataType, ok := row.Values[1].(string)
		if !ok {
			return nil, TypeCastErr
		}
		pi, ok := row.Values[2].([]interface{})
		if !ok {
			return nil, TypeCastErr
		}
		presets := make([]string, 0)
		for _, vi := range pi {
			val, ok := vi.(string)
			if !ok {
				return nil, TypeCastErr
			}
			if val == "" {
				continue
			}
			presets = append(presets, val)
		}
		manda, ok := row.Values[3].(string)
		if !ok {
			return nil, TypeCastErr
		}
		multi, ok := row.Values[4].(string)
		if !ok {
			return nil, TypeCastErr
		}
		fs[name] = Feature{
			Name:         name,
			DataType:     dataType,
			PresetValues: presets,
			Mandatory:    manda == "true",
			Multivalued:  multi == "true",
		}
	}
	return &Infomodel{
		Identifier:  identifier,
		StructureID: structureID,
		Features:    fs,
	}, nil
}

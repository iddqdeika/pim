package pim

import "strings"

var searchFabricInstance SearchFabric = &searchFabric{}

func GetSearchFabricInstance() SearchFabric {
	return searchFabricInstance
}

type searchFabric struct {
}

func (s *searchFabric) NewSearch(reportPath string) Search {
	return &search{
		reportPath: reportPath,
	}
}

type search struct {
	reportPath string
	predicates []SearchPredicate
	fields     []string
}

func (s *search) ReportPath() string {
	return s.reportPath
}

func (s *search) Query() string {
	preds := make([]string, len(s.predicates))
	for i, predicate := range s.predicates {
		preds[i] = predicate.Render()
	}
	return strings.Join(preds, " and ")
}

func (s *search) Fields() string {
	return strings.Join(s.fields, ";")
}

func (s *search) WithPredicate(predicate SearchPredicate) Search {
	s.predicates = append(s.predicates, predicate)
	return s
}

func (s *search) WithOutputField(field string) Search {
	s.fields = append(s.fields, field)
	return s
}

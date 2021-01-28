package pim

var searchPredicateFabricInstance SearchPredicateFabric = &searchPredicateFabric{}

func GetSearchPredicateFabricInstance() SearchPredicateFabric {
	return searchPredicateFabricInstance
}

type searchPredicateFabric struct {
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewEquals(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewContains(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewContainsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

func (s *searchPredicateFabric) Or(p1, p2 SearchPredicate) SearchPredicate {
	return simplePredicate{
		v1:       p1.Render(),
		operator: "or",
		v2:       p2.Render(),
	}
}

type simplePredicate struct {
	v1       string
	operator string
	v2       string
}

func (e simplePredicate) Render() string {
	return "(" + e.v1 + " " + e.operator + " " + e.v2 + ")"
}

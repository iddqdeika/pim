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
func (s *searchPredicateFabric) NewNotEquals(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

func (s *searchPredicateFabric) NewIsEmpty(field string) SearchPredicate {
	return simplePredicate{
		exclude:  false,
		v1:       field,
		operator: "is",
		v2:       "empty",
	}
}

func (s *searchPredicateFabric) NewNotIsEmpty(field string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "is",
		v2:       "empty",
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewEqualsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equalsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewNotEqualsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "equalsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewContains(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "contains",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewNotContains(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "contains",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewContainsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "containsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFabric) NewNotContainsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "containsIC",
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

func (s *searchPredicateFabric) And(p1, p2 SearchPredicate) SearchPredicate {
	return simplePredicate{
		v1:       p1.Render(),
		operator: "and",
		v2:       p2.Render(),
	}
}

type simplePredicate struct {
	exclude  bool
	v1       string
	operator string
	v2       string
}

func (e simplePredicate) Render() string {
	prefix := ""
	if e.exclude {
		prefix = "not "
	}
	return "(" + prefix + e.v1 + " " + e.operator + " " + e.v2 + ")"
}

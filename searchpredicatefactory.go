package pim

var SearchPredicates SearchPredicateFactory = &searchPredicateFactory{}

type searchPredicateFactory struct {
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) Equals(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) NotEquals(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "equals",
		v2:       value,
	}
}

func (s *searchPredicateFactory) IsEmpty(field string) SearchPredicate {
	return simplePredicate{
		exclude:  false,
		v1:       field,
		operator: "is",
		v2:       "empty",
	}
}

func (s *searchPredicateFactory) NotIsEmpty(field string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "is",
		v2:       "empty",
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) EqualsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "equalsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) NotEqualsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "equalsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) Contains(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "contains",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) NotContains(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "contains",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) ContainsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		v1:       field,
		operator: "containsIC",
		v2:       value,
	}
}

// учитывайте, что в значении числа надо указывать без кавычек, а строки и даты - в кавычках
func (s *searchPredicateFactory) NotContainsIC(field string, value string) SearchPredicate {
	return simplePredicate{
		exclude:  true,
		v1:       field,
		operator: "containsIC",
		v2:       value,
	}
}

func (s *searchPredicateFactory) Or(p1, p2 SearchPredicate) SearchPredicate {
	return simplePredicate{
		v1:       p1.Render(),
		operator: "or",
		v2:       p2.Render(),
	}
}

func (s *searchPredicateFactory) And(p1, p2 SearchPredicate) SearchPredicate {
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

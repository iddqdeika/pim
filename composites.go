package pim

var Updates updates = updates{
	Article:        &articleUpdateFactory{},
	StructureGroup: &structureGroupUpdateFactory{},
}

type updates struct {
	Article        ArticleUpdateFactory
	StructureGroup StructureGroupUpdateFactory
}

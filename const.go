package pim

import "time"

const (
	restUrlPath     = "/rest/V1.0/"
	restUrlListPath = "/rest/V1.0/list/"
	retryTimeout    = time.Second

	//default mandatory (0) is Mandatory, be careful
	UpdatePolicyMandatory UpdatePolicy = 0
	UpdatePolicyNormal    UpdatePolicy = 1

	ArticlePath           = "Article"
	ArticleAttributesPath = "Article/ArticleAttribute"
	InfomodelPath         = "StructureGroup/StructureGroupAttribute"
	StructureGroupPath    = "StructureGroup"
)

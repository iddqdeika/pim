package pim

import "time"

const (
	restUrlPath     = "/rest/V1.0/"
	restUrlListPath = "/rest/V1.0/list/"
	retryTimeout    = time.Second

	//default mandatory (0) is Mandatory, be careful
	UpdatePolicyMandatory UpdatePolicy = 0
	UpdatePolicyNormal    UpdatePolicy = 1

	ArticlePath              = "Article"
	ArticleAttributesPath    = "Article/ArticleAttribute"
	ArticleMediaAssetMapPath = "Article/ArticleMediaAssetMap"
	InfomodelPath            = "StructureGroup/StructureGroupAttribute"
	StructureGroupPath       = "StructureGroup"

	AssetMediaTypeDocID1 = iota - 3
	AssetMediaTypeDocID2
	AssetMediaTypeDocID3
	AssetMediaTypeImageID1
	AssetMediaTypeImageID2
	AssetMediaTypeImageID3
	AssetMediaTypeImageID4
	AssetMediaTypeImageID5
	AssetMediaTypeImageID6
	AssetMediaTypeImageID7
	AssetMediaTypeImageID8
	AssetMediaTypeImageID9
	AssetMediaTypeImageID10
)

func IsValidAssetMediaTypeID(id int) bool {
	if id < AssetMediaTypeDocID1 || id > AssetMediaTypeImageID10 {
		return false
	}

	return true
}

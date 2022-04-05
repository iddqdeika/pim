package pim

var Assets AssetFactory = &assetFactory{}

type assetFactory struct{}

type Asset struct {
	ID     string // id, который присваевается после register
	TempID string // id, который присваевается после upload
	Name   string
	Data   []byte
	Type   string
}

func (f *assetFactory) NewAssetFromData(name string, assetType string, data []byte) *Asset {
	return &Asset{Name: name, Type: assetType, Data: data}
}

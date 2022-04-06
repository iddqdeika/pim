package pim

var Assets AssetFactory = &assetFactory{}

type assetFactory struct{}

type File struct {
	Name string
	Data []byte
	Type string
}

type UploadedFile struct {
	ID   string // временный id, который присваевается после uploade
	Name string
	Type string
}

type Asset struct {
	ID   string // id, который присваевается после register
	Name string
	Type string
}

func (f *assetFactory) NewFile(name string, filetype string, data []byte) *File {
	return &File{Name: name, Type: filetype, Data: data}
}

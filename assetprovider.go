package pim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type assetProvider struct {
	c *Client
}

type uploadResponse struct {
	ID string `json:"id"`
}

func newAssetProvider(c *Client) AssetProvider {
	return &assetProvider{c: c}
}

// UploadAssetData загружает данные ассета, возвращает указатель на Asset с заполненным TempID
func (p *assetProvider) UploadFile(file *File) (*UploadedFile, error) {
	url := p.c.baseUrl() + "manage/file?originalFilename=" + file.Name
	resp, err := p.c.postOctet(url, file.Data)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("server returned unexpected code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	r := &uploadResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, fmt.Errorf("response parsing error: %w", err)
	}

	return &UploadedFile{Name: file.Name, Type: file.Type, ID: r.ID}, nil
}

func (p *assetProvider) RegisterAsset(file *UploadedFile) (*Asset, error) {
	url := p.c.baseUrl() + "media/" + file.ID
	resp, err := p.c.postText(url, file.ID)
	if err != nil {
		return nil, err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("server returned unexpected code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	return &Asset{Name: file.Name, Type: file.Type, ID: string(body)}, nil
}

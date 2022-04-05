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
func (p *assetProvider) UploadAssetData(asset *Asset) (*Asset, error) {
	url := p.c.baseUrl() + "manage/file?originalFilename=" + asset.Name
	resp, err := p.c.postOctet(url, asset.Data)
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

	asset.TempID = r.ID

	return asset, nil
}

func (p *assetProvider) RegisterAsset(asset *Asset) (*Asset, error) {
	url := p.c.baseUrl() + "media/" + asset.TempID
	resp, err := p.c.postText(url, asset.TempID)
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

	asset.ID = string(body)

	return asset, nil
}

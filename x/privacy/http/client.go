package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type Client struct {
	host string
}

func NewClient() *Client {
	return &Client{host: Host}
}

func (c *Client) sendRequest(method, url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) OutputCoin(index string) (*types.OutputCoin, error) {
	url := c.host + "/output_coin/" + index
	body, err := c.sendRequest("GET", url)
	if err != nil {
		return nil, err
	}
	type Result struct {
		OutputCoin *types.OutputCoin `json:"outputCoin"`
	}
	var temp Result
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	if temp.OutputCoin == nil {
		return nil, fmt.Errorf("cannot find outputCoin by index")
	}

	return temp.OutputCoin, nil
}

func (c *Client) AllOutputCoin() ([]types.OutputCoin, error) {
	url := c.host + "/output_coin"
	body, err := c.sendRequest("GET", url)
	if err != nil {
		return nil, err
	}
	type Result struct {
		OutputCoin []types.OutputCoin `json:"outputCoin"`
	}
	var temp Result
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}

	return temp.OutputCoin, nil
}

func (c *Client) OutputCoinLength() (*types.OutputCoinLength, error) {
	url := c.host + "/output_coin_length"
	body, err := c.sendRequest("GET", url)
	if err != nil {
		return nil, err
	}
	type Result struct {
		OutputCoinSerialNumber *types.OutputCoinLength `json:"OutputCoinLength"`
	}
	var temp Result
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	if temp.OutputCoinSerialNumber == nil {
		return nil, fmt.Errorf("Cannot find output coin serialNumber")
	}

	return temp.OutputCoinSerialNumber, nil
}

func (c *Client) OtaCoin(index string) (*types.OTACoin, error) {
	url := c.host + "/ota_coin/" + index
	body, err := c.sendRequest("GET", url)
	if err != nil {
		return nil, err
	}
	type Result struct {
		OTACoin *types.OTACoin `json:"oTACoin"`
	}
	var temp Result
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	if temp.OTACoin == nil {
		return nil, fmt.Errorf("Cannot find ota coin")
	}

	return temp.OTACoin, nil
}

func (c *Client) AllOtaCoins() ([]types.OTACoin, error) {
	url := c.host + "/ota_coin"
	body, err := c.sendRequest("GET", url)
	if err != nil {
		return nil, err
	}
	type Result struct {
		OTACoin []types.OTACoin `json:"oTACoin"`
	}
	var temp Result
	err = json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	if temp.OTACoin == nil {
		return nil, fmt.Errorf("Cannot find ota coin")
	}

	return temp.OTACoin, nil
}

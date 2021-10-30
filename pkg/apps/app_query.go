package apps

import (
	"encoding/json"
)

func (a *Apps) GetApps() ([]App, error) {
	resp, err := a.client.Request("GET", "admin/apps", nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	type Data struct {
		Data []App `json:"data"`
	}
	d := Data{}
	err = json.NewDecoder(resp.Body).Decode(&d)

	if err != nil {
		return nil, err
	}

	return d.Data, err
}

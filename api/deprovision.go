package api

type DeprovisionStatus struct {
	Status string `json:"-"`

	Operation string `json:"operation"`
}

func (c *Client) Deprovision(id string) (*DeprovisionStatus, error) {
	res, err := c.del("/v2/service_instances/" + id)
	if err != nil {
		return nil, err
	}

	var status DeprovisionStatus

	switch res.StatusCode {
	case 410:
		fallthrough
	case 200:
		status.Status = "deprovisioned"
		return &status, nil

	case 202:
		status.Status = "deprovisioning"
		return &status, c.parse(res, &status)
	}

	return nil, c.err(res)
}

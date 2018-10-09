package api

import (
	"fmt"
)

type Catalog struct {
	Services []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`

		Tags     []string `json:"tags,omitempty"`
		Requires []string `json:"requires,omitempty"`

		Bindable             bool `json:"bindable"`
		InstancesRetrievable bool `json:"instances_retrievable"`
		BindingsRetrievable  bool `json:"bindings_retrievable"`
		PlanUpdateable       bool `json:"plan_updateable"`

		Metadata interface{} `json:"metadata,omitempty"`

		DashboardClient *struct {
			ID          string `json:"id"`
			Secret      string `json:"secret"`
			RedirectURI string `json:"redirect_uri"`
		} `json:"dashboard_client,omitempty"`

		Plans []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`

			MaybeFree *bool `json:"free"`
			Free      bool  `json:"-"`

			MaybeBindable *bool `json:"bindable"`
			Bindable      bool  `json:"bindable"`

			Metadata interface{} `json:"metadata,omitempty"`

			Schemas []struct {
				ServiceInstance struct {
					Create struct {
						Parameters map[string]interface{} `json:"parameters"`
					} `json:"create"`
					Update struct {
						Parameters map[string]interface{} `json:"parameters"`
					} `json:"update"`
				} `json:"service_instance"`
				ServiceBinding struct {
					Create struct {
						Parameters map[string]interface{} `json:"parameters"`
					} `json:"create"`
				} `json:"service_binding"`
			} `json:"schemas,omitempty"`
		} `json:"plans,omitempty"`
	} `json:"services,omitempty"`
}

func (c *Client) GetCatalog() (*Catalog, error) {
	res, err := c.get("/v2/catalog")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, c.err(res)
	}

	var cat Catalog
	return &cat, c.parse(res, &cat)
}

func (cat Catalog) FindPlan(service, plan string) (string, string, error) {
	idx := -1
	for i, s := range cat.Services {
		if s.ID == service {
			idx = i
			break
		}
	}
	if idx < 0 {
		for i, s := range cat.Services {
			if s.Name == service {
				idx = i
				break
			}
		}
	}

	if idx >= 0 {
		for _, p := range cat.Services[idx].Plans {
			if p.ID == plan {
				return cat.Services[idx].ID, p.ID, nil
			}
		}
		for _, p := range cat.Services[idx].Plans {
			if p.Name == plan {
				return cat.Services[idx].ID, p.ID, nil
			}
		}
	}

	return "", "", fmt.Errorf("no such service / plan: %s / %s", service, plan)
}

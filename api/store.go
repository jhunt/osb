package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type binding struct {
	ID          string                 `yaml:"id"`
	Credentials map[string]interface{} `yaml:"credentials"`
}

type instance struct {
	ID        string `yaml:"id"`
	ServiceID string `yaml:"service_id"`
	PlanID    string `yaml:"plan_id"`

	Bindings []binding `yaml:"bindings"`
}

type broker struct {
	Broker    string     `yaml:"broker"`
	Instances []instance `yaml:"instances"`
}

type Store struct {
	Data []broker `yaml:"data"`
}

var DefaultStorePath string

func init() {
	DefaultStorePath = os.Getenv("HOME") + "/.osbrc"
}

func ReadStore(path string) (*Store, error) {
	if path == "" {
		path = DefaultStorePath
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{}, nil
		}
		return nil, err
	}

	var store Store
	return &store, yaml.Unmarshal(b, &store)
}

func (s *Store) Write(path string) error {
	if path == "" {
		path = DefaultStorePath
	}

	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0666)
}

func (s *Store) AddInstance(url, id, service, plan string) {
	url = strings.TrimSuffix(url, "/")
	inst := instance{
		ID:        id,
		ServiceID: service,
		PlanID:    plan,
	}

	for i, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			s.Data[i].Instances = append(broker.Instances, inst)
			return
		}
	}

	s.Data = append(s.Data, broker{
		Broker:    url,
		Instances: []instance{inst},
	})
}

func (s *Store) RemoveInstance(url, id string) {
	url = strings.TrimSuffix(url, "/")

	for i, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			for j, instance := range broker.Instances {
				if instance.ID == id {
					s.Data[i].Instances = append(broker.Instances[:j], broker.Instances[j+1:]...)
					return
				}
			}
		}
	}
}

func (s *Store) GetInstanceDetails(url, id string) (string, string, error) {
	url = strings.TrimSuffix(url, "/")

	for _, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			for _, instance := range broker.Instances {
				if instance.ID == id {
					return instance.ServiceID, instance.PlanID, nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("service instance '%s' not found", id)
}

func (s *Store) AddBinding(url, id, bid string, creds map[string]interface{}) {
	url = strings.TrimSuffix(url, "/")

	for i, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			for j, instance := range broker.Instances {
				if instance.ID == id {
					s.Data[i].Instances[j].Bindings = append(instance.Bindings, binding{
						ID:          bid,
						Credentials: creds,
					})
					return
				}
			}
		}
	}
}

func (s *Store) GetBindingDetails(url, id string) (string, string, string, error) {
	url = strings.TrimSuffix(url, "/")

	for _, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			for _, instance := range broker.Instances {
				for _, binding := range instance.Bindings {
					if binding.ID == id {
						return instance.ID, instance.ServiceID, instance.PlanID, nil
					}
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("service instance binding '%s' not found", id)
}

func (s *Store) RemoveBinding(url, id, bid string) {
	url = strings.TrimSuffix(url, "/")

	for i, broker := range s.Data {
		if strings.TrimSuffix(broker.Broker, "/") == url {
			for j, instance := range broker.Instances {
				if instance.ID == id {
					for k, binding := range instance.Bindings {
						if binding.ID == bid {
							s.Data[i].Instances[j].Bindings = append(instance.Bindings[:k], instance.Bindings[k+1:]...)
							return
						}
					}
				}
			}
		}
	}
}

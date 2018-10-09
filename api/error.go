package api

import (
	"fmt"
	"net/http"
)

type Error struct {
	HTTP        string
	Description string `json:"description"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (HTTP %s)", e.Description, e.HTTP)
}

func (c *Client) err(res *http.Response) error {
	var e Error
	if err := c.parse(res, &e); err != nil {
		return err
	}
	e.HTTP = res.Status
	if e.Description == "" {
		e.Description = "an unknown error has occurred"
	}

	return e
}

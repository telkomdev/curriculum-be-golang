package services

import (
	"100-ticketing/app/models"
	"net/url"
)

func (c *Services) GetRouteByID(token string, id string) (res models.Route, err error) {
	u, err := url.JoinPath(c.EndpointRoute, "/", id)
	if err != nil {
		return res, err
	}

	err = c.Get(token, u, map[string]string{}, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

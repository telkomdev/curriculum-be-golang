package services

import (
	"100-ticketing/app/models"
	"net/url"
)

func (c *Services) GetUserByID(token string, id string) (res models.User, err error) {
	u, err := url.JoinPath(c.EndpointUser, "/", id)
	if err != nil {
		return res, err
	}

	err = c.Get(token, u, map[string]string{}, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

package services

import (
	"100-ticketing/app/models"
	"net/url"
)

func (c *Services) CreateNewTicket(token string, data models.CreateTicketRequest) (res models.Ticket, err error) {
	u, err := url.JoinPath(c.EndpointTicket)
	if err != nil {
		return res, err
	}

	err = c.Post(token, u, map[string]string{}, data, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

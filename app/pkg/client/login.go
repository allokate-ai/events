package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Login struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Successful bool   `json:"successful"`
}

func (c *Client) EmitLoginEvent(source string, payload Login) (GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(payload)
	if err != nil {
		return GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "login",
		Source:    source,
		Body:      data,
	}

	return c.Publish(e)
}

package board

import (
	"encoding/json"
	"errors"
	"time"
)

type Changeset struct {
	ID      string    `json:"string"`
	Changes []Changes `json:"changes"`
}

type Changes struct {
	Command Command     `json:"command"`
	Payload interface{} `json:"payload"`
}

type Command struct {
	Type     string    `json:"type"`
	Datetime time.Time `json:"datetime"`
}

type PayloadCreateThread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadEditThread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Changes) UnmarshalJSON(data []byte) error { // TODO(thomasmarlow): move to aux?

	// type aliasing to avoid infinite recursion on Changes.UnmarshallJSON
	type Alias Changes
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch c.Command.Type {
	case "createThread": // TODO(thomasmarlow): de-duplicate these lines using generics
		var payloadCreateThread PayloadCreateThread
		if err := json.Unmarshal(data, &payloadCreateThread); err != nil {
			return err
		}
		c.Payload = payloadCreateThread
	case "editThread":
		var payloadEditThread PayloadEditThread
		if err := json.Unmarshal(data, &payloadEditThread); err != nil {
			return err
		}
		c.Payload = payloadEditThread
	default:
		return errors.New("unknown command type") // TODO(thomasmarlow): i'm not being able to achieve this error; besides, implement actual custom error type
	}

	return nil
}

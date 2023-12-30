package board

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type Changeset struct {
	ID      string    `json:"string"`
	Changes []Changes `json:"changes"`
}

type Changes struct {
	Command Command `json:"command"`
	Payload any     `json:"payload"`
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

// TODO(thomasmarlow): move to aux?
// TODO(thomasmarlow): this code is not being ran...
func (c *Changes) UnmarshalJSON(data []byte) error {

	log.Println(1)

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

	log.Println(2)

	payload, validCommand := map[string]any{
		"createThread": PayloadCreateThread{},
		"editThread":   PayloadEditThread{},
	}[c.Command.Type]
	if !validCommand {

		log.Println(3)

		return errors.New("pepelio")
		// if the error is handled now, the code becomes less clear/readable;
		// it is somehow duplicate logic,
		// but this code here is unavoidable,
		// as the unmarshaling is the presentation layer's responsibility,
		// and must not be delegated to the logic layer
	}

	log.Println(4)

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	log.Println(5)

	c.Payload = payload

	return nil
}

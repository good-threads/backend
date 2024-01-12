package board

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
)

type Command struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Datetime mongoClient.NanoTime `json:"datetime"` // TODO(thomasmarlow): unused
	Payload  any                  `json:"payload"`
}

type PayloadCreateThread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadEditThreadName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadHideThread struct {
	ID string `json:"id"`
}

type PayloadRelocateThread struct {
	ID       string `json:"id"`
	NewIndex uint   `json:"newIndex"`
}

type PayloadCreateKnot struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
	KnotBody string `json:"knotBody"`
}

type PayloadEditKnot struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
	KnotBody string `json:"knotBody"`
}

type PayloadDeleteKnot struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
}

// TODO(thomasmarlow): move to aux?
func (c *Command) UnmarshalJSONAltVersion(data []byte) error {
	type PayloadUnmarshalFunc func(data []byte, c *Command) error

	payloadFuncs := map[string]PayloadUnmarshalFunc{
		"createThread":   unmarshalPayload[PayloadCreateThread],
		"editThread":     unmarshalPayload[PayloadEditThreadName],
		"relocateThread": unmarshalPayload[PayloadRelocateThread],
		"hideThread":     unmarshalPayload[PayloadHideThread],
		"createKnot":     unmarshalPayload[PayloadCreateKnot],
		"editKnot":       unmarshalPayload[PayloadEditKnot],
		"deleteKnot":     unmarshalPayload[PayloadDeleteKnot],
	}

	log.Println(1)

	// type aliasing to avoid infinite recursion on Commands.UnmarshallJSON
	type Alias Command
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	log.Println(2)

	payloadFunc, ok := payloadFuncs[c.Type]
	if !ok {
		log.Println(3)
		return errors.New("unknown command type") // TODO(thomasmarlow): proper error handling
		// if the error is handled now, the code becomes less clear/readable;
		// it is somehow duplicate logic,
		// but this code here is unavoidable,
		// as the unmarshaling is the presentation layer's responsibility,
		// and must not be delegated to the logic layer
	}

	err := payloadFunc(data, c)
	if err != nil {
		return err
	}

	log.Println(4)

	return nil
}

func unmarshalPayload[Payload any](data []byte, c *Command) error {
	log.Println(data)
	payload := new(Payload)
	if err := json.Unmarshal(data, payload); err != nil {
		return err
	}
	log.Println(*payload)
	c.Payload = *payload
	return nil
}

func (c *Command) UnmarshalJSONSSSS(data []byte) error { // TODO(thomasmarlow): benchmark and make a choice
	log.Println(1)

	// type aliasing to avoid infinite recursion on Commands.UnmarshallJSON
	type Alias Command
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	log.Println(2)

	payloadTypes := map[string]reflect.Type{
		"createThread":   reflect.TypeOf(PayloadCreateThread{}),
		"editThread":     reflect.TypeOf(PayloadEditThreadName{}),
		"relocateThread": reflect.TypeOf(PayloadRelocateThread{}),
		"hideThread":     reflect.TypeOf(PayloadHideThread{}),
		"createKnot":     reflect.TypeOf(PayloadCreateKnot{}),
		"editKnot":       reflect.TypeOf(PayloadEditKnot{}),
		"deleteKnot":     reflect.TypeOf(PayloadDeleteKnot{}),
	}

	payloadType, ok := payloadTypes[c.Type]
	if !ok {
		log.Println(3)
		return errors.New("unknown command type") // TODO(thomasmarlow): proper error handling
		// if the error is handled now, the code becomes less clear/readable;
		// it is somehow duplicate logic,
		// but this code here is unavoidable,
		// as the unmarshaling is the presentation layer's responsibility,
		// and must not be delegated to the logic layer
	}

	payload := reflect.New(payloadType).Interface()
	if err := json.Unmarshal(data, payload); err != nil {
		return err
	}

	c.Payload = reflect.ValueOf(payload).Elem().Interface()

	log.Println(4)

	return nil
}

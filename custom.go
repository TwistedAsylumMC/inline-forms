package form

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/df-mc/dragonfly/server/player/form"
)

// Custom represents a form that may be sent to a player and has fields that should be filled out by the player that the
// form is sent to.
type Custom struct {
	// Title is the title of the form that is displayed at the very top of the form.
	Title string
	// Buttons is a slice of elements that can be modified by a player. There must be at least one element for the client
	// to render the form.
	Elements []Element
	// Submit is called when the form is closed or if a player pressed the submit button. This is always called after the
	// Submit of every Element.
	Submit func(closed bool)
}

// Element appends an element to the bottom of the form.
func (form *Custom) Element(element Element) {
	form.Elements = append(form.Elements, element)
}

// SubmitJSON ...
func (form *Custom) SubmitJSON(data []byte, _ form.Submitter) error {
	if data == nil {
		if form.Submit != nil {
			form.Submit(true)
		}
		return nil
	}
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	var inputData []any
	if err := dec.Decode(&inputData); err != nil {
		return fmt.Errorf("error decoding JSON data to slice: %w", err)
	} else if len(form.Elements) != len(inputData) {
		return fmt.Errorf("form JSON data array does not have enough values")
	}
	for i, element := range form.Elements {
		err := element.submit(inputData[i])
		if err != nil {
			return fmt.Errorf("error parsing form response value: %w", err)
		}
	}
	if form.Submit != nil {
		form.Submit(false)
	}
	return nil
}

// MarshalJSON ...
func (form *Custom) MarshalJSON() ([]byte, error) {
	if len(form.Elements) == 0 {
		return nil, errors.New("menu form requires at least one element")
	}
	return json.Marshal(map[string]any{
		"type":    "custom_form",
		"title":   form.Title,
		"content": form.Elements,
	})
}

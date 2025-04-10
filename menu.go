package form

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

// Menu represents a menu form. These menus are made up of a title and a body, with a number of buttons which
// come below the body. These buttons may also have images on the side of them.
type Menu struct {
	// Title is the title of the form that is displayed at the very top of the form.
	Title string
	// Content is the content that is displayed underneath the title and before any buttons.
	Content string
	// Elements is a slice of elements that can be displayed in the form. These elements currently include
	// Header, Label, Divider and Button. Only buttons can be clicked.
	Elements []MenuElement
	// Submit is called when the form is closed or if a player clicks a button. This is always called after
	// the clicked Button's Submit.
	Submit func(closed bool, tx *world.Tx)

	// buttons is a slice of buttons that are present in the form. This is used to determine which button was
	// clicked when the form is submitted. It is populated when the form is marshalled to JSON.
	buttons []Button
}

// Element appends a MenuElement to the bottom of the form.
func (form *Menu) Element(e MenuElement) {
	form.Elements = append(form.Elements, e)
}

// SubmitJSON ...
func (form *Menu) SubmitJSON(data []byte, _ form.Submitter, tx *world.Tx) error {
	if data == nil {
		if form.Submit != nil {
			form.Submit(true, tx)
		}
		return nil
	}
	var index uint
	err := json.Unmarshal(data, &index)
	if err != nil {
		return fmt.Errorf("cannot parse button index as int: %w", err)
	}
	if index >= uint(len(form.buttons)) {
		return fmt.Errorf("button index points to invalid button: %v (only %v buttons present)", index, len(form.buttons))
	}
	button := form.buttons[index]
	if button.Submit != nil {
		button.Submit(tx)
	}
	if form.Submit != nil {
		form.Submit(false, tx)
	}
	return nil
}

// MarshalJSON ...
func (form *Menu) MarshalJSON() ([]byte, error) {
	form.buttons = nil
	for _, element := range form.Elements {
		if button, ok := element.(Button); ok {
			form.buttons = append(form.buttons, button)
		}
	}
	return json.Marshal(map[string]any{
		"type":     "form",
		"title":    form.Title,
		"content":  form.Content,
		"elements": form.Elements,
	})
}

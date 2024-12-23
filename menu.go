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
	// Buttons is a slice of buttons that can be clicked by a player. There must be at least one button for the client
	// to render the form.
	Buttons []Button
	// Submit is called when the form is closed or if a player clicks a button. This is always called after the clicked
	// Button's Submit.
	Submit func(closed bool)
}

// Button appends a button to the bottom of the form.
func (form *Menu) Button(button Button) {
	form.Buttons = append(form.Buttons, button)
}

// SubmitJSON ...
func (form *Menu) SubmitJSON(data []byte, _ form.Submitter, _ *world.Tx) error {
	if data == nil {
		if form.Submit != nil {
			form.Submit(true)
		}
		return nil
	}
	var index uint
	err := json.Unmarshal(data, &index)
	if err != nil {
		return fmt.Errorf("cannot parse button index as int: %w", err)
	}
	if index >= uint(len(form.Buttons)) {
		return fmt.Errorf("button index points to inexistent button: %v (only %v buttons present)", index, len(form.Buttons))
	}
	button := form.Buttons[index]
	if button.Submit != nil {
		button.Submit()
	}
	if form.Submit != nil {
		form.Submit(false)
	}
	return nil
}

// MarshalJSON ...
func (form *Menu) MarshalJSON() ([]byte, error) {
	if form.Buttons == nil {
		form.Buttons = make([]Button, 0)
	}
	return json.Marshal(map[string]any{
		"type":    "form",
		"title":   form.Title,
		"content": form.Content,
		"buttons": form.Buttons,
	})
}

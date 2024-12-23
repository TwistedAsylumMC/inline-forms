package form

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

// Modal represents a modal form. These forms have a body with text and two buttons at the end, typically one for Yes
// and one for No. These buttons may have custom text, but can, unlike with a Menu form, not have images next to them.
type Modal struct {
	// Title is the title of the form that is displayed at the very top of the form.
	Title string
	// Content is the content that is displayed underneath the title and before any buttons.
	Content string
	// Button1 is the top button in the form.
	Button1 Button
	// Button2 is the bottom button in the form.
	Button2 Button
	// Submit is called when the form is closed or if a player clicks a button. This is always called after the clicked
	// Button's Submit.
	Submit func(closed bool)
}

// SubmitJSON ...
func (form *Modal) SubmitJSON(data []byte, _ form.Submitter, _ *world.Tx) error {
	if data == nil {
		if form.Submit != nil {
			form.Submit(true)
		}
		return nil
	}
	var value bool
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("error parsing JSON as bool: %w", err)
	}
	button := form.Button1
	if !value {
		button = form.Button2
	}
	if button.Submit != nil {
		button.Submit()
	}
	if form.Submit != nil {
		form.Submit(false)
	}
	return nil
}

// MarshalJSON ...
func (form *Modal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "modal",
		"title":   form.Title,
		"content": form.Content,
		"button1": form.Button1.Text,
		"button2": form.Button2.Text,
	})
}

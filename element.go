package form

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"strings"
	"unicode/utf8"
)

// Element represents an element that may be added to a Form. Any of the types in this package that implement
// the element interface may be added to a form before it is sent to a player.
type Element interface {
	json.Marshaler
	submit(value any) error
}

// Label represents a static label on a form. It serves only to display a box of text, and users cannot
// submit values to it.
type Label struct {
	// Text is the text held by the label. The text may contain Minecraft formatting codes.
	Text string
}

// MarshalJSON ...
func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "label",
		"text": l.Text,
	})
}

// Submit ...
func (l Label) submit(any) error {
	return nil
}

// Input represents a text input box element. Submitters may write any text in these boxes with no specific
// length.
type Input struct {
	// Text is the text displayed over the input element. The text may contain Minecraft formatting codes.
	Text string
	// Default is the default value filled out in the input. The user may remove this value and fill out its
	// own text. The text may contain Minecraft formatting codes.
	Default string
	// Placeholder is the text displayed in the input box if it does not contain any text filled out by the
	// user. The text may contain Minecraft formatting codes.
	Placeholder string
	// Submit is called with the value provided by the player whenever they submit the form. If the form is closed, this
	// method is not called. This is always called before the Form's Submit.
	Submit func(text string)
}

// MarshalJSON ...
func (i Input) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":        "input",
		"text":        i.Text,
		"default":     i.Default,
		"placeholder": i.Placeholder,
	})
}

// Submit ...
func (i Input) submit(value any) error {
	if i.Submit == nil {
		return nil
	}
	text, ok := value.(string)
	if !ok {
		return fmt.Errorf("value %v is not allowed for input element", value)
	} else if !utf8.ValidString(text) {
		return fmt.Errorf("value %v is not valid UTF8", value)
	}
	i.Submit(text)
	return nil
}

// Toggle represents an on-off button element. Submitters may either toggle this on or off, which will then
// hold a value of true or false respectively.
type Toggle struct {
	// Text is the text displayed over the toggle element. The text may contain Minecraft formatting codes.
	Text string
	// Default is the default value filled out in the input. The user may remove this value and fill out its
	// own text. The text may contain Minecraft formatting codes.
	Default bool
	// Submit is called with the value provided by the player whenever they submit the form. If the form is closed, this
	// method is not called. This is always called before the Form's Submit.
	Submit func(enabled bool)
}

// MarshalJSON ...
func (t Toggle) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "toggle",
		"text":    t.Text,
		"default": t.Default,
	})
}

// Submit ...
func (t Toggle) submit(value any) error {
	if t.Submit == nil {
		return nil
	}
	enabled, ok := value.(bool)
	if !ok {
		return fmt.Errorf("value %v is not allowed for toggle element", value)
	}
	t.Submit(enabled)
	return nil
}

// Slider represents a slider element. Submitters may move the slider to values within the range of the slider
// to select a value.
type Slider struct {
	// Text is the text displayed over the slider element. The text may contain Minecraft formatting codes.
	Text string
	// Min and Max are used to specify the minimum and maximum range of the slider. A value lower or higher
	// than these values cannot be selected.
	Min, Max float64
	// StepSize is the size that one step of the slider takes up. When set to 1.0 for example, a submitter
	// will be able to select only whole values.
	StepSize float64
	// Default is the default value filled out for the slider.
	Default float64
	// Submit is called with the value provided by the player whenever they submit the form. If the form is closed, this
	// method is not called. This is always called before the Form's Submit.
	Submit func(value float64)
}

// MarshalJSON ...
func (s Slider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "slider",
		"text":    s.Text,
		"min":     s.Min,
		"max":     s.Max,
		"step":    s.StepSize,
		"default": s.Default,
	})
}

// Submit ...
func (s Slider) submit(value any) error {
	if s.Submit == nil {
		return nil
	}
	number, ok := value.(json.Number)
	val, err := number.Float64()
	if !ok || err != nil {
		return fmt.Errorf("value %v is not allowed for toggle element", value)
	} else if val < s.Min || val > s.Max {
		return fmt.Errorf("slider value %v is out of range %v-%v", val, s.Min, s.Max)
	}
	s.Submit(val)
	return nil
}

// Dropdown represents a dropdown which, when clicked, opens a window with the options set in the Options
// field. Submitters may select one of the options.
type Dropdown struct {
	// Text is the text displayed over the dropdown element. The text may contain Minecraft formatting codes.
	Text string
	// Options holds a list of options that a Submitter may select. The order of these options is retained
	// when shown to the submitter of the form.
	Options []string
	// DefaultIndex is the index in the Options slice that is used as default. When sent to a Submitter, the
	// value at this index in the Options slice will be selected.
	DefaultIndex int
	// Submit is called with the value provided by the player whenever they submit the form. If the form is closed, this
	// method is not called. This is always called before the Form's Submit.
	Submit func(index int, option string)
}

// MarshalJSON ...
func (d Dropdown) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "dropdown",
		"text":    d.Text,
		"default": d.DefaultIndex,
		"options": d.Options,
	})
}

// Submit ...
func (d Dropdown) submit(value any) error {
	if d.Submit == nil {
		return nil
	}
	number, ok := value.(json.Number)
	val, err := number.Int64()
	if !ok || err != nil {
		return fmt.Errorf("value %v is not allowed for dropdown element", value)
	}
	if val < 0 || int(val) >= len(d.Options) {
		return fmt.Errorf("dropdown value %v is out of range %v-%v", val, 0, len(d.Options)-1)
	}
	d.Submit(int(val), d.Options[val])
	return nil
}

// StepSlider represents a slider that has a number of options that may be selected. It is essentially a
// combination of a Dropdown and a Slider, looking like a slider but having properties like a dropdown.
type StepSlider Dropdown

// MarshalJSON ...
func (s StepSlider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "step_slider",
		"text":    s.Text,
		"default": s.DefaultIndex,
		"steps":   s.Options,
	})
}

// Submit ...
func (s StepSlider) submit(value any) error {
	if s.Submit == nil {
		return nil
	}
	number, ok := value.(json.Number)
	val, err := number.Int64()
	if !ok || err != nil {
		return fmt.Errorf("value %v is not allowed for step slider element", value)
	}
	if val < 0 || int(val) >= len(s.Options) {
		return fmt.Errorf("dropdown value %v is out of range %v-%v", val, 0, len(s.Options)-1)
	}
	s.Submit(int(val), s.Options[val])
	return nil
}

// Button represents a button added to a Menu or Modal form. The button has text on it and an optional image,
// which may be either retrieved from a website or the local assets of the game.
type Button struct {
	// Text holds the text displayed on the button. It may use Minecraft formatting codes and may have
	// newlines.
	Text string
	// Image holds a path to an image for the button. The Image may either be a URL pointing to an image,
	// such as 'https://someimagewebsite.com/someimage.png', or a path pointing to a local asset, such as
	// 'textures/blocks/grass_carried'.
	Image string
	// Submit is called when a player clicks on the button in a form. This is always called before the Form's Submit.
	Submit func(tx *world.Tx)
}

// MarshalJSON ...
func (b Button) MarshalJSON() ([]byte, error) {
	m := map[string]any{"text": b.Text}
	if b.Image != "" {
		buttonType := "path"
		if strings.HasPrefix(b.Image, "http:") || strings.HasPrefix(b.Image, "https:") {
			buttonType = "url"
		}
		m["image"] = map[string]any{"type": buttonType, "data": b.Image}
	}
	return json.Marshal(m)
}

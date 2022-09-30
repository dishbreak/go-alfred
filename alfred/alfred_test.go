package alfred

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func badOption(sr *ScriptFilterResponse) error {
	return errors.New("boooo!")
}

func TestNewScriptFilterResponse(t *testing.T) {

	sr := NewScriptFilterResponse()
	assert.Empty(t, sr.Items)
}

func TestScriptFilterResponseAddItems(t *testing.T) {
	output := bytes.NewBufferString("")

	sr := NewScriptFilterResponse(ScriptFilterWithOutput(output))

	sr.AddItem(ListItem{
		Title:     "foo",
		Subtitle:  "bar",
		Arg:       "fool",
		Variables: map[string]string{},
	})

	sr.SendFeedback()
	assert.Equal(t, `{"items":[{"title":"foo","subtitle":"bar","arg":"fool","valid":false,"variables":{}}]}
`, output.String())
}

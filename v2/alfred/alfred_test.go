package alfred

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestScriptOutputWithCache(t *testing.T) {
	output := bytes.NewBufferString("")

	sr := NewScriptFilterResponse(ScriptFilterWithOutput(output), ScriptFilterWithCache(60, false))
	sr.AddItem(ListItem{
		Title:     "foo",
		Subtitle:  "bar",
		Arg:       "fool",
		Variables: map[string]string{},
	})

	sr.SendFeedback()
	assert.Equal(t, `{"items":[{"title":"foo","subtitle":"bar","arg":"fool","valid":false,"variables":{}}],"cache":{"seconds":60,"loosereload":false}}
`, output.String())
}

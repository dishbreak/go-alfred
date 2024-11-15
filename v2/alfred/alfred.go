package alfred

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Response interface {
	SetError(error)
	SendFeedback()
}

// ListItem represents a single item in the script filter.
type ListItem struct {
	Title     string            `json:"title"`
	Subtitle  string            `json:"subtitle"`
	Arg       string            `json:"arg"`
	Valid     bool              `json:"valid"`
	Variables map[string]string `json:"variables"`
	Uid       string            `json:"uid,omitempty"`
}

// ScriptFilterResponse represents a list of items that the script filter will provide to Alfred.
type ScriptFilterResponse struct {
	Items  []ListItem   `json:"items"`
	Cache  *CacheConfig `json:"cache,omitempty"`
	output io.Writer
}

type CacheConfig struct {
	Seconds     int  `json:"seconds"`
	LooseReload bool `json:"loosereload"`
}

type ScriptFilterResponseOption func(*ScriptFilterResponse)

func ScriptFilterWithOutput(writer io.Writer) ScriptFilterResponseOption {
	return func(sr *ScriptFilterResponse) {
		sr.output = writer
	}
}

func ScriptFilterWithCache(seconds int, looseReload bool) ScriptFilterResponseOption {
	return func(sr *ScriptFilterResponse) {
		sr.Cache = &CacheConfig{
			Seconds:     seconds,
			LooseReload: looseReload,
		}
	}
}

func NewScriptFilterResponse(opts ...ScriptFilterResponseOption) *ScriptFilterResponse {
	sr := &ScriptFilterResponse{
		output: os.Stdout,
		Items:  make([]ListItem, 0),
	}

	for _, opt := range opts {
		opt(sr)
	}

	return sr
}

func (sr *ScriptFilterResponse) AddItem(item ListItem) {
	sr.Items = append(sr.Items, item)
}

// SetError will write the error back to Alfred.
// Callers must return after calling this function!
func (sr *ScriptFilterResponse) SetError(err error) {
	sr.Items = []ListItem{
		{
			Title:    "Encountered Error!",
			Subtitle: err.Error(),
			Valid:    false,
		},
	}
	sr.SendFeedback()
}

func (sr *ScriptFilterResponse) SendFeedback() {
	encoder := json.NewEncoder(sr.output)
	encoder.Encode(sr)
}

type ScriptActionResponse struct {
	Contents struct {
		Arg       string            `json:"arg"`
		Config    map[string]string `json:"config"`
		Variables map[string]string `json:"variables"`
	} `json:"alfredworkflow"`
	encoder *json.Encoder
	output  io.Writer
}

type ScriptActionResponseOption func(*ScriptActionResponse)

func ScriptActionWithOutput(w io.Writer) ScriptActionResponseOption {
	return func(sar *ScriptActionResponse) {
		sar.output = w
	}
}

func NewScriptActionResponse(opts ...ScriptActionResponseOption) *ScriptActionResponse {
	a := &ScriptActionResponse{}
	a.output = os.Stdout
	for _, opt := range opts {
		opt(a)
	}

	a.Contents.Config = make(map[string]string)
	a.Contents.Variables = make(map[string]string)
	a.encoder = json.NewEncoder(a.output)

	return a
}

func (a *ScriptActionResponse) SetConfig(key, value string) {
	a.Contents.Config[key] = value
}

func (a *ScriptActionResponse) SetVariable(key, value string) {
	a.Contents.Variables[key] = value
}

func (a *ScriptActionResponse) SetError(err error) {
	a.Contents.Variables[ExecStatus] = StatusFailed
	a.Contents.Variables[MessageTitle] = "Encountered an error!"
	a.Contents.Variables[MessageBody] = err.Error()
	a.SendFeedback()
}

func RecoverIfErr(a Response) func() {
	return func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				a.SetError(err)
			} else {
				a.SetError(errors.New("undefined error"))
			}
		}
	}
}

func (a *ScriptActionResponse) SendFeedback() {
	encoder := json.NewEncoder(a.output)
	encoder.Encode(a)
}

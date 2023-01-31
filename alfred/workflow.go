package alfred

import (
	"errors"

	"github.com/zalando/go-keyring"
)

type Workflow struct {
	applicationName string
	settingsClient  *Settings
	keyringName     string
	config          map[string]interface{}
}

func NewWorkflow(name string) (*Workflow, error) {
	w := &Workflow{
		applicationName: name,
	}

	settings, err := NewKvSettings(w.applicationName)
	if err != nil {
		return w, err
	}
	w.settingsClient = settings

	data, err := w.settingsClient.Load()
	if err != nil {
		return w, err
	}
	cfg, ok := data.(*map[string]interface{})
	if !ok {
		return w, errors.New("unexpected type assertion failure")
	}
	w.config = *cfg

	w.keyringName = "go-alfred." + w.applicationName

	return w, nil
}

func (w *Workflow) GetConfigString(key, defaultVal string) string {
	v, ok := w.config[key].(string)
	if !ok {
		return defaultVal
	}

	return v
}

func (w *Workflow) GetConfigInt(key string, defaultVal int) int {
	v, ok := w.config[key].(int)
	if !ok {
		return defaultVal
	}

	return v
}

func (w *Workflow) SetConfig(key string, value interface{}) error {
	w.config[key] = value
	return w.settingsClient.Save(w.config)
}

func (w *Workflow) GetSecret(key string) (string, error) {
	return keyring.Get(w.keyringName, key)
}

func (w *Workflow) SaveSecret(key, value string) error {
	return keyring.Set(w.keyringName, key, value)
}

func (w *Workflow) DeleteSecret(key string) error {
	return keyring.Delete(w.keyringName, key)
}

func RunScriptAction(action func(*ScriptActionResponse) error, opts ...ScriptActionResponseOption) {
	sar := NewScriptActionResponse(opts...)
	defer RecoverIfErr(sar)()

	err := action(sar)

	if err != nil {
		sar.SetError(err)
	}

	sar.SendFeedback()
}

func RunScriptFilter(action func(*ScriptFilterResponse) error, opts ...ScriptFilterResponseOption) {
	sfr := NewScriptFilterResponse(opts...)
	defer RecoverIfErr(sfr)()

	err := action(sfr)

	if err != nil {
		sfr.SetError(err)
	}

	sfr.SendFeedback()
}

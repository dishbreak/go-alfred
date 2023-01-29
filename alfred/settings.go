package alfred

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/user"
	"path"
)

type Settings struct {
	payload  interface{}
	filePath string
}

func NewKvSettings(name string) (*Settings, error) {
	p := make(map[string]interface{})
	return NewSettings(name, p)
}

func NewSettings(name string, payload interface{}) (*Settings, error) {
	s := &Settings{
		payload: payload,
	}

	current, err := user.Current()
	if err != nil {
		return s, err
	}

	s.filePath = path.Join(current.HomeDir, SettingsDir, name, SettingsFilename)
	return s, nil
}

func (s *Settings) Load() (interface{}, error) {

	fp, err := os.Open(s.filePath)
	if errors.Is(err, fs.ErrNotExist) {
		err := s.Save(s.payload)
		return s.payload, err
	} else if err != nil {
		return nil, err
	}

	defer fp.Close()

	decoder := json.NewDecoder(fp)
	return s.payload, decoder.Decode(s.payload)
}

func (s *Settings) Save(payload interface{}) error {
	s.payload = payload
	if err := os.MkdirAll(path.Dir(s.filePath), 0700); err != nil {
		return err
	}

	fp, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	encoder := json.NewEncoder(fp)
	return encoder.Encode(s.payload)
}

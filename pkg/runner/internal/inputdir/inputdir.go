package inputdir

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// InputDir represents an input directory for ansible-runner.
type InputDir struct {
	Path         string
	PlaybookPath string
	RolePath     string
	Parameters   map[string]interface{}
	EnvVars      map[string]string
}

// makeDirs creates the required directory structure.
func (i *InputDir) makeDirs() error {
	for _, path := range []string{"env", "project", "inventory"} {
		fullPath := filepath.Join(i.Path, path)
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			logrus.Errorf("unable to create directory %v", fullPath)
			return err
		}
	}
	return nil
}

// addFile adds a file to the given relative path within the input directory.
func (i *InputDir) addFile(path string, content []byte) error {
	fullPath := filepath.Join(i.Path, path)
	err := ioutil.WriteFile(fullPath, content, 0644)
	if err != nil {
		logrus.Errorf("unable to write file %v", fullPath)
	}
	return err
}

// Write commits the object's state to the filesystem at i.Path.
func (i *InputDir) Write() error {
	paramBytes, err := json.Marshal(i.Parameters)
	if err != nil {
		return err
	}
	envVarBytes, err := json.Marshal(i.EnvVars)
	if err != nil {
		return err
	}

	err = i.makeDirs()
	if err != nil {
		return err
	}

	err = i.addFile("env/envvars", envVarBytes)
	if err != nil {
		return err
	}
	err = i.addFile("env/extravars", paramBytes)
	if err != nil {
		return err
	}
	err = i.addFile("inventory/hosts", []byte("localhost ansible_connection=local"))
	if err != nil {
		return err
	}

	if i.PlaybookPath != "" {
		f, err := os.Open(i.PlaybookPath)
		if err != nil {
			logrus.Errorf("failed to open playbook file %v", i.PlaybookPath)
			return err
		}
		defer f.Close()

		playbookBytes, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		err = i.addFile("project/playbook.yaml", playbookBytes)
		if err != nil {
			return err
		}
	}
	return nil
}
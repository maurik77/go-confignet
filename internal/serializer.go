package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// UnmarshalFromFile deserializes the content of the input file in the target interface using the unmarshal function passed in input
func UnmarshalFromFile(path string, target interface{}, unmarshal func(in []byte, out interface{}) (err error)) (err error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("UnmarshalFromFile:File not found %v", path)
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return fmt.Errorf("UnmarshalFromFile:Error when opening file '%v': '%v'", path, err)
	}

	err = unmarshal(content, target)
	if err != nil {
		return fmt.Errorf("UnmarshalFromFile:Error during Unmarshal(): %v", err)
	}

	return nil
}

// MarshalToFile serializes the interface in the target file path using the marshal function passed in input
func MarshalToFile(path string, source interface{}, marshal func(v interface{}) ([]byte, error)) (err error) {

	data, err := marshal(&source)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0644)

	return err
}

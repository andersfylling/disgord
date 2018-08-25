package resource

import (
	"encoding/json"
	"testing"
	"errors"
)

// ValidateJSONMarshalling
func validateJSONMarshalling(b []byte, v interface{}) error {
	var err error

	// convert to struct
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	// back to json
	prettyJSON, err := json.MarshalIndent(&v, "", "    ")
	if err != nil {
		return err
	}

	// sort the data by keys
	omg := make(map[string]interface{})
	err = json.Unmarshal(prettyJSON, &omg)
	if err != nil {
		return err
	}

	omgAgain := make(map[string]interface{})
	err = json.Unmarshal(b, &omgAgain)
	if err != nil {
		return err
	}

	// it is expected that v will contain at least all keys that exists in b
	var missing []string
	for ki, _ := range omgAgain {
		var found bool
		for kj, _ := range omg {
			if kj == ki {
				found = true
				break
			}
		}

		if !found {
			missing = append(missing, ki)
		}
	}

	if len(missing) > 0 {
		missingRows := "JSON data differs. Missing the following keys in struct:\n"
		for _, i := range missing {
			missingRows += i + "\n"
		}
		return errors.New(missingRows)
	}

	// Note: this test doesn't compare the values. I won't bother with handling string pointers and such.. yet.
	// TODO: could still create a test for this. check if it's either nil or an empty value. Although I assume
	// TODO-1: that json.Unmarshal will always fill in the values given the correct data type
	return nil
}

func check(err error, t *testing.T) {
	// Hide function from stacktrace, PR#3
	t.Helper()

	if err != nil {
		t.Error(err)
	}
}

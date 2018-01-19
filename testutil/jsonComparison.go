package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func ValidateJSONMarshalling(b []byte, v interface{}) error {
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

	// minify for comparison
	dst1 := bytes.Buffer{}
	err = json.Compact(&dst1, b)
	if err != nil {
		return err
	}
	dst2 := bytes.Buffer{}
	err = json.Compact(&dst2, prettyJSON)
	if err != nil {
		return err
	}

	// compare
	if dst2.String() != dst1.String() {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(b), string(prettyJSON), false)
		return fmt.Errorf("json data differs. \nDifference \n%s", dmp.DiffPrettyText(diffs))
	}

	return nil
}

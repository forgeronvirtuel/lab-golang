package errorexample

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewNotFoundError(t *testing.T) {
	var err error
	if err = DoStuff("/"); err != nil {
		t.Error(err)
	}

	err = DoStuff("/not-found")
	if err == nil {
		t.Error("No error raised !")
	}

	// Print it in a string format
	fmt.Println(err)

	// Print it in a JSON format
	jsonified, _ := json.Marshal(err)
	fmt.Println(string(jsonified))
}

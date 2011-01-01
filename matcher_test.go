package hamcrest

import (
	"testing"
)

func Test_NewDescription(t *testing.T) {
	description := NewDescription("%v %v", "foo", "bar")
	descriptionString := description.String()
	if descriptionString != "foo bar" {
		t.Errorf("Description should be 'foo bar', was %v", descriptionString)
	}
}


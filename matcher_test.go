// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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


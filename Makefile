# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=hamcrest
GOFILES=\
	core.go\
	matcher.go\
	
include $(GOROOT)/src/Make.pkg

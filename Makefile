# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

PREFIX = github.com/rdrdr/hamcrest

DEPS=\
	$(PREFIX)/base \
	$(PREFIX)/asserter \
	$(PREFIX)/core \
	$(PREFIX)/logic \
	$(PREFIX)/collections \
	$(PREFIX)/reflect \
	$(PREFIX)/strings \


.PHONY: all bench clean install nuke test
all bench clean install nuke test: $(DEPS)

all: TARGET=nuke install test
bench: TARGET=bench
clean: TARGET=clean
nuke: TARGET=nuke
test: TARGET=test
install: TARGET=install

$(DEPS): force
	make -C ../../../$@ $(TARGET)

.PHONY: force
force :;

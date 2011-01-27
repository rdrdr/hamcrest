# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

DEPS=\
	hamcrest \
	hamcrest/core \
	hamcrest/asserter \
	hamcrest/logic \
	hamcrest/comparison \
	hamcrest/collections \
	hamcrest/reflect \
	hamcrest/strings \


.PHONY: bench clean install nuke test
bench clean install nuke test: $(DEPS)

bench: TARGET=bench
clean: TARGET=clean
nuke: TARGET=nuke
test: TARGET=test
install: TARGET=install

$(DEPS): force
	make -C $@ $(TARGET)

.PHONY: force
force :;

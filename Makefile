# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

PREFIX = github.com/rdrdr/hamcrest

all: clean install test

.PHONY: all bench clean install nuke test



DEPS=\
	$(PREFIX)/base \
	$(PREFIX)/asserter \
	$(PREFIX)/core \
	$(PREFIX)/reflect \
	$(PREFIX)/logic \
	$(PREFIX)/slices \
	$(PREFIX)/collections \
	$(PREFIX)/strings \


.PHONY: all bench clean install nuke test

all: clean install test bench

bench: install
	make -C base bench
	make -C asserter bench
	make -C core bench
	make -C reflect bench
	make -C logic bench
	make -C slices bench
	make -C collections bench
	make -C strings bench

clean: 
	make -C base clean
	make -C asserter clean
	make -C core clean
	make -C reflect clean
	make -C logic clean
	make -C slices clean
	make -C collections clean
	make -C strings clean

install:
	make -C base install
	make -C asserter install
	make -C core install
	make -C reflect install
	make -C logic install
	make -C slices install
	make -C collections install
	make -C strings install

nuke: 
	make -C base nuke
	make -C asserter nuke
	make -C core nuke
	make -C reflect nuke
	make -C logic nuke
	make -C slices nuke
	make -C collections nuke
	make -C strings nuke

test: install
	make -C base test
	make -C asserter test
	make -C core test
	make -C reflect test
	make -C logic test
	make -C slices test
	make -C collections test
	make -C strings test

.PHONY: force
force :;

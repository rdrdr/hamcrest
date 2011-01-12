# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

DEPS=\
	hamcrest \
	hamcrest/asserter \
	hamcrest/logic \
	hamcrest/comparison \
	hamcrest/collections \
	hamcrest/reflect \
	hamcrest/strings \

TARG=hamcrest_demo

GOFILES=\
	hamcrest_demo.go \

include $(GOROOT)/src/Make.cmd


all:
	for dep in $(DEPS); do \
		echo "Cleaning $${dep}" ; \
		pushd $${dep} ; gomake clean ; gomake ; gomake install ; gotest ; popd ; \
	done;

clean:
	for dep in $(DEPS); do \
		echo "Cleaning $${dep}" ; \
		pushd $${dep} ; gomake clean ; popd ; \
	done;


install:
	for dep in $(DEPS); do \
		echo "Installing $${dep}" ; \
		pushd $${dep} ; gomake install ; popd ; \
	done


test:
	for dep in $(DEPS); do \
		echo "Testing $${dep}" ; \
		pushd $${dep} ; gotest ; popd ; \
	done


uninstall:
	for dep in $(DEPS); do \
		echo "Uninstalling $${dep}" ; \
		rm -rf $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/"$${dep}".a ; \
	done


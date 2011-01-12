# Copyright 2011 Mick Killianey.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

PACKAGES=\
	hamcrest \
	hamcrest/asserter \
	hamcrest/logic \
	hamcrest/comparison \
	hamcrest/collections \
	hamcrest/reflect \
	hamcrest/strings \


TARG=pkg
GOFILES=\
	doc.go\


all:
	for package in $(PACKAGES); do \
		echo "Cleaning $${package}" ; \
		pushd $${package} ; gomake clean ; gomake ; gomake install ; gotest ; popd ; \
	done;

clean:
	for package in $(PACKAGES); do \
		echo "Cleaning $${package}" ; \
		pushd $${package} ; gomake clean ; popd ; \
	done;


build:
	for package in $(PACKAGES); do \
		echo "Installing $${package}" ; \
		pushd $${package} ; gomake install ; popd ; \
	done


test:
	for package in $(PACKAGES); do \
		echo "Testing $${package}" ; \
		pushd $${package} ; gotest ; popd ; \
	done

include $(GOROOT)/src/Make.pkg

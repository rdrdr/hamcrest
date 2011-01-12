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


install:
	for package in $(PACKAGES); do \
		echo "Installing $${package}" ; \
		pushd $${package} ; gomake install ; popd ; \
	done


test:
	for package in $(PACKAGES); do \
		echo "Testing $${package}" ; \
		pushd $${package} ; gotest ; popd ; \
	done


uninstall:
	for package in $(PACKAGES); do \
		echo "Uninstalling $${package}" ; \
		rm -rf $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/"$${package}".a ; \
	done


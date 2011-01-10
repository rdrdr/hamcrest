#!/usr/bin/env bash

maketest() {
	for i
	do
		(
			pushd $i
			gomake clean
			gomake
			gomake install
			gomake test
			popd
		) || exit $?
	done
}

maketest \
	. \
	asserter \
	logic \
	comparison \
	collections \
	reflect \
	strings \


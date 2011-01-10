#!/usr/bin/env bash

makeclean() {
	for i
	do
		(
			pushd $i
			gomake clean
			popd
		) || exit $?
	done
}

makeclean \
	. \
	asserter \
	collections \
	comparison \
	logic \
	reflect \
	strings \

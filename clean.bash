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
	reflect \
	strings \
	comparison \
#!/bin/sh

# Travis errors when trying to use the "normal" 'go generate' way of making this file.
# I don't know why.

set -x

if [ "$CI" = "true" ] && [ "$TRAVIS" = "true" ]; then
	echo "In Travis CI. Building SO with an adhoc Makefile."
	# Use a temporary adhoc Makefile located in a child-dir of the evmc submodule to build the required example_vm.so file.
	# Once finished, remove the adhoc Makefile.
	# The important part of this command (vs. the 'go generate' command)
	# is the -fPIC flags.
	> ./evmc/bindings/go/evmc/Makefile \
	echo -e 'example_vm.so:\n\tgcc -fPIC -shared ../../../examples/example_vm/example_vm.c -I../../../include -o example_vm.so'
	make -C ./evmc/bindings/go/evmc/ example_vm.so
	rm -f ./evmc/bindings/go/evmc/Makefile
else
    go generate ./evmc/bindings/go/evmc/
fi

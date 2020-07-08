#!/usr/bin/env bash

# Travis errors when trying to use the "normal" 'go generate' way of making this file.
# I don't know why.

if [[ $CI == true ]] && [[ $TRAVIS == true ]]; then
	# Use a temporary adhoc Makefile located in a child-dir of the evmc submodule to build the required example_vm.so file.
	# Once finished, remove the adhoc Makefile.
	> ./evmc/bindings/go/evmc/Makefile \
	echo -e 'example_vm.so:\n\tgcc -fPIC -shared ../../../examples/example_vm/example_vm.c -I../../../include -o example_vm.so'
	make -C ./evmc/bindings/go/evmc/ example_vm.so
	rm -f ./evmc/bindings/go/evmc/Makefile
else
    go generate ./evmc/bindings/go/evmc/
fi

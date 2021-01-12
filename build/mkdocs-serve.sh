#!/bin/sh

# This script installs requirements for and starts the live-reloading development server for MkDocs.
# MkDocs generates a static site from the markdown files in docs/.

py_command=python
pip_command="${py_command} -m pip"

setup_py_command() {
	if [ -z $(which ${py_command}) ]; then
		py_command=python3
		return
	fi
	py_command_major_version=$(${py_command} 2>&1 --version | cut -d' ' -f2 | head -c1)
	if [ $py_command_major_version != 3 ]; then
		py_command=python3
	fi
}
setup_py_command

set -ex
${py_command} --version
${pip_command} --version
set +x

${pip_command} install -r requirements-mkdocs.txt

mkdocs serve

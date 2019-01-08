#!/usr/bin/env bash

GETH_ARCHIVE_NAME="multi-geth-$TRAVIS_OS_NAME"
zip -j "$GETH_ARCHIVE_NAME.zip" build/bin/geth

shasum -a 256 $GETH_ARCHIVE_NAME.zip
shasum -a 256 $GETH_ARCHIVE_NAME.zip > $GETH_ARCHIVE_NAME.zip.sha256

ALLTOOLS_ARCHIVE_NAME="multi-geth-alltools-$TRAVIS_OS_NAME"
zip -j "$ALLTOOLS_ARCHIVE_NAME.zip" build/bin/*

shasum -a 256 $ALLTOOLS_ARCHIVE_NAME.zip
shasum -a 256 $ALLTOOLS_ARCHIVE_NAME.zip > $ALLTOOLS_ARCHIVE_NAME.zip.sha256


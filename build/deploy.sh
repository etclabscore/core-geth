#!/usr/bin/env bash


GETH_ARCHIVE_NAME="core-geth-${BUILD_OS_NAME}-$(git describe --abbrev=0 --tags)"
ALLTOOLS_ARCHIVE_NAME="core-geth-alltools-${BUILD_OS_NAME}-$(git describe --abbrev=0 --tags)"

if [[ "${BUILD_OS_NAME}" == "win64" ]]; then
  7z a "$GETH_ARCHIVE_NAME.zip" ./build/bin/geth.exe

  sha256sum $GETH_ARCHIVE_NAME.zip
  sha256sum $GETH_ARCHIVE_NAME.zip > $GETH_ARCHIVE_NAME.zip.sha256

  7z a "$ALLTOOLS_ARCHIVE_NAME.zip" ./build/bin/*

  sha256sum $ALLTOOLS_ARCHIVE_NAME.zip
  sha256sum $ALLTOOLS_ARCHIVE_NAME.zip > $ALLTOOLS_ARCHIVE_NAME.zip.sha256

else
  zip -j "$GETH_ARCHIVE_NAME.zip" build/bin/geth

  shasum -a 256 $GETH_ARCHIVE_NAME.zip
  shasum -a 256 $GETH_ARCHIVE_NAME.zip > $GETH_ARCHIVE_NAME.zip.sha256

  zip -j "$ALLTOOLS_ARCHIVE_NAME.zip" build/bin/*

  shasum -a 256 $ALLTOOLS_ARCHIVE_NAME.zip
  shasum -a 256 $ALLTOOLS_ARCHIVE_NAME.zip > $ALLTOOLS_ARCHIVE_NAME.zip.sha256
fi

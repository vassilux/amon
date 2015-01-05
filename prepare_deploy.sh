#!/bin/bash
#
# 
# Description : Prepare deploy amon. 
# Author : vassilux
# Last modified : 2015-01-05 14:53:54  
#

set -e

VER_MAJOR="1"
VER_MINOR="0"
VER_PATCH="0"

make

DEPLOY_DIR="amon_${VER_MAJOR}.${VER_MINOR}.${VER_PATCH}"
DEPLOY_FILE_NAME="amon_${VER_MAJOR}.${VER_MINOR}.${VER_PATCH}.tar.gz"

if [ -d "$DEPLOY_DIR" ]; then
    rm -rf  "$DEPLOY_DIR"
fi
#
#
mkdir "$DEPLOY_DIR"
mkdir "$DEPLOY_DIR/samples"
mkdir "$DEPLOY_DIR/logs"
mkdir "$DEPLOY_DIR/tmp"

cp -aR ./bin/* "$DEPLOY_DIR"
cp -aR ./samples/* "${DEPLOY_DIR}/samples"
#
mkdir "$DEPLOY_DIR/docs"
pandoc -o "$DEPLOY_DIR/docs/INSTALL.html" ./docs/INSTALL.md
pandoc -o "$DEPLOY_DIR/docs/ReleaseNotes.html" ./docs/ReleaseNotes.md
cp "$DEPLOY_DIR/docs/INSTALL.html" .
cp "$DEPLOY_DIR/docs/ReleaseNotes.html" .

tar cvzf "${DEPLOY_FILE_NAME}" "${DEPLOY_DIR}"

if [ ! -f "$DEPLOY_FILE_NAME" ]; then
    echo "Deploy build failed."
    exit 1
fi

if [ ! -d releases ]; then
        mkdir releases
fi

mv ${DEPLOY_FILE_NAME} ./releases
mv INSTALL.* ./releases
mv ReleaseNotes.* ./releases

rm -rf "$DEPLOY_DIR"



rm -rf "$DEPLOY_DIR"

echo "Deploy build complete."
echo "Live well"

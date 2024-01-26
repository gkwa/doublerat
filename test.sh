#!/usr/bin/env bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $SCRIPT_DIR

rm -rf $SCRIPT_DIR/mytest1706299472
rm -rf /tmp/doublerat/t1
rm -rf /tmp/doublerat/t2

# test1
cd $SCRIPT_DIR

./doublerat test1
cd $SCRIPT_DIR/mytest1706299472
git add .gitmodules
git add darksheep
git commit -am sub

rm -rf /tmp/doublerat/t1
mkdir -p /tmp/doublerat/t1
cp -r $SCRIPT_DIR/mytest1706299472 /tmp/doublerat/t1

# test2
cd $SCRIPT_DIR

./doublerat test1
cd $SCRIPT_DIR/mytest1706299472
git submodule add https://github.com/taylormonacelli/darksheep
git status

rm -rf /tmp/doublerat/t2
mkdir -p /tmp/doublerat/t2

cd $SCRIPT_DIR/mytest1706299472
git commit -am sub

cd $SCRIPT_DIR
cp -r $SCRIPT_DIR/mytest1706299472 /tmp/doublerat/t2

diff -uw --brief --recursive /tmp/doublerat/t1 /tmp/doublerat/t2

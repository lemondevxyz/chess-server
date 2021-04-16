#!/bin/sh
export PROD_ID="107.191.62.233"

echo "Building"
./build.sh

echo "Deploying"
rsync -a bin "chess@$PROD_ID:~/" --progress

echo "Cleaning"
rm ./bin

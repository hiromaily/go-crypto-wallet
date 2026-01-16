#!/bin/bash

gofiles=$(find . -name "*.go" | grep -v "/vendor/")

for gofile in $gofiles; do
	echo $gofile
	sed '/^import/,/^[[:space:]]*)/ { /^[[:space:]]*$/ d; }' $gofile >tmp
	mv tmp $gofile
done

go list ./... | grep -v "/vendor/" | xargs -r go fmt
goimports -local github.com/hiromaily/ -l ./ | grep -v "/vendor/" | xargs -r goimports -local github.com/hiromaily/ -w

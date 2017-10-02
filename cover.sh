#!/bin/bash

set -e

COVER=.cover
ROOT_PKG=github.com/getfider/fider/

if [[ -d "$COVER" ]]; then
	rm -rf "$COVER"
fi
mkdir -p "$COVER"

i=0
for pkg in "$@"; do
	i=$((i + 1))

	extracoverpkg=""
	if [[ -f "$GOPATH/src/$pkg/.extra-coverpkg" ]]; then
		extracoverpkg=$( \
			sed -e "s|^|$pkg/|g" < "$GOPATH/src/$pkg/.extra-coverpkg" \
			| tr '\n' ',')
	fi

	coverpkg=$(go list -json "$pkg" | jq -r '
		.Deps
		| . + ["'"$pkg"'"]
		| map
			( select(startswith("'"$ROOT_PKG"'"))
			| select(contains("/vendor/") | not)
			)
		| join(",")
	')
	testcoverpkg=$(go list -json "$pkg" | jq -r '
		.XTestImports
		| . + ["'"$pkg"'"]
		| map
			( select(startswith("'"$ROOT_PKG"'"))
			| select(contains("/vendor/") | not)
			)
		| join(",")
	')
	if [[ -n "$extracoverpkg" ]]; then
		coverpkg="$extracoverpkg$coverpkg"
	fi

	args=""
	if [[ -n "$coverpkg" ]]; then
		args="-p=1 -coverprofile $COVER/cover.${i}.out -covermode=atomic -coverpkg $coverpkg,$testcoverpkg"
	fi

	echo go test $args -v -race "$pkg"
	go test -i $args -v -race "$pkg"
done

gocovmerge "$COVER"/*.out > cover.out
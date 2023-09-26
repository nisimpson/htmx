#!/bin/sh
GOPROXY=proxy.golang.org
TAG=$1

git tag $TAG
git push origin --tags
go list -m "github.com/nisimpson/htmx@${TAG}"
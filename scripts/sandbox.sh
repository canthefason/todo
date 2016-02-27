#!/usr/bin/env bash
set -e

if [ "$TRAVIS_BRANCH" != "master" ]; then
  exit 0
fi

git config --global branch.autosetupmerge always
git checkout -b sandbox
git merge master
git push origin sandbox

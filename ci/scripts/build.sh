#!/bin/bash -eux

pushd dp-sitemap
  make build
  cp build/dp-sitemap Dockerfile.concourse ../build
popd

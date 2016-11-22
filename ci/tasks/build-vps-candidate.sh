#!/usr/bin/env bash
set -e

semver=`cat version-semver/number`

pushd vps-release

  mkdir -p blobs/golang_1.7/
  pushd blobs/golang_1.7/
    wget http://golangtc.com/static/go/1.7/go1.7.linux-amd64.tar.gz
  popd

# for debug, remove other jobs and packages
  rm -rf jobs/redis
  rm -rf jobs/baremetal-provision-server
  rm -rf jobs/xcat-server
  rm -rf packages/baremetal-provision-server
  rm -rf packages/libyaml
  rm -rf packages/redis
  rm -rf packages/ruby-2.2
  rm -rf packages/xcat
  
  echo "installing bosh CLI"
  gem install bosh_cli --no-ri --no-rdo c

  echo "using bosh CLI version..."
  bosh version

  vps_release_name="vps"


  echo "building baremetal server dev release..."
  bosh create release --name $vps_release_name --version $semver --with-tarball --force
popd

mv vps-release/dev_releases/$vps_release_name/$vps_release_name-$semver.tgz candidate/
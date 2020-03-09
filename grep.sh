#!/usr/bin/env bash
set -ex

grep -Rn cmd models modules routers services templates -e "$1"

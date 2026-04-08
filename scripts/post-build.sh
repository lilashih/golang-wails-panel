#!/bin/sh

set -eu

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <source-binary>" >&2
  exit 1
fi

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(dirname "$SCRIPT_DIR")
SOURCE_BINARY=$1

case "$SOURCE_BINARY" in
  /*) SOURCE_PATH=$SOURCE_BINARY ;;
  *) SOURCE_PATH=$REPO_ROOT/$SOURCE_BINARY ;;
esac

if [ ! -f "$SOURCE_PATH" ]; then
  echo "Build output not found: $SOURCE_PATH" >&2
  exit 1
fi

TARGET_DIR=$REPO_ROOT/release
TARGET_PATH=$TARGET_DIR/$(basename "$SOURCE_PATH")

mkdir -p "$TARGET_DIR"
cp "$SOURCE_PATH" "$TARGET_PATH"
echo "Copied build output to $TARGET_PATH"

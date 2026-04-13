#!/bin/sh

set -eu

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <source-binary> [output-dir]" >&2
  exit 1
fi

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(dirname "$SCRIPT_DIR")
SOURCE_BINARY=$1
OUTPUT_DIR=${2:-"$REPO_ROOT/release"}

case "$SOURCE_BINARY" in
  /*) SOURCE_PATH=$SOURCE_BINARY ;;
  *) SOURCE_PATH=$REPO_ROOT/$SOURCE_BINARY ;;
esac

case "$OUTPUT_DIR" in
  /*) OUTPUT_PATH=$OUTPUT_DIR ;;
  *) OUTPUT_PATH=$REPO_ROOT/$OUTPUT_DIR ;;
esac

if [ ! -f "$SOURCE_PATH" ]; then
  echo "Build output not found: $SOURCE_PATH" >&2
  exit 1
fi

if ! command -v linuxdeploy >/dev/null 2>&1; then
  echo "linuxdeploy not found in PATH" >&2
  exit 1
fi

if ! command -v linuxdeploy-plugin-gtk.sh >/dev/null 2>&1; then
  echo "linuxdeploy-plugin-gtk.sh not found in PATH" >&2
  exit 1
fi

APP_NAME=$(basename "$SOURCE_PATH")
APPDIR_ROOT=$(mktemp -d)
APPDIR=$APPDIR_ROOT/$APP_NAME.AppDir
ICON_FILE=$APPDIR/usr/share/icons/hicolor/512x512/apps/$APP_NAME.png
DEPLOY_TARGETS_FILE=$APPDIR_ROOT/deploy-targets.txt

cleanup() {
  rm -rf "$APPDIR_ROOT"
}

trap cleanup EXIT

copy_with_parents() {
  SOURCE_FILE=$1
  TARGET_DIR=$APPDIR$(dirname "$SOURCE_FILE")

  mkdir -p "$TARGET_DIR"
  cp -a "$SOURCE_FILE" "$TARGET_DIR/"
}

register_deploy_target() {
  DEPLOY_TARGET=$1
  printf '%s\n' "$DEPLOY_TARGET" >> "$DEPLOY_TARGETS_FILE"
}

bundle_if_exists() {
  BUNDLE_PATH=$1

  if [ -f "$BUNDLE_PATH" ]; then
    copy_with_parents "$BUNDLE_PATH"
    register_deploy_target "$APPDIR$BUNDLE_PATH"
  fi
}

bundle_first_match() {
  MATCH_NAME=$1
  FOUND_PATH=$(find /usr/lib /usr/lib64 /usr/libexec -type f -name "$MATCH_NAME" 2>/dev/null | head -n 1 || true)

  if [ -n "$FOUND_PATH" ]; then
    copy_with_parents "$FOUND_PATH"
    register_deploy_target "$APPDIR$FOUND_PATH"
  fi
}

mkdir -p \
  "$APPDIR/usr/bin" \
  "$APPDIR/usr/share/applications" \
  "$APPDIR/usr/share/icons/hicolor/512x512/apps"

cp "$SOURCE_PATH" "$APPDIR/usr/bin/$APP_NAME"
chmod +x "$APPDIR/usr/bin/$APP_NAME"

if [ -f "$REPO_ROOT/appicon.png" ]; then
  cp "$REPO_ROOT/appicon.png" "$ICON_FILE"
fi

cat > "$APPDIR/usr/share/applications/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=golang-wails-panel
Exec=$APP_NAME
Icon=$APP_NAME
Categories=Utility;
Terminal=false
StartupNotify=true
EOF

cat > "$APPDIR/AppRun" <<EOF
#!/bin/sh
set -eu

APPDIR=\${APPDIR:-\$(CDPATH= cd -- "\$(dirname -- "\$0")" && pwd)}

export PATH="\$APPDIR/usr/bin:\$PATH"
export LD_LIBRARY_PATH="\$APPDIR/usr/lib:\$APPDIR/usr/lib/x86_64-linux-gnu:\$APPDIR/usr/lib64\${LD_LIBRARY_PATH:+:\$LD_LIBRARY_PATH}"

if [ -d "\$APPDIR/apprun-hooks" ]; then
  for hook in "\$APPDIR"/apprun-hooks/*.sh; do
    if [ -f "\$hook" ]; then
      # shellcheck disable=SC1090
      . "\$hook"
    fi
  done
fi

for gio_dir in "\$APPDIR/usr/lib/gio/modules" "\$APPDIR/usr/lib/x86_64-linux-gnu/gio/modules" "\$APPDIR/usr/lib64/gio/modules"; do
  if [ -d "\$gio_dir" ]; then
    export GIO_MODULE_DIR="\$gio_dir"
    break
  fi
done

cd "\$APPDIR/usr"
exec "./bin/$APP_NAME" "\$@"
EOF

chmod +x "$APPDIR/AppRun"
ln -sf "usr/share/applications/$APP_NAME.desktop" "$APPDIR/$APP_NAME.desktop"

if [ -f "$ICON_FILE" ]; then
  ln -sf "usr/share/icons/hicolor/512x512/apps/$APP_NAME.png" "$APPDIR/$APP_NAME.png"
fi

bundle_first_match WebKitNetworkProcess
bundle_first_match WebKitWebProcess
bundle_first_match WebKitGPUProcess
bundle_if_exists /usr/lib/x86_64-linux-gnu/webkit2gtk-4.0/injected-bundle/libwebkit2gtkinjectedbundle.so
bundle_if_exists /usr/lib/aarch64-linux-gnu/webkit2gtk-4.0/injected-bundle/libwebkit2gtkinjectedbundle.so

mkdir -p "$OUTPUT_PATH"

(
  cd "$OUTPUT_PATH"
  export APPIMAGE_EXTRACT_AND_RUN=1
  export OUTPUT=${APP_NAME}.AppImage
  set -- \
    --appdir "$APPDIR" \
    --desktop-file "$APPDIR/usr/share/applications/$APP_NAME.desktop" \
    --executable "$APPDIR/usr/bin/$APP_NAME"

  if [ -f "$ICON_FILE" ]; then
    set -- "$@" --icon-file "$ICON_FILE"
  fi

  if [ -f "$DEPLOY_TARGETS_FILE" ]; then
    while IFS= read -r deploy_target; do
      if [ -n "$deploy_target" ]; then
        set -- "$@" --deploy-deps-only "$deploy_target"
      fi
    done < "$DEPLOY_TARGETS_FILE"
  fi

  set -- "$@" --plugin gtk --output appimage
  linuxdeploy "$@"
)

echo "AppImage build completed: $OUTPUT_PATH/${APP_NAME}.AppImage"

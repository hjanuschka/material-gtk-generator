#!/bin/bash

PREF_FILE="$HOME/.config/chromium/Default/Preferences"
PREF_DIR="$(dirname "$PREF_FILE")"
if [ ! -f "$PREF_FILE" ]; then
      mkdir -p "$(dirname "$PREF_FILE")"
      echo '{"extensions":{"theme":{"system_theme":1}}}' > "$PREF_FILE"
fi

jq '.extensions.theme.system_theme = 1' ~/.config/chromium/Default/Preferences > tmp && mv tmp ~/.config/chromium/Default/Preferences

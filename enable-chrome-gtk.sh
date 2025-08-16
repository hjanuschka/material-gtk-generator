#!/bin/bash

# Force Chrome to use GTK theme by modifying preferences
# This script enables "Use GTK+ theme" programmatically

CHROME_CONFIG="$HOME/.config/chromium/Default/Preferences"
CHROME_LOCAL_STATE="$HOME/.config/chromium/Local State"

enable_gtk_theme() {
    if [[ -f "$CHROME_CONFIG" ]]; then
        echo "🔧 Enabling GTK theme in Chrome preferences..."
        
        # Backup original
        cp "$CHROME_CONFIG" "$CHROME_CONFIG.backup"
        
        # Enable GTK theme (system_theme: true)
        if grep -q '"system_theme"' "$CHROME_CONFIG"; then
            # Update existing setting
            sed -i 's/"system_theme":[^,}]*/"system_theme":1/g' "$CHROME_CONFIG"
            echo "✅ Updated existing system_theme setting"
        else
            # Add new setting to browser object
            sed -i 's/"browser":{/"browser":{"system_theme":1,/g' "$CHROME_CONFIG"
            echo "✅ Added system_theme setting"
        fi
        
        echo "🎨 GTK theme enabled in Chrome preferences"
    else
        echo "❌ Chrome config not found. Launch Chrome once first."
        return 1
    fi
}

# Main execution
echo "🚀 Enabling Chrome GTK theme..."

# Check if Chrome is running
if pgrep -x "chromium" > /dev/null; then
    echo "⚠️  Chrome is running. Close it first for changes to take effect."
    read -p "Close Chrome and press Enter to continue..."
fi

enable_gtk_theme

echo
echo "💡 Now launch Chrome normally - it should automatically use GTK theme!"
echo "   chromium"
echo
echo "🔍 To verify: chrome://settings/appearance should show 'Use GTK+ theme' enabled"
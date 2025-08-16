# 🎨 Material 3 GTK Theme Generator

A Go implementation of Chrome's exact Material Design 3 color system for generating GTK themes that match Chrome's internal theming.

## ✨ Features

- **Chrome's Exact Algorithm**: Ports Chrome's C++ Material Color Utilities directly from `ui/color/dynamic_color/palette_factory.cc`
- **Material 3 Variants**: Supports TonalSpot, Vibrant, Expressive, and Neutral color schemes
- **Proper Tone Mappings**: Uses Chrome's actual neutral98 base colors with primary accents
- **Hot-Reload**: Automatically switches themes to trigger Chrome reload without restart
- **Perfect Color Science**: HCT color space with proper chroma and hue rotations

## 🚀 Installation

```bash
git clone https://github.com/hjanuschka/material-gtk-generator.git
cd material-gtk-generator
go build -o material-gtk
```

## 📖 Usage

```bash
# Basic usage - generate and apply theme
./material-gtk -apply 255,255,0

# Different variants
./material-gtk -apply -variant vibrant 255,0,0
./material-gtk -apply -variant expressive 0,150,255

# Output to file
./material-gtk 28,32,39 > my-theme.css
```

## 🎯 Background

This tool replaces the functionality of Chrome CL 6832165 (`--set-theme-color` flag) which was rejected by the Chrome team. Instead of a CLI flag, this generates GTK themes that Chrome can read when "Use GTK+ theme" is enabled.

## 🧪 Chrome Setup

Enable GTK theming in Chrome:

```bash
# Method 1: Force GTK theme on launch
chromium --force-system-theme

# Method 2: Use the included script
./enable-chrome-gtk.sh
```

## 🎨 Material 3 Variants

- **TonalSpot** (default): Moderate saturation (chroma 40)
- **Vibrant**: High saturation (chroma 200) with hue rotations
- **Expressive**: Creative color combinations with varied rotations
- **Neutral**: Muted, sophisticated palette

## 🔬 Technical Details

This implementation:
1. Converts RGB input to HCT color space
2. Generates Material 3 palettes using Chrome's exact chroma values
3. Maps to Chrome's neutral98 base + primary accent architecture
4. Creates GTK CSS that Chrome's theme system can parse

## 🤝 Contributing

Issues and PRs welcome! This tool aims to provide pixel-perfect Chrome color matching.

## 📄 License

MIT License - see LICENSE file for details.

---

*🤖 Generated with [Claude Code](https://claude.ai/code)*
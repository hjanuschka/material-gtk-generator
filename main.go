package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func parseRGB(rgbStr string) (uint8, uint8, uint8, error) {
	parts := strings.Split(rgbStr, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid RGB format. Use R,G,B (e.g., 28,32,39)")
	}

	r, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || r < 0 || r > 255 {
		return 0, 0, 0, fmt.Errorf("invalid red value: %s", parts[0])
	}

	g, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || g < 0 || g > 255 {
		return 0, 0, 0, fmt.Errorf("invalid green value: %s", parts[1])
	}

	b, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil || b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("invalid blue value: %s", parts[2])
	}

	return uint8(r), uint8(g), uint8(b), nil
}

func rgbaToARGB(c color.RGBA) int {
	return int(c.A)<<24 | int(c.R)<<16 | int(c.G)<<8 | int(c.B)
}

func argbToHex(argb int) string {
	r := (argb >> 16) & 0xFF
	g := (argb >> 8) & 0xFF
	b := argb & 0xFF
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func colorToHex(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func generateGTKTheme(seedColor color.RGBA, variant string) string {
	// Special case: handle black like Chrome does
	if seedColor.R == 0 && seedColor.G == 0 && seedColor.B == 0 {
		// Chrome converts black to near-black to avoid pink tones
		seedColor = color.RGBA{1, 1, 1, 255}
	}
	
	// Use Chrome's exact Material 3 algorithm
	var chromeVariant SchemeVariant
	switch variant {
	case "vibrant":
		chromeVariant = Vibrant
	case "expressive":
		chromeVariant = Expressive
	case "neutral":
		chromeVariant = Neutral
	case "monochrome":
		chromeVariant = Neutral // Use neutral for monochrome
	default: // tonal_spot
		chromeVariant = TonalSpot
	}
	
	// Generate Chrome's Material 3 palette
	chromePalette := GenerateChromePalette(seedColor, chromeVariant)
	
	// Chrome's actual browser UI tone mappings (light mode)
	// Chrome uses kColorSysBase (neutral98) for toolbar, NOT primary colors!
	
	// Primary colors (for accents, highlights, focus)
	primary := colorToHex(chromePalette.Primary.Tone(40))           // Primary accent
	onPrimary := colorToHex(chromePalette.Primary.Tone(100))        // On Primary (white)
	primaryContainer := colorToHex(chromePalette.Primary.Tone(90))  // Primary Container
	onPrimaryContainer := colorToHex(chromePalette.Primary.Tone(10)) // On Primary Container
	
	// Chrome's browser chrome colors - use neutral base!
	chromeBase := colorToHex(chromePalette.Neutral.Tone(98))         // kColorSysBase - very light!
	chromeOnBase := colorToHex(chromePalette.Neutral.Tone(10))       // Text on base
	
	// Primary tones for highlights and accents
	primary80 := colorToHex(chromePalette.Primary.Tone(80))         // Lighter accent
	primary90 := colorToHex(chromePalette.Primary.Tone(90))         // Very light accent
	
	// Neutral colors (Chrome's actual surface colors)
	surface := colorToHex(chromePalette.Neutral.Tone(99))              // Surface (almost white)
	onSurface := colorToHex(chromePalette.Neutral.Tone(10))            // On Surface (dark text)
	
	// Neutral Variant colors
	surfaceVariant := colorToHex(chromePalette.NeutralVariant.Tone(90))        // Surface Variant
	onSurfaceVariant := colorToHex(chromePalette.NeutralVariant.Tone(30))      // On Surface Variant
	outlineVariant := colorToHex(chromePalette.NeutralVariant.Tone(80))        // Outline Variant

	// Generate GTK CSS with Material 3 colors
	css := fmt.Sprintf(`/*
 * Material 3 GTK Theme - Auto-generated using Material Color Utilities
 * Seed: RGB(%d,%d,%d)
 * Variant: %s
 * Generated: %s
 * 
 * This theme uses Google's Material Design 3 color system
 * with proper HCT color space calculations for harmonious colors
 */

/* Base window styling */
window {
    background-color: %s;      /* Material 3 surface */
    color: %s;                  /* Material 3 on-surface */
    background-image: none;
}

/* Header bar - Chrome uses neutral base (tone 98) for toolbar */
headerbar {
    background-color: %s;       /* Chrome kColorSysBase (neutral98) */
    color: %s;                  /* Chrome kColorSysOnBase (neutral10) */
    background-image: none;
    border-color: %s;           /* Primary accent for borders */
}

/* Button styling - Chrome uses light base with primary accents */
button {
    background-color: %s;       /* Light base background */
    color: %s;                  /* Dark text */
    background-image: none;
    border-color: %s;           /* Primary accent border */
    border-radius: 4px;
}

button:hover {
    background-color: %s;       /* Primary container on hover */
    color: %s;                  /* On primary container */
    background-image: none;
}

button:active {
    background-color: %s;       /* Primary accent when pressed */
    color: %s;                  /* White on primary */
    background-image: none;
}

/* Entry fields - Chrome uses for address bar and input styling */
entry {
    background-color: %s;       /* Material 3 surface variant */
    color: %s;                  /* Material 3 on-surface */
    border-color: %s;           /* Material 3 primary */
    background-image: none;
}

entry:focus {
    border-color: %s;           /* Material 3 primary (tone 80) */
    box-shadow: 0 0 0 1px %s;
}

/* Chrome-targeted selectors */
.titlebar {
    background-color: %s;       /* Chrome neutral base */
    color: %s;                  /* Chrome on-base */
    background-image: none;
}

/* Menu and toolbar styling - match Chrome's light base */
menubar {
    background-color: %s;       /* Chrome neutral base */
    color: %s;                  /* Chrome on-base */
    background-image: none;
}

toolbar {
    background-color: %s;       /* Chrome neutral base */
    color: %s;                  /* Chrome on-base */
    background-image: none;
}

/* Selection colors */
selection {
    background-color: %s;       /* Material 3 primary container */
    color: %s;                  /* Material 3 on-primary container */
}

/* Scrollbar styling */
scrollbar {
    background-color: %s;       /* Material 3 surface */
}

scrollbar slider {
    background-color: %s;       /* Material 3 outline variant */
    border-radius: 8px;
}

scrollbar slider:hover {
    background-color: %s;       /* Material 3 primary */
}

/* Tab styling - important for Chrome tabs */
notebook {
    background-color: %s;       /* Material 3 surface */
}

notebook header {
    background-color: %s;       /* Material 3 primary */
    background-image: none;
}

notebook tab {
    background-color: %s;       /* Material 3 surface variant */
    color: %s;                  /* Material 3 on-surface variant */
    background-image: none;
}

notebook tab:checked {
    background-color: %s;       /* Material 3 primary container */
    color: %s;                  /* Material 3 on-primary container */
    background-image: none;
}

/* Chrome-specific Material 3 additions */
headerbar.titlebar {
    background-color: %s;       /* Chrome neutral base for frame */
}

/* Use Material 3 tones for inactive elements */
.tab:not(:checked) {
    background-color: %s;       /* Material 3 primary tone 90 */
    color: %s;                  /* Material 3 on-surface */
}

/* Ensure all backgrounds are solid colors */
* {
    background-image: none;
}
`,
		seedColor.R, seedColor.G, seedColor.B,
		variant,
		time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
		// Base window
		surface, onSurface,
		// Header bar - use Chrome's neutral base!
		chromeBase, chromeOnBase, primary,
		// Button - use Chrome's neutral base with primary accents
		chromeBase, chromeOnBase, primary,
		// Button hover
		primaryContainer, onPrimaryContainer,
		// Button active - use primary accent
		primary, onPrimary,
		// Entry
		surfaceVariant, onSurface, primary,
		// Entry focus
		primary80, primary,
		// Chrome selectors - use neutral base
		chromeBase, chromeOnBase,
		// Menubar - use neutral base
		chromeBase, chromeOnBase,
		// Toolbar - use neutral base
		chromeBase, chromeOnBase,
		// Selection
		primaryContainer, onPrimaryContainer,
		// Scrollbar
		surface,
		// Scrollbar slider
		outlineVariant,
		// Scrollbar hover
		primary,
		// Notebook
		surface,
		// Notebook header
		primary,
		// Notebook tab
		surfaceVariant, onSurfaceVariant,
		// Notebook tab checked
		primaryContainer, onPrimaryContainer,
		// Chrome-specific - use neutral base
		chromeBase,
		// Inactive tab
		primary90, onSurface,
	)

	return css
}

func main() {
	var (
		rgbInput string
		variant  string
		output   string
		apply    bool
	)

	flag.StringVar(&rgbInput, "rgb", "", "RGB values as R,G,B (e.g., 28,32,39)")
	flag.StringVar(&variant, "variant", "tonal_spot", "Material 3 variant: tonal_spot, vibrant, expressive, neutral, monochrome")
	flag.StringVar(&output, "output", "", "Output file path (default: stdout)")
	flag.BoolVar(&apply, "apply", false, "Automatically apply theme to Chrome via gsettings")
	flag.Parse()

	if rgbInput == "" && len(flag.Args()) > 0 {
		rgbInput = flag.Args()[0]
	}

	if rgbInput == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] R,G,B\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample: %s 28,32,39\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "         %s -variant vibrant -apply 255,0,0\n", os.Args[0])
		os.Exit(1)
	}

	// Parse RGB input
	r, g, b, err := parseRGB(rgbInput)
	if err != nil {
		log.Fatalf("Error parsing RGB: %v", err)
	}

	seedColor := color.RGBA{r, g, b, 255}

	// Generate GTK theme with Material 3 colors
	css := generateGTKTheme(seedColor, variant)

	// Output the CSS
	if output != "" {
		// Create directory if needed
		dir := filepath.Dir(output)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}

		// Write to file
		if err := os.WriteFile(output, []byte(css), 0644); err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("✅ Theme written to %s\n", output)
	} else if !apply {
		// Print to stdout if not applying
		fmt.Print(css)
	}

	// Apply theme if requested
	if apply {
		// First, write to temporary theme
		tempThemePath := filepath.Join(os.Getenv("HOME"), ".themes/OmarchyThemeTemp/gtk-3.0/gtk.css")
		
		// Create temp theme directory
		tempThemeDir := filepath.Dir(tempThemePath)
		if err := os.MkdirAll(tempThemeDir, 0755); err != nil {
			log.Fatalf("Failed to create temp theme directory: %v", err)
		}

		// Write temp theme file
		if err := os.WriteFile(tempThemePath, []byte(css), 0644); err != nil {
			log.Fatalf("Failed to write temp theme file: %v", err)
		}

		// Create temp index.theme
		tempIndexPath := filepath.Join(os.Getenv("HOME"), ".themes/OmarchyThemeTemp/index.theme")
		tempIndexContent := fmt.Sprintf(`[Desktop Entry]
Type=X-GNOME-Metatheme
Name=OmarchyThemeTemp
Comment=Material 3 Theme Temp - RGB(%d,%d,%d)
Encoding=UTF-8

[X-GNOME-Metatheme]
GtkTheme=OmarchyThemeTemp
IconTheme=Adwaita
CursorTheme=Adwaita
`, r, g, b)
		
		if err := os.WriteFile(tempIndexPath, []byte(tempIndexContent), 0644); err != nil {
			log.Fatalf("Failed to write temp index.theme: %v", err)
		}

		// Now write to main theme
		mainThemePath := filepath.Join(os.Getenv("HOME"), ".themes/OmarchyTheme/gtk-3.0/gtk.css")
		
		// Create main theme directory
		mainThemeDir := filepath.Dir(mainThemePath)
		if err := os.MkdirAll(mainThemeDir, 0755); err != nil {
			log.Fatalf("Failed to create main theme directory: %v", err)
		}

		// Write main theme file
		if err := os.WriteFile(mainThemePath, []byte(css), 0644); err != nil {
			log.Fatalf("Failed to write main theme file: %v", err)
		}

		// Create main index.theme
		mainIndexPath := filepath.Join(os.Getenv("HOME"), ".themes/OmarchyTheme/index.theme")
		mainIndexContent := fmt.Sprintf(`[Desktop Entry]
Type=X-GNOME-Metatheme
Name=OmarchyTheme
Comment=Material 3 Theme - RGB(%d,%d,%d)
Encoding=UTF-8

[X-GNOME-Metatheme]
GtkTheme=OmarchyTheme
IconTheme=Adwaita
CursorTheme=Adwaita
`, r, g, b)
		
		if err := os.WriteFile(mainIndexPath, []byte(mainIndexContent), 0644); err != nil {
			log.Fatalf("Failed to write main index.theme: %v", err)
		}

		fmt.Printf("🎨 Material 3 theme created with RGB(%d,%d,%d)\n", r, g, b)
		fmt.Printf("   Variant: %s\n", variant)
		fmt.Printf("   Seed color: %s\n", argbToHex(rgbaToARGB(seedColor)))
		fmt.Printf("✅ Themes saved to ~/.themes/OmarchyTheme and ~/.themes/OmarchyThemeTemp\n")
		
		// Trigger Chrome to reload by switching between our own themes (no flicker)
		fmt.Println("🔄 Triggering theme reload...")
		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "OmarchyThemeTemp")
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to switch to OmarchyThemeTemp: %v", err)
		}
		
		// Wait 1 second to ensure the switch is registered
		time.Sleep(1 * time.Second)
		
		cmd = exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "OmarchyTheme")
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Failed to switch back to OmarchyTheme: %v", err)
		}
		
		fmt.Println("🎉 Chrome should now display with your Material 3 colors!")
		fmt.Println("\nTo use this theme permanently, make sure 'Use GTK+ theme' is enabled")
		fmt.Println("in chrome://settings/appearance")
	}
}
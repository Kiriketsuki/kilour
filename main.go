package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
)

//go:embed VERSION
var versionFile string

func main() {
	compFlag := flag.Bool("c", false, "Generate complementary colors")
	versionFlag := flag.Bool("v", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("kilour version %s\n", strings.TrimSpace(versionFile))
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go [-c] [-v] <image_path> [output_path]")
		flag.PrintDefaults()
		return
	}
	inputFile := args[0]

	// Determine output path
	var outputPath string
	if len(args) >= 2 {
		outputPath = args[1]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		outputPath = filepath.Join(homeDir, "dots/styles/palette.css")
	}

	// 1. Load the Image
	img, err := loadImage(inputFile)
	if err != nil {
		panic(err)
	}

	// 2. Extract Dominant Colors (K-Means)
	// K=7 to get the 7 most dominant colors.
	// prominentcolor handles resizing internally for speed.
	centroids, err := prominentcolor.KmeansWithAll(7, img, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, prominentcolor.GetDefaultMasks())
	if err != nil {
		panic(err)
	}

	// 3. Prepare CSS Content
	var cssBuilder strings.Builder
	cssBuilder.WriteString(":root {\n")

	// We use a map to handle duplicate names (e.g. if two shades map to "Teal")
	nameCounts := make(map[string]int)

	// Limit to top 7 colors (or fewer if K-means returns less)
	count := 7
	if len(centroids) < count {
		count = len(centroids)
	}

	for i := 0; i < count; i++ {
		c := centroids[i]

		// Convert to colorful.Color
		col := colorful.Color{
			R: float64(c.Color.R) / 255.0,
			G: float64(c.Color.G) / 255.0,
			B: float64(c.Color.B) / 255.0,
		}

		if *compFlag {
			h, s, l := col.Hsl()
			h += 180.0
			if h >= 360.0 {
				h -= 360.0
			}
			col = colorful.Hsl(h, s, l)
		}

		// Convert standard Go color to hex string
		hexCode := col.Hex()

		// Find closest name
		name := findClosestColorName(col)
		
		// Handle duplicates (e.g., deep-blue, deep-blue-2)
		cleanName := toCssVarName(name)
		nameCounts[cleanName]++
		if nameCounts[cleanName] > 1 {
			cleanName = fmt.Sprintf("%s-%d", cleanName, nameCounts[cleanName])
		}

		// Append to CSS string
		cssBuilder.WriteString(fmt.Sprintf("    --%s: %s;\n", cleanName, hexCode))
	}

	cssBuilder.WriteString("}")

	// 4. Write to File
	err = os.WriteFile(outputPath, []byte(cssBuilder.String()), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully generated %s\n", outputPath)
	fmt.Println(cssBuilder.String())
}

// loadImage helper
func loadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// findClosestColorName iterates through the SVG/CSS color list and finds
// the one with the shortest distance in Lab color space.
func findClosestColorName(target colorful.Color) string {
	minDist := math.MaxFloat64
	closestName := "color"

	for name, rgba := range colornames.Map {
		// Convert map color to colorful.Color
		candidate, _ := colorful.MakeColor(rgba)

		// Calculate perception distance (CIE94 is fast and good for this)
		dist := target.DistanceCIE94(candidate)

		if dist < minDist {
			minDist = dist
			closestName = name
		}
	}
	return closestName
}

// toCssVarName formats "DarkOliveGreen" to "dark-olive-green"
func toCssVarName(name string) string {
	// Simple camelCase to kebab-case conversion
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
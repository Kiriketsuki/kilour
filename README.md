# kilour

![Version](https://img.shields.io/badge/version-0.0.2-blue.svg)

kilour is a simple command-line tool written in Go that extracts the dominant colors from an image and generates a CSS file with CSS variables for those colors. It uses K-means clustering to find the most prominent colors and names them using the nearest standard web color name.

## Features

-   Extracts up to 7 dominant colors from an image.
-   Optionally generates complementary colors for the extracted palette.
-   Generates a CSS file defining these colors as CSS variables (e.g., `--dark-olive-green: #556b2f;`).
-   Automatically handles duplicate color names by appending a counter.
-   Defaults to saving the palette to `~/dots/styles/palette.css`.
-   Allows specifying a custom output path.

## Installation

You can run the tool directly using `go run` or compile it into a binary.

### Prerequisites

-   [Go](https://golang.org/dl/) (1.16 or later recommended)

### Build

```bash
git clone https://github.com/Kiriketsuki/kilour.git
cd kilour
go build -o kilour main.go
```

## Usage

### Running directly with Go

```bash
go run main.go [-c] [-v] <path_to_image> [output_path]
```

### Running the binary

```bash
./kilour [-c] [-v] <path_to_image> [output_path]
```

### Arguments

1. **`-c`** (Optional): Generate complementary colors instead of the dominant colors.
2. **`-v`** (Optional): Print the version and exit.
3. **`<path_to_image>`** (Required): The path to the input image file (JPEG or PNG).
4. **`[output_path]`** (Optional): The path where the generated CSS file should be saved.
    - **Default:** `~/dots/styles/palette.css`

## Examples

**1. Generate palette to default location:**

```bash
./kilour my_wallpaper.jpg
```

This will generate `~/dots/styles/palette.css`.

**2. Generate palette to a specific file:**

```bash
./kilour my_wallpaper.jpg ./my_theme.css
```

This will generate `./my_theme.css`.

**3. Generate complementary colors:**

```bash
./kilour -c my_wallpaper.jpg
```

This will generate the complementary colors of the dominant palette to `~/dots/styles/palette.css`.

## Output Format

The generated CSS file looks like this:

```css
:root {
    --dark-slate-gray: #2f4f4f;
    --dim-gray: #696969;
    --silver: #c0c0c0;
    --dark-olive-green: #556b2f;
    --black: #050505;
    --dark-slate-gray-2: #3a5f5f;
    --gray: #808080;
}
```

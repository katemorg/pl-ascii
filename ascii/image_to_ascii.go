package ascii

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

type Ascii struct {
	Filename    string
	AspectRatio string
	Art         string
}

func ImageToAscii(w http.ResponseWriter, r *http.Request) {
	// parse multipart form and get file
	// iterate to allow for multiple files
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse multipart message: %v", err), http.StatusBadRequest)
		return
	}

	var payload []Ascii

	for _, fh := range r.MultipartForm.File["media"] {
		// open file
		file, err := fh.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to open file: %v", err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// decode image
		img, imageType, err := image.Decode(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode image: %v", err), http.StatusBadRequest)
			return
		}

		// check for supported file type
		if imageType != "jpeg" && imageType != "png" {
			http.Error(w, "please upload a jpeg or png", http.StatusBadRequest)
		}

		var asciiArt Ascii

		// TODO: allow client to request differet aspect ratios and inverted colors
		// intensity := " .:-=+*#%@"
		intensity := "@%#*+=-:. "
		min, max := img.Bounds().Min, img.Bounds().Max
		scaleX, scaleY := 8, 4
		output := ""

		// transform pixels
		for y := min.Y; y < max.Y; y += scaleX {
			for x := min.X; x < max.X; x += scaleY {
				c := getAvgPixelColor(img, x, y, scaleX, scaleY, max)
				// 16-bit color format has 65536 possible values
				output = output + string(intensity[len(intensity)*c/65536])
			}
			output = output + string("\n")
		}

		asciiArt = Ascii{
			Filename:    fh.Filename,
			AspectRatio: fmt.Sprintf("%v:%v", scaleX, scaleY),
			Art:         output,
		}

		payload = append(payload, asciiArt)

	}
	if len(payload) > 0 {
		respondJSON(w, payload)
	} else {
		http.Error(w, "no images found", http.StatusBadRequest)
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG!"))
}

// getAvgPixelColor returns the average color of a specified range of pixels.
// This allows us to effectively "resize" the image in ascii art.
func getAvgPixelColor(img image.Image, x, y, width, height int, max image.Point) int {
	count, sum := 0, 0
	for i := x; i < x+width && i < max.X; i++ {
		for j := y; j < y+height && j < max.Y; j++ {
			sum += pixelToGreyscale(img.At(i, j))
			count++
		}
	}
	return (sum / count)
}

// pixelToGreyscale converts pixel RGB values to grayscale using the NTSC formula, which represents the
// average person's relative perception of the brightness of red, green, and blue light.
// Allows us to better assign a more accurate ascii character intensity.
func pixelToGreyscale(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}

// respondJSON makes and sends the response with json payload and default 200 status code.
// Can extend method signature to include cache behavior and other headers.
func respondJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "\n")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal json response: %v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write([]byte(response))
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to write response: %v", err), http.StatusInternalServerError)
	}
}

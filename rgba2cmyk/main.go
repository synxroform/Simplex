package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tiff "github.com/andeling/tiff"
)

func main() {

	fmt.Println("This program replace RGBA photometric tag in tiff image with CMYK.")
	fmt.Println("Only uncompreseed 8bit images supported.")
	fmt.Println("------------------------------------------------------------------")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n tiff file to process: ")
	fname, _ := reader.ReadString('\n')

	file, err := os.Open(strings.TrimSpace(fname))
	if err != nil {
		panic(err)
	}

	dec, _ := tiff.NewDecoder(file)
	itr := dec.Iter()

	for itr.Next() {
		img := itr.Image()
		w, h := img.WidthHeight()
		spp := img.SamplesPerPixel()

		if spp != 4 {
			panic("image must contain four channels")
		}

		buf := make([]uint8, w*h*spp)
		img.DecodeImage(buf)

		fmt.Println("output image to : ")
		oname, _ := reader.ReadString('\n')

		out, _ := os.Create(strings.TrimSpace(oname))
		enc := tiff.NewEncoder(out)
		defer enc.Close()
		defer out.Close()

		img2 := enc.NewImage()
		img2.SetPixelFormat(tiff.PhotometricSeparated, 4, []int{8, 8, 8, 8})
		img2.SetWidthHeight(w, h)
		img2.SetCompression(tiff.CompressionNone)

		bu2 := make([]uint8, w*h*4)
		for n := range bu2 {
			bu2[n] = buf[n]
		}

		img2.EncodeImage(bu2)
	}
	fmt.Println("all done")
}

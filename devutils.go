package go_printer

import "os"

// For debugging purposes ONLY

func MustValue[T any](something T, err error) T {
	if err != nil {
		panic(err)
	}
	return something
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustDeferred(fn func() error) {
	func() {
		Must(fn())
	}()
}

func WritePrintableImageDebugBin(filename string, img *PrintableImage) {
	outFile := MustValue(os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644))
	defer MustDeferred(outFile.Close)

	_, err := outFile.Write([]byte{uint8(img.width), uint8(img.height >> 8), uint8(img.height % 256)})
	Must(err)

	for i := 0; i < img.height; i++ {
		_, err := outFile.Write(*(*img.lines)[i])
		Must(err)
	}
}

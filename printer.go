package go_printer

import (
	"fmt"
	"strings"
)

type UsbPrinterState struct {
	doubleWidth, doubleHeight bool
	underlineLevel            UnderlineLevel
	emphasis                  bool
	invert                    bool
	fontType                  FontType
	alignmentLevel            AlignmentLevel
	leftMargin                uint16
}

type UsbPrinter struct {
	write func(data []byte) (wrote int, err error)
	state *UsbPrinterState
	debug bool
	Done  func()
}

const (
	lineWidthNormal      = 32
	lineWidthSmall       = 42
	cmdESC          byte = 0x1b
	cmdGS           byte = 0x1d
)

// TODO: barcodes / QR codes

func (p *UsbPrinter) initialize() {
	p.Reset()
}

func (p *UsbPrinter) SetDebug(enabled bool) {
	p.debug = enabled
}

func (p *UsbPrinter) GetLineWidth() int {
	switch p.state.fontType {
	case FontSmall:
		return lineWidthSmall
	default:
		return lineWidthNormal
	}
}

func (p *UsbPrinter) WriteOrPanic(data []byte) {
	_, err := p.write(data)
	if p.debug {
		printHexWithFooter(data, fmt.Sprintf("Wrote %d bytes to printer\n", len(data)))
	}
	if err != nil {
		panic(err)
	}
}

func (p *UsbPrinter) EscCommandOrPanic(data ...byte) {
	p.WriteOrPanic(append([]byte{cmdESC}, data...))
}

func (p *UsbPrinter) GsCommandOrPanic(data ...byte) {
	p.WriteOrPanic(append([]byte{cmdGS}, data...))
}

func (p *UsbPrinter) Reset() {
	p.EscCommandOrPanic(0x40)
}

func (p *UsbPrinter) SoftReset() {
	p.Reset()

	state := p.state

	if state.doubleWidth || state.doubleHeight {
		p.Double(state.doubleWidth, state.doubleHeight)
	}
	if state.underlineLevel > 0 {
		p.Underline(state.underlineLevel)
	}
	if state.emphasis {
		p.Emphasis(state.emphasis)
	}
	if state.invert {
		p.Invert(state.invert)
	}
	if state.fontType > 0 {
		p.Font(state.fontType)
	}
	if state.alignmentLevel > 0 {
		p.Alignment(state.alignmentLevel)
	}
	if state.leftMargin > 0 {
		p.LeftMargin(state.leftMargin)
	}
}

func (p *UsbPrinter) Print(text string) {
	p.WriteOrPanic([]byte(text))
}

func (p *UsbPrinter) Println(text string) {
	p.WriteOrPanic(append([]byte(text), '\n'))
}

func (p *UsbPrinter) LineFeed(lines uint8) {
	p.EscCommandOrPanic(0x64, lines)
}

func (p *UsbPrinter) Divider(divider string) {
	var output strings.Builder
	for output.Len() < p.GetLineWidth() {
		output.WriteString(divider)
	}
	p.Println(output.String()[:p.GetLineWidth()])
}

func (p *UsbPrinter) Double(width, height bool) {
	var flags uint8 = 0
	if width {
		flags += 32
	}
	if height {
		flags += 16
	}
	p.EscCommandOrPanic(0x21, flags)
	p.state.doubleWidth = width
	p.state.doubleHeight = height
}

type UnderlineLevel = uint8
type AlignmentLevel = uint8
type FontType = uint8

const (
	UnderlineNone   UnderlineLevel = 0
	UnderlineMedium UnderlineLevel = 1
	UnderlineStrong UnderlineLevel = 2
	AlignmentLeft   AlignmentLevel = 0
	AlignmentCenter AlignmentLevel = 1
	AlignmentRight  AlignmentLevel = 2
	FontNormal      FontType       = 0
	FontSmall       FontType       = 1
)

func (p *UsbPrinter) Underline(level UnderlineLevel) {
	p.EscCommandOrPanic(0x2d, level)
	p.state.underlineLevel = level
}

func (p *UsbPrinter) Emphasis(enabled bool) {
	var flags uint8 = 0
	if enabled {
		flags = 1
	}
	p.EscCommandOrPanic(0x45, flags)
	p.state.emphasis = enabled
}

func (p *UsbPrinter) Invert(enabled bool) {
	var flags uint8 = 0
	if enabled {
		flags = 1
	}
	p.GsCommandOrPanic(0x42, flags)
	p.state.invert = enabled
}

func (p *UsbPrinter) Font(fontType FontType) {
	p.EscCommandOrPanic(0x4d, fontType)
	p.state.fontType = fontType
}

func (p *UsbPrinter) Alignment(level AlignmentLevel) {
	p.EscCommandOrPanic(0x61, level)
	p.state.alignmentLevel = level
}

func (p *UsbPrinter) LeftMargin(margin uint16) {
	p.GsCommandOrPanic(0x4c, byte(margin%0xff), byte(margin>>8))
	p.state.leftMargin = margin
}

func (p *UsbPrinter) PrintImage(img *PrintableImage) {
	buffer := make([]byte, 0, 8+img.width*img.height)

	buffer = append(buffer,
		cmdGS, 0x76, 0x30, 0x00,
		uint8(img.width%256), uint8(img.width>>8),
		uint8(img.height%256), uint8(img.height>>8),
	)

	for _, line := range *img.lines {
		buffer = append(buffer, *line...)
	}

	p.WriteOrPanic(buffer)

	p.SoftReset()
}

func (p *UsbPrinter) PrintText(text string) {
	lines := splitTextIntoLines(text, p.GetLineWidth())

	for _, line := range lines {
		p.Println(line)
	}
}

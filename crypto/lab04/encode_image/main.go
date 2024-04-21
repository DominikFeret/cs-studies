package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"slices"
)

func main() {
	img, s1, s2, perm, key, _ := getInput()
	imgBits := jpegToBits(img)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y

	encImg := encTextEcb(imgBits, s1, s2, perm, key)

	newImg := createImage([]byte(encImg), w, h)
	writeImage(newImg)
}

// cannot run in parallel
func encTextCbc(str []byte, s1, s2 [][]byte, perm []int, key []byte, iv []byte) string {
	enc := ""
	blockSize := 12

	previousBlock := iv
	for i := 0; i < len(str); i += blockSize {
		fmt.Printf("%.02f", float32(i)/float32(len(str))*100)
		fmt.Println("%")
		if i+blockSize > len(str) {
			break
		}

		currentBlock := intToBin12(binToInt(str[i:i+blockSize]) ^ binToInt(previousBlock))

		keycpy := slices.Clone(key)
		previousBlock = []byte(miniDesEnc(currentBlock, s1, s2, perm, keycpy))
		enc += string(previousBlock)
	}
	return enc
}

// can run in parallel
func encTextEcb(str []byte, s1, s2 [][]byte, perm []int, key []byte) string {
	enc := ""
	blockSize := 12

	for i := 0; i < len(str); i += blockSize {
		fmt.Printf("%.02f", float32(i)/float32(len(str))*100)
		fmt.Println("%")
		if i+blockSize > len(str) {
			break
		}
		keycpy := slices.Clone(key)
		enc += miniDesDec(str[i:i+blockSize], s1, s2, perm, keycpy)
	}

	return enc
}

func miniDesDec(str []byte, s1, s2 [][]byte, perm []int, key []byte) string {
	bLen := len(str) / 2
	pl, pr := binToInt(str[:bLen]), binToInt(str[bLen:])
	cl, cr := 0, 0

	for i := 0; i < 7; i++ {
		// shift key slice left by 1 with each iteration
		e := permute(intToBin6(pr), perm)
		shiftRight(key)

		xored := intToBin8(binToInt(e) ^ binToInt(key))

		i1, i2 := binToInt(xored[:4]), binToInt(xored[4:])

		f := binToInt([]byte(string(s1[i1]) + string(s2[i2])))

		cl = pr
		cr = pl ^ f
		pl, pr = cl, cr
	}

	return string(append(intToBin6(cr), intToBin6(cl)...))
}

func miniDesEnc(str []byte, s1, s2 [][]byte, perm []int, key []byte) string {
	bLen := len(str) / 2
	pl, pr := binToInt(str[:bLen]), binToInt(str[bLen:])
	cl, cr := 0, 0

	for i := 0; i < 7; i++ {
		// shift key slice left by 1 with each iteration
		e := permute(intToBin6(pr), perm)
		shiftLeft(key)

		xored := intToBin8(binToInt(e) ^ binToInt(key))

		i1, i2 := binToInt(xored[:4]), binToInt(xored[4:])

		f := binToInt([]byte(string(s1[i1]) + string(s2[i2])))

		cl = pr
		cr = pl ^ f
		pl, pr = cl, cr
	}

	return string(append(intToBin6(cr), intToBin6(cl)...))
}

func writeImage(img image.Image) {
	out, err := os.Create("out.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, img, nil)
}

func createImage(img []byte, w, h int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	i := 0
	for y := range h {
		for x := range w {
			r := binToInt(img[i : i+8])
			g := binToInt(img[i+8 : i+16])
			b := binToInt(img[i+16 : i+24])

			newImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			i += 36
		}
	}

	return newImg
}

func binToInt(bin []byte) int {
	s := string(bin)

	n := 0

	for i := 0; i < len(s); i++ {
		if s[len(s)-i-1] == '1' {
			n += 1 << i
		}
	}

	return int(n)
}

func intToBin6(n int) []byte {
	s := fmt.Sprintf("%06b", n)

	return []byte(s)
}

func intToBin8(n int) []byte {
	s := fmt.Sprintf("%08b", n)

	return []byte(s)
}

func intToBin12(n int) []byte {
	s := fmt.Sprintf("%012b", n)

	return []byte(s)
}

func intToBin(n int) []byte {
	s := fmt.Sprintf("%b", n)

	return []byte(s)
}

func permute(bin []byte, perm []int) []byte {
	permuted := make([]byte, len(perm))
	for i, v := range perm {
		permuted[i] = bin[v]
	}

	return permuted
}

func shiftLeft(key []byte) {
	first := key[0]
	for i := 0; i < len(key)-1; i++ {
		key[i] = key[i+1]
	}
	key[len(key)-1] = first
}

func shiftRight(s []byte) {
	last := s[len(s)-1]
	for i := len(s) - 1; i > 0; i-- {
		s[i] = s[i-1]
	}
	s[0] = last
}

func getInput() (image.Image, [][]byte, [][]byte, []int, []byte, []byte) {
	img := readInput()

	s1 := [][]byte{
		[]byte("101"),
		[]byte("010"),
		[]byte("001"),
		[]byte("110"),
		[]byte("011"),
		[]byte("100"),
		[]byte("111"),
		[]byte("000"),
		[]byte("001"),
		[]byte("100"),
		[]byte("110"),
		[]byte("010"),
		[]byte("000"),
		[]byte("111"),
		[]byte("101"),
		[]byte("011"),
	}
	s2 := [][]byte{
		[]byte("100"),
		[]byte("000"),
		[]byte("110"),
		[]byte("101"),
		[]byte("111"),
		[]byte("001"),
		[]byte("011"),
		[]byte("010"),
		[]byte("101"),
		[]byte("011"),
		[]byte("000"),
		[]byte("111"),
		[]byte("110"),
		[]byte("010"),
		[]byte("001"),
		[]byte("100"),
	}

	k := []byte{'1', '0', '1', '0', '1', '0', '1', '0'}

	p := []int{0, 4, 3, 2, 3, 5, 4, 1}
	iv := []byte{'1', '1', '1', '0', '1', '1', '0', '1', '0', '0', '1', '0'}

	return img, s1, s2, p, k, iv
}

func jpegToBits(img image.Image) []byte {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// each pixel is represented by binary values of 3 colors with the length of 12 bits each
	bits := make([]byte, 0, width*height*3*8)

	for y := range height {
		for x := range width {
			r, g, b, _ := img.At(x, y).RGBA()
			bits = append(bits, intToBin8(int(r))...)
			bits = append(bits, intToBin8(int(g))...)
			bits = append(bits, intToBin8(int(b))...)
		}
	}

	return bits
}

func readInput() image.Image {
	file, err := os.Open("img.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	return img
}

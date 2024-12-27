package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateHexColor() string {
	// Membuat generator angka acak baru dengan sumber yang diinisialisasi dengan waktu saat ini
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Menghasilkan angka acak antara 0 dan 255 untuk setiap komponen RGB
	r := randGen.Intn(256)
	g := randGen.Intn(256)
	b := randGen.Intn(256)

	// Mengubah komponen RGB menjadi format hex dan menggabungkannya
	hexColor := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return hexColor
}

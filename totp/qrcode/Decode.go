package qrcode

import "image"
import "github.com/makiuchi-d/gozxing"
import gozxing_qrcode "github.com/makiuchi-d/gozxing/qrcode"

func Decode(image image.Image) string {

	var result string

	bitmap, err1 := gozxing.NewBinaryBitmapFromImage(image)

	if err1 == nil {

		reader := gozxing_qrcode.NewQRCodeReader()
		data, err2 := reader.Decode(bitmap, nil)

		if err2 == nil {
			result = data.GetText()
		}

	}

	return result

}

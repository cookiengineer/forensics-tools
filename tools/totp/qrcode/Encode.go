package qrcode

import "github.com/makiuchi-d/gozxing"
import "github.com/makiuchi-d/gozxing/qrcode"
import "totp/otpauth"
import "bytes"
import "image/png"

func Encode(otp otpauth.OTP) []byte {

	var result []byte

	data := otp.ToURL()
	writer := qrcode.NewQRCodeWriter()
	hints := map[gozxing.EncodeHintType]interface{}{gozxing.EncodeHintType_CHARACTER_SET: "UTF-8"}
	image, err1 := writer.Encode(data,
		gozxing.BarcodeFormat_QR_CODE,
		1024,
		1024,
		hints,
	)

	if err1 == nil {

		buffer := new(bytes.Buffer)
		err2 := png.Encode(buffer, image)

		if err2 == nil {
			result = buffer.Bytes()
		}

	}

	return result

}

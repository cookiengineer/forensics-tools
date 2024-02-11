package otpauth

import "encoding/base32"
import net_url "net/url"
import "strconv"

type OTP struct {
	Issuer string `json:"issuer"`
	Name   string `json:"name"`
	Secret string `json:"secret"`

	Algorithm string `json:"algorithm"`
	Type      string `json:"type"`
	Digits    int    `json:"digits"`
	Period    uint64 `json:"period"`
}

func NewOTP() OTP {

	var otp OTP

	otp.Secret = "FOOBAR"
	otp.Digits = 6
	otp.Period = 30
	otp.Algorithm = "SHA1"

	return otp

	// TODO: Create secret

}

func (otp *OTP) SetSecret(secret []byte) {
	otp.Secret = base32.StdEncoding.EncodeToString(secret)
}

func (otp *OTP) ToURL() string {

	var url string

	url = "otpauth://" + otp.Type + "/" + otp.Name
	url = url + "?secret=" + otp.Secret

	if otp.Issuer != "" {
		url = url + "&issuer=" + net_url.QueryEscape(otp.Issuer)
	}

	if otp.Algorithm != "" {
		url = url + "&algorithm=" + otp.Algorithm
	}

	if otp.Digits != 0 {
		url = url + "&digits=" + strconv.Itoa(otp.Digits)
	}

	if otp.Period != 0 {
		url = url + "&period=" + strconv.FormatUint(otp.Period, 10)
	}

	return url

}

// otpauth://totp/ACME%20Co:john.doe@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&issuer=ACME%20Co&algorithm=SHA1&digits=6&period=30

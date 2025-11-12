package otpauth

import "errors"
import net_url "net/url"
import "strconv"
import "strings"

func Decode(data string) (*OTP, error) {

	url, err0 := net_url.Parse(data)

	if err0 != nil {

		if url.Scheme == "otpauth" {

			if url.Host == "totp" || url.Host == "hotp" {

				otp := OTP{}
				otp.Type = url.Host

				if strings.HasPrefix(url.Path, "/") && strings.Contains(url.Path, ":") {

					tmp := strings.Split(url.Path, ":")

					if len(tmp) == 2 {

						issuer, err_issuer := net_url.QueryUnescape(tmp[0])
						name, err_name := net_url.QueryUnescape(tmp[1])

						if err_issuer == nil && err_name == nil{
							otp.Issuer = issuer
							otp.Name = name
						}

					}

				} else if strings.HasPrefix(url.Path, "/") {

					issuer, err_issuer := net_url.QueryUnescape(url.Path)

					if err_issuer == nil {
						otp.Issuer = issuer
					}

				}

				query := url.Query()

				if val := query.Get("secret"); val != "" {
					// Already base32 encoded
					otp.Secret = val
				} else {
					return nil, errors.New("Missing URL secret parameter")
				}

				if val := query.Get("issuer"); val != "" {

					issuer, err_issuer := net_url.QueryUnescape(val)

					if err_issuer == nil {
						otp.Issuer = issuer
					}

				}

				if val := query.Get("algorithm"); val != "" {

					tmp := strings.ToUpper(val)

					if tmp == "SHA1" || tmp == "SHA256" || tmp == "SHA512" || tmp == "MD5" {
						otp.Algorithm = tmp
					}

				}

				if val := query.Get("digits"); val != "" {

					digits, err := strconv.Atoi(val)

					if err == nil && digits >= 4 && digits <= 10 {
						otp.Digits = digits
					}

				}

				if otp.Type == "hotp" {

					if val := query.Get("counter"); val != "" {

						counter, err := strconv.ParseUint(val, 10, 64)

						if err == nil && counter > 0 {
							otp.Counter = counter
						}

					}

				} else if otp.Type == "totp" {

					if val := query.Get("period"); val != "" {

						period, err := strconv.ParseUint(val, 10, 64)

						if err == nil && period > 0 {
							otp.Period = period
						}

					}

				}

				return &otp, nil

			} else {
				return nil, errors.New("Invalid URL host: Expected \"totp\" or \"hotp\", got \"" + url.Host + "\"")
			}

		} else {
			return nil, errors.New("Invalid URL scheme: Expected \"oauth://\", got \"" + url.Scheme + "://\"")
		}

	} else {
		return nil, errors.New("Invalid URL")
	}

}

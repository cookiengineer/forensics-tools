package otpauth

import "github.com/dim13/otpauth/migration"

func DecodeMigration(data string) []OTP {

	var result []OTP

	decoded, err1 := migration.Data(data)

	if err1 == nil {

		payload, err2 := migration.Unmarshal(decoded)

		if err2 == nil {

			if len(payload.OtpParameters) > 0 {

				for o := 0; o < len(payload.OtpParameters); o++ {

					var otp OTP

					parameters := payload.OtpParameters[o]

					if parameters.Algorithm.Name() == "SHA1" {
						otp.Algorithm = "SHA1"
					} else if parameters.Algorithm.Name() == "SHA256" {
						otp.Algorithm = "SHA256"
					} else if parameters.Algorithm.Name() == "SHA512" {
						otp.Algorithm = "SHA512"
					} else if parameters.Algorithm.Name() == "MD5" {
						otp.Algorithm = "MD5"
					}

					if len(parameters.Secret) > 0 {
						otp.SetSecret(parameters.Secret)
					}

					if parameters.Type.Name() == "hotp" {
						otp.Type = "hotp"
					} else if parameters.Type.Name() == "totp" {
						otp.Type = "totp"
					}

					if parameters.Counter > 0 {
						otp.Period = parameters.Counter
					}

					if parameters.Digits.Count() == 6 {
						otp.Digits = 6
					} else if parameters.Digits.Count() == 8 {
						otp.Digits = 8
					}

					if parameters.Name != "" {
						otp.Name = parameters.Name
					}

					if parameters.Issuer != "" {
						otp.Issuer = parameters.Issuer
					}

					result = append(result, otp)

				}

			}

		}

	}

	return result

}

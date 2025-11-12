package otpauth

import "github.com/dim13/otpauth/migration"

func DecodeMigration(data string) []*OTP {

	result := make([]*OTP, 0)

	decoded, err1 := migration.Data(data)

	if err1 == nil {

		payload, err2 := migration.Unmarshal(decoded)

		if err2 == nil {

			if len(payload.OtpParameters) > 0 {

				for o := 0; o < len(payload.OtpParameters); o++ {

					otp := NewOTP()
					parameters := payload.OtpParameters[o]

					algorithm := parameters.Algorithm.Name()

					if algorithm == "SHA1" || algorithm == "SHA256" || algorithm == "SHA512" || algorithm == "MD5" {
						otp.Algorithm = algorithm
					}

					if len(parameters.Secret) > 0 {
						otp.SetSecret(parameters.Secret)
					}

					typ := parameters.Type.Name();

					if typ == "hotp" {

						otp.Type = typ

						if parameters.Counter > 0 {
							otp.Counter = parameters.Counter
						}

					} else if typ == "totp" {

						otp.Type = typ

						// XXX: Bug in upstream otpauth implementation schema, missing Period property
						if parameters.Counter > 0 {
							otp.Period = parameters.Counter
						}

					}

					if digits := parameters.Digits.Count(); digits >= 4 && digits <= 10 {
						otp.Digits = digits
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

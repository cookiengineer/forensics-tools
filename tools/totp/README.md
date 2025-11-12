
# TOTP Tools

The TOTP tools are useful to extract and generate QR code images from camera photos of 2FA export
dialogs, for example from a smartphone or computer screen. The TOTP tools support both the Google
Authenticator `otp-migration://` URL scheme and the `otp-auth://` URL scheme that's provided in
most 2FA secret dialogs that show QR codes for import in a password manager.

## Dependencies

- [gozxing](https://github.com/makiuchi-d/gozxing)
- [otpauth](https://github.com/dim13/otpauth)

## Usage

```bash
# Extract TOTP qrcodes into totp-qrcode-<number>-<issuer>.png and totp-qrcode-<number>-<issuer>.json
totp-extract ./DCIM_1337_google_authenticator_dialog.jpeg;

# Extract TOTP qrcode into totp-qrcode-<issuer>.png and totp-qrcode-<issuer>.json
totp-extract ./DCIM_1338_keepass_dialog.jpeg;
```

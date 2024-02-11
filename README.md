
# Forensics Tools

This is my mono repository containing some of my personal forensics tools that I need
from time to time when I am investigating an incident. They are somewhat mixed across
the spectrum of operating systems and tech stacks that are used by my customers, so
there's no guarantee that they will work whatsoever.


## Git Tools

The [Git Tools](./git-tools) are useful for OSINT gathering on public repositories,
such as on BitBucket, GitHub or GitLab.

## SQL Tools

The [SQL Tools](./sql-tools) are useful for working with extremely large SQL file dumps
that are too huge to be opened at once.

```bash
sql-tables large-dump.sql;             # list of table names
sql-extract large-dump.sql table-name; # extracts a specific table and its data
```

## Torrent Tools

The [Torrent Tools](./torrent-tools) allow to inspect and modify `magnet:` URLs.

## TOTP Tools

The [TOTP Tools](totp-tools) allow to export encoded `otp-migration://` 2FA seeds
from screenshots or camera photos into a `json` file and a QRCode `png` that can
be used to import them within other password managers.

```bash
totp-extract ./path/to/camera-photo-of-qrcode.jpg;
```

## Chrome/Chromium Extension Extractor

The [uncrx](./uncrx) tool is useful for extracting a `.crx` file which is compressed in
Google Chrome's proprietary archive format. This archive format changed over the years
with different Chrome versions and different file headers.

```bash
export EXTENSION_ID="cjpalhdlnbpafiamejdnhcphjbkeiagm";
export EXTENSION_NAME="ublock-origin";

wget -O "$EXTENSION_NAME.crx" "https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion=100&x=id%3D$EXTENSION_ID%26uc";

uncrx "$EXTENSION_NAME.crx":                        # creates the $EXTENSION_NAME.zip file in the same folder
unzip "$EXTENSION_NAME.zip" -d "./$EXTENSION_NAME"; # unpack the extension, so that it can be loaded in Developer Mode
```


## License

GPL3


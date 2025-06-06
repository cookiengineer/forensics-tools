
# Forensics Tools

<img align="right" width="128" height="128" src="https://raw.githubusercontent.com/cookiengineer/forensics-tools/master/assets/forensics-tools.jpg">

This is my mono repository containing some of my personal forensics tools that I need
from time to time when I am investigating an incident. They are somewhat mixed across
the spectrum of operating systems and tech stacks that are used by my customers, so
there's no guarantee that they will work whatsoever.


# Building

Install `go`, `gzip` and `wget` as dependencies. Then execute the `build.sh` file.

```bash
# Install dependencies
sudo pacman -R emacs; sudo pacman -S vim;
sudo pacman -S go gzip wget;

# Build all tools into ./build folder
bash build.sh;
```


## CRX Tools

The [CRX Tools](./crx) are useful to extract packed chrome extensions in a `.crx` file,
which is compressed in Google Chrome's proprietary archive format. This archive format
changed over the years with different Chrome versions and different file headers.

```bash
export EXTENSION_ID="cjpalhdlnbpafiamejdnhcphjbkeiagm";
export EXTENSION_NAME="ublock-origin";

wget -O "$EXTENSION_NAME.crx" "https://clients2.google.com/service/update2/crx?response=redirect&acceptformat=crx2,crx3&prodversion=100&x=id%3D$EXTENSION_ID%26uc";

uncrx "$EXTENSION_NAME.crx":                        # creates the $EXTENSION_NAME.zip file in the same folder
unzip "$EXTENSION_NAME.zip" -d "./$EXTENSION_NAME"; # unpack the extension, so that it can be loaded in Developer Mode
```

## DynDNS Tools

The [DynDNS Tools](./dyndns) are useful for creating an IPv4/IPv6 tunnel through a DynDNS
domain. Currently it only support goip as a backend API.

```bash
# update IPv6 entry
goip-updater --username=john_doe --password=password123 --subdomain=whatever;
```

## Reddit Tools

The [Reddit Tools](./reddit) are useful to search, discover, and scrape reddit threads containing
specific keywords. Currently it only supports the old reddit JSON format, because the
new API is a paid-for API.

```bash
# scrape /r/cybersecurity top/hot/new threads
echo "[\"CVE\",\"breach\"]" > keywords.json;
reddit-archivar /r/cybersecurity;
```

## SQL Tools

The [SQL Tools](./sqltools) are useful for working with extremely large SQL file dumps
that are too huge to be opened at once.

```bash
sql-tables large-dump.sql;             # list of table names
sql-extract large-dump.sql table-name; # extracts a specific table and its data
```

## Torrent Tools

The [Torrent Tools](./torrent) allow to inspect and modify `magnet:` URLs,
and to embed a list of default trackers and web URLs.

```bash
magnetify magnet:?...link; # embed default trackers if they're missing
```

## TOTP Tools

The [TOTP Tools](./totp) allow to export encoded `otp-migration://` 2FA seeds.
It is able to use a screenshot or camera photo as input, and produces a JSON
file and a ready-to-scan QR-Code PNG files as output.

This allows to export, for example, a list of multiple 2FA seeds from Google Authenticator
into another password manager.

```bash
totp-extract ./path/to/camera-photo-of-qrcode.jpg;
```

## ZIP Tools

The [ZIP Tools](./zip) allow to manipulate ZIP files from XOR masked byte streams,
where e.g. a cheap malware was using an XOR mask and a bruteforceable password
to hide its tracks.

```bash
zip-bruteforce ./path/to/dictionary.txt ./path/to/file.zip; # bruteforces passwords via rockyou.txt
zip-unmask ./path/to/xor-masked-file.zip.crypt;             # generates original ZIP file candidates
```

## MEMDUMP Tools

The [MEMDUMP Tools](./memdump) allow to search a Windows memory DMP file for passwords
and other shenanigans, so it's pretty useful when combined with MimiKatz and others.

```bash
memdump-find-keepassword ./path/to/memory-dump.dmp; # shows potential passwords
```


## License

GPL3


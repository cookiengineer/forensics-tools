
# Cookie Engineer's Forensics Tools

<p align="center">
    <img width="256" height="256" src="https://raw.githubusercontent.com/cookiengineer/forensics-tools/master/assets/forensics-tools.jpg">
</p>

This is my mono repository containing my personal forensics tools that I need when I'm
investigating an incident. They are something like a better `bashrc`, implemented in Go,
without any kind of warranty to work at all.

The [tools](/tools) folder is separated by use-case. Each of the tools' root folder contains
a `README.md` explaining the purpose of the forensics tools inside them. Make sure to read
them carefully.

The [toolchain](/toolchain) folder contains one entry point:

- The [build.go](/toolchain/build.go) which builds all binaries and a separate `install-forensics-tools` program.


## Tools / Features

- [ ] [archive-pack](/tools/archive/cmds/archive-pack/main.go) packs any known archive files
- [ ] [archive-unpack](/tools/archive/cmds/archive-unpack/main.go) unpacks any known archive files
- [x] [crx-dl](/tools/crx/cmds/crx-dl/main.go) downloads Chromium extension files
- [x] [crx-extract](/tools/crx/cmds/crx-extract/main.go) extracts Chromium extension files
- [x] [dns-iscensored](/tools/dns/cmds/dns-iscensored/main.go) checks whether a domain is censored
- [ ] [ffmpeg-to720p](/tools/ffmpeg/cmds/ffmpeg-to720p/main.go) converts videos to x264 720p videos
- [ ] [ffmpeg-to1080p](/tools/ffmpeg/cmds/ffmpeg-to1080p/main.go) converts videos to x264 1080p videos
- [ ] [ffmpeg-tomp3](/tools/ffmpeg/cmds/ffmpeg-tomp3/main.go) converts videos to mp3 files
- [ ] [git-serve](/tools/git/cmds/git-serve/main.go) serves a local git server
- [x] [dyndns-goip](/tools/dyndns/cmds/dyndns-goip/main.go) updates `goip.de` DynDNS domains
- [ ] [gs-totiff](/tools/gs/cmds/gs-totiff/main.go) converts documents to tiff images
- [x] [http-serve](/tools/http/cmds/http-serve/main.go) serves a folder via HTTP
- [x] [memdump-keepass](/tools/memdump/cmds/memdump-keepass/main.go) finds a KeePass(XC) password in memory dump files
- [x] [npm-dl](/tools/npm/cmds/npm-dl/main.go) downloads and extracts specific packages from the NPM registry
- [ ] [reddit-archive](/tools/reddit/cmds/reddit-archive/main.go) downloads subreddits and threads
- [x] [sql-extract](/tools/sql/cmds/sql-extract/main.go) extracts a specific table from large SQL dump files
- [x] [sql-tables](/tools/sql/cmds/sql-tables/main.go) lists a table index of large SQL dump files
- [x] [torrent-magnetify](/tools/torrent/cmds/torrent-magnetify/main.go) adds default trackers to torrent magnet links
- [x] [totp-extract](/tools/totp/cmds/totp-extract/main.go) extracts OTP password seeds from screenshots or camera photos of QR codes
- [ ] [yt-mp3](/tools/yt-dlp/cmds/yt-mp3/main.go) downloads streams as MP3 files
- [ ] [yt-mp4](/tools/yt-dlp/cmds/yt-mp4/main.go) downloads streams as MP4 files
- [ ] [yt-opus](/tools/yt-dlp/cmds/yt-opus/main.go) downloads streams as OPUS files
- [x] [zip-bruteforce](/tools/totp/cmds/zip-bruteforce/main.go) bruteforces the password of a ZIP file
- [x] [zip-unmask](/tools/totp/cmds/zip-unmask/main.go) unmasks ZIP files that have been XOR obfuscated


## Building

The [build.go](/toolchain/build.go) script builds all tool binaries and a separate `install-forensics-tools`
program that can be deployed to another machine and executed there to install all the contained binaries.

```bash
# Build all tools and the installer
cd /path/to/forensics-tools/toolchain;
go run build.go;
```

Alternatively, you can also simply use `go install` on any of the tools in an isolated capacity:

```bash
cd /path/to/forensics-tools;

cd /tools/npm;
go install ./cmds/npm-dl;
```


## Installation

```bash
# Build the installer which contains all tool binaries
cd /path/to/forensics-tools/toolchain;
go run build.go;

# Install all binaries to /usr/local/bin
cd /path/to/forensics-tools/build;
export PREFIX="/usr/local"; sudo install-forensics-tools;
```


## License

GPL3


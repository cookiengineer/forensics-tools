
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

The [toolchain](/toolchain) folder contains two entry points:

- The [build.go](/toolchain/build.go) which builds all binaries and a separate `install-forensics-tools` program.
- The [install.go](/toolchain/install.go) which builds and installs a specific binary of a `<tool>/<cmd>` path.


## Tools / Features

- [archive-pack](/tools/archive/cmds/archive-pack/main.go) packs any known archive files
- [archive-unpack](/tools/archive/cmds/archive-unpack/main.go) unpacks any known archive files
- [crx-dl](/tools/crx/cmds/crx-dl/main.go) downloads Chromium extension files
- [crx-extract](/tools/crx/cmds/crx-extract/main.go) extracts Chromium extension files
- [ffmpeg-to720p](/tools/ffmpeg/cmds/ffmpeg-to720p) converts videos to x264 720p videos
- [ffmpeg-to1080p](/tools/ffmpeg/cmds/ffmpeg-to1080p) converts videos to x264 1080p videos
- [ffmpeg-tomp3](/tools/ffmpeg/cmds/ffmpeg-tomp3) converts videos to mp3 files
- [git-serve](/tools/git/cmds/git-extract/main.go) serves a local git server
- [dyndns-goip](/tools/dyndns/cmds/dyndns-goip/main.go) updates `goip.de` DynDNS domains
- [gs-totiff](/tools/gs/cmds/gs-totiff/main.go) converts documents to tiff images
- [http-serve](/tools/http/cmds/http-serve/main.go) serves a folder via HTTP
- [memdump-keepass](/tools/memdump/cmds/memdump-keepass/main.go) finds a KeePass(XC) password in memory dump files
- [npm-dl](/tools/npm/cmds/npm-dl/main.go) downloads and extracts specific package versions from NPM
- [reddit-archive](/tools/reddit/cmds/reddit-archive/main.go) downloads subreddits and threads
- [sql-extract](/tools/sql/cmds/sql-extract/main.go) extracts a specific table from SQL dump files
- [torrent-magnetify](/tools/torrent/cmds/torrent-magnetify/main.go) adds default trackers to torrent magnet links
- [totp-extract](/tools/totp/cmds/totp-extract/main.go) extracts OTP password seeds from screenshots or camera photos of QR codes
- [youtube-mp3](/tools/yt-dlp/cmds/youtube-mp3/main.go) downloads streams as MP3 files
- [youtube-mp4](/tools/yt-dlp/cmds/youtube-mp4/main.go) downloads streams as MP4 files
- [youtube-opus](/tools/yt-dlp/cmds/youtube-opus/main.go) downloads streams as OPUS files
- [zip-bruteforce](/tools/totp/cmds/zip-bruteforce/main.go) bruteforces the password of a ZIP file
- [zip-unmask](/tools/totp/cmds/zip-unmask/main.go) unmasks ZIP files that have been XOR obfuscated


## License

GPL3


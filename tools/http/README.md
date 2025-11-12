
# HTTP Tools

The HTTP tools are useful when you want to share or modify things via a local webserver but you don't
want to install a REPL for that. It's also great for exfil purposes, as HTTP usually is pretty much
ignored by Blueteams.

## Dependencies

none (Pure Go implementation)

## Usage

```bash
# Serve anything for other devices
cd ~/Videos;
http-serve 8080;

xdg-open "http://localhost:8080";
```

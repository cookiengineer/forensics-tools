
# Torrent Tools

The Torrent tools are useful when you receive an outdated `magnet:` link and the initially seeding
torrent tracker isn't online anymore. For this purpose, the `torrent-magnetify` tool allows to embed
all known default torrent trackers, so that chances increase of being able to download the (censored
or taken down) torrent.

## Dependencies

none (Pure Go implementation)

## Usage

```bash
# Output the fixed magnet link with more torrent trackers
torrent-magnetify "magnet:?xt=urn:btih:a55ac5f0580c53777cfa765d1f86e50dcac50fc0";
```

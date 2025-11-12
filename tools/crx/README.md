
# CRX Tools

The CRX tools are useful to download and extract packed Chrome Extensions from the proprietary format
into a ZIP file that you can either open as a Developer Mode Extension locally or inspect manually
for forensics purposes.

## Dependencies

none (Pure Go implementation)

## Usage

```bash
# Download uBlock Origin Lite via Chrome WebStore URL
crx-dl "https://chromewebstore.google.com/detail/ublock-origin-lite/ddkjiahejlhfcafbddmgiahcphecmpfh?hl=en&pli=1";

# Alternatively: Download uBlock Origin Lite via Extension ID
# crx-dl "ddkjiahejlhfcafbddmgiahcphecmpfh";

# Extract CRX file into ZIP file
crx-extract "ddkjiahejlhfcafbddmgiahcphecmpfh.crx";

# Open ZIP file
xdg-open "ddkjiahejlhfcafbddmgiahcphecmpfh.zip";
```

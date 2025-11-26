
# NPM Tools

The NPM tools are useful to download and extract NPM packages without having dangerous executions
that could compromise your host machine. It will also show dangerous lifecycle hooks that could
potentially execute malware or worms.

## Dependencies

- [brotli](https://github.com/andybalholm/brotli)

## Usage

```bash
# Download unscoped package and Shai Hulud worm sample
npm-dl web3-providers-http 4.1.0;

# Download scoped package and Shai Hulud worm sample
npm-dl @accordproject concerto-analysis 3.24.1;
```

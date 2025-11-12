
# DynDNS Tools

The DynDNS tools are useful for updating the IPv4/IPv6 entries of DynDNS subdomains to the ones that
the current local machine has. In practice, this isn't failsafe and can be detected by Blueteams,
that's why the tools default to using IPv6 for exfil purposes.

### Dependencies

none (Pure Go implementation)

## Usage

```bash
# Update AAAA entry to machine's current IPv6
dyndns-goip john_doe "password123#secure!" johnsonlineshop.goip.de;
```

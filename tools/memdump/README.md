
# MEMDUMP Tools

The MEMDUMP Tools are useful for investigating a memory dump that was provided by WinDbg, dd, or
other tools that use a 1:1 raw memory format without encapsulation.

## Dependencies

none (Pure Go implementation)

## Usage

```bash
# Find KeePass(XC) password in memory dump
memdump-keepass ~/Downloads/evidence-1337.mem;
```

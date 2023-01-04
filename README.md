# gnabglib-go

Basic tools for any app

### Checksum

- Block check character (BCC)
- Fletcher (16,32,64)
- Longitudinal redundancy check (LRC)
- Luhn

### CodeGen

- BytesToHexSep: Format a byte slice in rows of `bytesPerSection` values, formatted in hexadecimal format
- BytesToString: Format a byte slice as a utf8 string, useful for constants in go
- BytesToStringSep: Format a byte slice in rows of `bytesPerSection` utf8 strings

### Encoding

- hex: Convert byte slices to/from hex strings.  Includes a tag:tiny version that doesn't use a 256 byte lookup table for use on embedded devices (~50% slower than regular). Similar to go's built-in hex encoded, except errors include location and value of invalid hex-values on decode.

### Endian

- Detect platform endianness

### Hash

- RipeMD (128,160,256,320): For secure hashing RipeMD 128/256 are no longer recommended
    [Preimage Attacks on Step-Reduced RIPEMD-128 and RIPEMD-160](https://link.springer.com/chapter/10.1007/978-3-642-21518-6_13)
- Streebog (256,512): Subject to a [rebound attack](https://www.sciencedirect.com/science/article/abs/pii/S0020019014001458?via%3Dihub) and [second-preimage attack](https://eprint.iacr.org/2014/675)
- Whirlpool: Subject to a [rebound attack](https://www.iacr.org/archive/fse2009/56650270/56650270.pdf)

### Net

- Cidr, Ipv4, Mask types
- IpTree: Add CIDR and IP addresses to a collection and get back the shortest description of the composition (repeat/overlapping CIDR will merge, sequential CIDR will join into larger sets).  Useful for eg firewall rules

## Testing

`go test ./...`
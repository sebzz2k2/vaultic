# Vaultic Storage Encoding Format
## Data Encoding Layout

Each entry in the key-value (KV) store is encoded using the following structured format:

```
4 bytes   | Total length of the entry (including this field)
1 byte    | Encoding format version
1 byte    | Flags for metadata or additional features
4 bytes   | CRC32 checksum (computed over the key and value)
8 bytes   | Timestamp (Unix epoch time in milliseconds)
2 bytes   | Key length (in bytes)
4 bytes   | Value length (in bytes)
<key>     | Key data (variable length, as specified by the key length)
<value>   | Value data (variable length, as specified by the value length)
```

### Flags

The `Flags` byte is structured as follows:

- **Bit 0**: Indicates if the entry is deleted (`0` = not deleted, `1` = deleted).
- **Bit 1**: Indicates if compression is applied (`0` = uncompressed, `1` = compressed).
- **Bits 2-7**: Reserved for future use (currently set to `0`).

This layout ensures efficient storage and retrieval while maintaining data integrity and flexibility for future enhancements.
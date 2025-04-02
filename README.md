# Vaultic Store Encoding Format

## Data Encoding Layout
Each entry in the KV store is stored in the following structured format:

```
| Version (1 byte) | Flags (1 byte) | Checksum (4 bytes) | Timestamp (8 bytes) | Key Size (2 bytes) | Value Size (4 bytes) | Key (N bytes) | Value (M bytes) |
```

## Field Breakdown
### **1. Version (1 byte)**
- Defines the **encoding format version**, allowing future upgrades.
- Uses only **1 byte (0-255 values)**, as frequent format changes are unlikely.

### **2. Flags (1 byte)**
A bit-packed field storing various boolean properties:
- **Bit 0**: Deleted flag (`0 = Active`, `1 = Deleted` - Tombstone).
- **Bit 1**: Compression flag (`0 = Uncompressed`, `1 = Compressed`).
- **Bits 2-7**: Reserved for future use.

### **3. Checksum (4 bytes)**
- **CRC32 checksum** to verify data integrity.
- Helps detect corruption before parsing variable-length data.

### **4. Timestamp (8 bytes)**
- **Unix epoch timestamp (`int64`)**, storing last modification time.
- Useful for conflict resolution and TTL-based deletions.

### **5. Key Size (2 bytes)**
- **Unsigned 16-bit integer (`0-65535`)**, indicating the length of the key.

### **6. Value Size (4 bytes)**
- **Unsigned 32-bit integer (`0-4,294,967,295`)**, representing value length.
- Supports large values while keeping metadata compact.

### **7. Key (N bytes)**
- Variable-length **key** (`N` bytes, defined by Key Size).

### **8. Value (M bytes)**
- Variable-length **value** (`M` bytes, defined by Value Size).

## Design Justifications
### **Why Place Fixed-Size Fields First?**
- Fixed-size fields (`Version, Flags, Checksum, Timestamp, Key Size, Value Size`) are placed at the beginning for **efficient parsing**.
- Variable-length fields (`Key, Value`) come last to avoid complex scanning.

### **Why Use a Flags Byte?**
- Saves space by packing multiple boolean flags into **1 byte**.
- Allows future extensibility without altering the structure.

### **Why CRC32 for Checksum?**
- Provides a fast and effective method for **corruption detection**.
- **4 bytes is a good balance** between accuracy and overhead.

### **Why Use a Timestamp?**
- Enables **time-based conflict resolution**.
- Supports **TTL-based expirations**.

## Future Enhancements
- **Compression Support:** Add Snappy/Zstd compression.
- **Encryption Flag:** Indicate whether data is encrypted.
- **Sharding Support:** Include partitioning metadata.

This encoding format ensures **efficient, future-proof, and resilient storage** for a high-performance key-value store.


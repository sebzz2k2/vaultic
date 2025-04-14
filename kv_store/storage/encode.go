package storage

import "github.com/sebzz2k2/vaultic/utils"

/*
0th bit deleted 0 or 1
1st bit version 0 or 1
3rd - 7th bit reserved 0 for now
*/
func encodeFlags(flags ...bool) byte {
	var encoded byte
	for i, flag := range flags {
		if flag {
			encoded |= 1 << i
		}
	}
	return encoded
}

/*
4 bytes length
1 byte version
1 byte flags
4 bytes key+value CRC
8 bytes timestamp
2 bytes key length
4 bytes value length
<key length> bytes key
<value length> bytes value
*/
func EncodeData(version int, deleted bool, ts uint64, checkpoint bool, key, value string) ([]byte, int) {
	keyLen := uint16(len(key))     // 2 bytes for key length
	valueLen := uint32(len(value)) // 4 bytes for value length

	// Calculate total length (including the 4-byte length field)
	totalLength := 4 + 1 + 1 + 4 + 8 + 2 + 4 + len(key) + len(value) // 4 bytes for length field

	encoded := make([]byte, 0, totalLength) // Preallocate memory

	// Store total length (first 4 bytes)
	encoded = append(encoded,
		byte(totalLength>>24), byte(totalLength>>16), byte(totalLength>>8), byte(totalLength))

	flags := encodeFlags(deleted, false, checkpoint)
	encoded = append(encoded, byte(version))
	encoded = append(encoded, flags)

	// Store keyValCRC (4 bytes)
	keyValCRC := utils.Crc32(key + value)
	encoded = append(encoded, byte(keyValCRC>>24), byte(keyValCRC>>16), byte(keyValCRC>>8), byte(keyValCRC))

	// Store timestamp (8 bytes)
	encoded = append(encoded,
		byte(ts>>56), byte(ts>>48), byte(ts>>40), byte(ts>>32),
		byte(ts>>24), byte(ts>>16), byte(ts>>8), byte(ts))

	// Store key length (2 bytes)
	encoded = append(encoded, byte(keyLen>>8), byte(keyLen))

	// Store value length (4 bytes)
	encoded = append(encoded,
		byte(valueLen>>24), byte(valueLen>>16), byte(valueLen>>8), byte(valueLen))

	// Append key and value
	encoded = append(encoded, key...)
	encoded = append(encoded, value...)

	return encoded, totalLength
}

# Memtable + SSTable (LSM Base) TODO

## 1️⃣ Memtable
- [x] Choose data structure:
  - Skiplist (better incremental inserts)
  - Map + sort on flush (simpler start)
- [x] Implement CRUD:
  - [x] `GET` → search memtable
  - [x] `SET` → insert/update key
  - [x] `DEL` → insert tombstone marker
- [x] Track memory usage in bytes
- [ ] Trigger flush when size limit reached

---

## 2️⃣ Flush Logic
- [ ] Sort entries (if not already sorted)
- [ ] Write to SSTable:
  - Immutable file
  - Sequential write for fast reads
- [ ] Create new empty memtable after flush
- [ ] Optionally keep immutable memtable during flush for reads

---

## 3️⃣ SSTable Format
- [ ] Store sorted key-value pairs (binary format)
- [ ] Add index block at file end for fast lookup
- [ ] Store metadata:
  - Min/max key
  - Entry count
  - Offset table for block seeks
- [ ] Append Bloom filter for that SSTable

---

## 4️⃣ Read Path
- [ ] Search order:
  1. Active memtable
  2. Immutable memtable (if flush in progress)
  3. SSTables (newest → oldest)
- [ ] Use Bloom filter to skip non-existent keys
- [ ] Merge results (last write wins, tombstones respected)

---

## 5️⃣ Compaction
- [ ] Merge multiple SSTables into one:
  - Drop duplicate keys (keep latest)
  - Remove tombstoned keys older than TTL
- [ ] Replace old SSTables atomically
- [ ] Choose compaction strategy:
  - Levelled (LevelDB-style)
  - Size-tiered (Cassandra-style)

---

## 6️⃣ Metadata & Manifest
- [ ] Maintain manifest file with:
  - SSTable list + levels
  - Bloom filter locations
  - Last sequence number
- [ ] Load manifest on startup

---

## 7️⃣ Optional Enhancements
- [ ] Background compaction thread
- [ ] CRC checksums for data blocks
- [ ] Compression (Snappy/LZ4) for blocks
- [ ] Metrics for flush time, read latency, compaction stats
package kvstore

import (
	"sync"

	"math/rand"
)

// SkipListNode represents a node in the skip list. Each node contains a key,
// a value, and an array of forward pointers to other nodes at different levels.
// The number of forward pointers is determined by the height of the node,
// which is randomly assigned during insertion. This allows for efficient
// traversal of the skip list.
//
// Fields:
//   - Key: The unique identifier for the node.
//   - Value: The value associated with the key.
//   - Next: An array of pointers to the next nodes at different levels.
type SkipListNode struct {
	Key   string
	Value string
	Next  []*SkipListNode
}

// SkipList represents a probabilistic data structure that allows for fast
// search, insertion, and deletion operations. It is composed of multiple
// levels, where each level is a sorted linked list. Higher levels allow
// for skipping over larger portions of the list, improving efficiency.
//
// Fields:
//   - Head: A pointer to the head node of the skip list.
//   - Height: The total number of levels in the skip list. This determines
//     the maximum number of forward pointers a node can have.
//   - Length: The total number of elements currently stored in the skip list.
//   - Level: The current highest level in the skip list. This determines the
//     height of the tallest "tower" of nodes in the structure.
//   - Mutex: A mutex used to ensure thread-safe operations on the skip list.
type SkipList struct {
	Head   *SkipListNode
	Height int
	Length int
	Level  int
	Mutex  *sync.Mutex
}

// NewSkipList creates a new skip list with the specified height and initializes
// the head node. The head node serves as a sentinel node that simplifies
// insertion and deletion operations. The height of the skip list determines
// the maximum number of levels in the structure.
// The skip list is initialized with a height of 1, and the head node's
// forward pointers are set to nil. The length of the skip list is also
// initialized to 0, indicating that it is empty.
func NewSkipList(height int) *SkipList {
	head := &SkipListNode{
		Key:   "",
		Value: "",
		Next:  make([]*SkipListNode, height),
	}
	return &SkipList{
		Head:   head,
		Height: height,
		Length: 0,
		Level:  1,
		Mutex:  &sync.Mutex{},
	}
}

// randomHeight generates a random height for a new node in the skip list.
// The height is determined by flipping a coin until it lands on tails.
// The maximum height is limited to the height of the skip list.
// This probabilistic approach helps maintain a balanced structure,
// ensuring that the skip list remains efficient for search, insertion,
// and deletion operations.
func (s *SkipList) randomHeight() int {
	height := 1
	for height < s.Height && rand.Intn(2) == 0 {
		height++
	}
	return height
}

// Insert adds a new key-value pair to the skip list. The function first
// acquires a lock to ensure thread safety. It then searches for the
// appropriate position to insert the new node. If the key already exists,
// the value is updated. If the key does not exist, a new node is created
// with a random height. The new node's forward pointers are set to point
// to the appropriate nodes in the skip list. The function updates the
// forward pointers of the nodes that precede the new node at each level.
// Finally, the length of the skip list is incremented.
func (s *SkipList) Insert(key, value string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	update := make([]*SkipListNode, s.Height)
	current := s.Head

	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}
		update[i] = current
	}

	if current.Next[0] != nil && current.Next[0].Key == key {
		current.Next[0].Value = value
		return
	}

	newHeight := s.randomHeight()
	if newHeight > s.Level {
		for i := s.Level; i < newHeight; i++ {
			update[i] = s.Head
		}
		s.Level = newHeight
	}

	newNode := &SkipListNode{
		Key:   key,
		Value: value,
		Next:  make([]*SkipListNode, newHeight),
	}

	for i := 0; i < newHeight; i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}

	s.Length++
}

// Search retrieves the value associated with a given key in the skip list.
// The function first acquires a lock to ensure thread safety. It then
// traverses the skip list, moving through the levels and nodes until
// it finds the desired key. If the key is found, the associated value
// is returned. If the key is not found, an empty string and a boolean
// indicating the absence of the key are returned.
func (s *SkipList) Search(key string) (string, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	current := s.Head
	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}
	}

	current = current.Next[0]
	if current != nil && current.Key == key {
		return current.Value, true
	}
	return "", false
}

// Delete removes a key-value pair from the skip list. The function first
// acquires a lock to ensure thread safety. It then searches for the
// specified key and updates the forward pointers of the nodes that
// precede the key. If the key is found, it is removed from the skip list.
// The function also updates the levels of the skip list if necessary.
// Finally, the length of the skip list is decremented.
func (s *SkipList) Delete(key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	update := make([]*SkipListNode, s.Height)
	current := s.Head

	for i := s.Level - 1; i >= 0; i-- {
		for current.Next[i] != nil && current.Next[i].Key < key {
			current = current.Next[i]
		}
		update[i] = current
	}

	current = current.Next[0]
	if current == nil || current.Key != key {
		return
	}

	for i := 0; i < s.Level; i++ {
		if update[i].Next[i] != current {
			break
		}
		update[i].Next[i] = current.Next[i]
	}

	s.Length--

	for s.Level > 1 && s.Head.Next[s.Level-1] == nil {
		s.Level--
	}
}

// GetLength returns the current number of elements in the skip list.
// This function is useful for monitoring the size of the skip list
// and can be used to determine when to resize or rehash the structure.
func (s *SkipList) GetLength() int {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Length
}

// GetLevel returns the current highest level in the skip list.
// This function is useful for understanding the structure of the skip list
// and can help in analyzing the performance of search, insertion,
// and deletion operations.
func (s *SkipList) GetLevel() int {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Level
}

// GetHeight returns the maximum height of the skip list.
// This function is useful for understanding the potential maximum
// number of levels in the skip list and can help in analyzing
// the performance of search, insertion, and deletion operations.
func (s *SkipList) GetHeight() int {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Height
}

// GetHead returns the head node of the skip list.
// This function is useful for accessing the sentinel node,
// which simplifies insertion and deletion operations.
func (s *SkipList) GetHead() *SkipListNode {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	return s.Head
}

// GetAllKeys returns a slice of all keys currently stored in the skip list.
// This function is useful for iterating over the keys in the skip list
// and can be used for various operations, such as exporting or
// displaying the contents of the skip list.
func (s *SkipList) GetAllKeys() []string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	keys := make([]string, 0, s.Length)
	current := s.Head.Next[0]
	for current != nil {
		keys = append(keys, current.Key)
		current = current.Next[0]
	}
	return keys
}

// GetAllValues returns a slice of all values currently stored in the skip list.
// This function is useful for iterating over the values in the skip list
// and can be used for various operations, such as exporting or
// displaying the contents of the skip list.
func (s *SkipList) GetAllValues() []string {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	values := make([]string, 0, s.Length)
	current := s.Head.Next[0]
	for current != nil {
		values = append(values, current.Value)
		current = current.Next[0]
	}
	return values
}

// Clear removes all elements from the skip list.
// This function is useful for resetting the skip list
// and can be used when the skip list is no longer needed.
func (s *SkipList) Clear() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Head.Next = make([]*SkipListNode, s.Height)
	s.Length = 0
	s.Level = 1
}

// Print prints the contents of the skip list.
// This function is useful for debugging and understanding
// the structure of the skip list. It displays the keys and values
// at each level, allowing for a visual representation of the skip list.
func (s *SkipList) Print() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for i := s.Level - 1; i >= 0; i-- {
		current := s.Head.Next[i]
		print("Level ", i, ": ")
		for current != nil {
			print(current.Key, " -> ")
			current = current.Next[i]
		}
		print("nil\n")
	}
}

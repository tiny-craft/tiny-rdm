package convutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type BitSetConvert struct{}

func (BitSetConvert) Enable() bool {
	return true
}

func (BitSetConvert) Encode(str string) (string, bool) {
	var result strings.Builder

	str = strings.ReplaceAll(str, "\r\n", "\n") // CRLF → LF
	str = strings.ReplaceAll(str, "\r", "\n")   // CR → LF

	lines := strings.Split(str, "\n")
	bytes := encodeToRedisBitset(lines)
	result.Write(bytes)

	return result.String(), true
}

func (BitSetConvert) Decode(str string) (string, bool) {
	bitset := getBitSet([]byte(str))

	var binBuilder strings.Builder
	for pos, value := range bitset {
		if value {
			if binBuilder.Len() > 0 {
				binBuilder.WriteByte('\n')
			}
			binBuilder.WriteString(fmt.Sprintf("%d", pos))
		}
	}

	return binBuilder.String(), true
}

// encodeToRedisBitset encodes a list of strings with integers (positions) into a Redis bitset byte array.
// The bit at position 'n' will be set to 1 if n is in the input array.
// The resulting byte slice can be stored in Redis using SET command.
func encodeToRedisBitset(numbers []string) []byte {
	if len(numbers) == 0 {
		return []byte{}
	}

	// Find the maximum number to determine the required bit length and convert strings to numbers
	maxNum := uint64(0)
	var validNumbers []uint64
	for _, s := range numbers {
		if s == "" {
			continue
		}
		num, err := strconv.ParseUint(s, 10, 64)
		if err != nil || num < 0 || num > math.MaxUint32 {
			fmt.Printf("Warning: skipping invalid number '%s': %v\n", s, err)
			continue
		}
		validNumbers = append(validNumbers, num)
		if num > maxNum {
			maxNum = num
		}
	}

	if len(validNumbers) == 0 {
		return []byte{}
	}

	// Calculate required byte length (8 bits per byte)
	byteLen := ((maxNum + 7) / 8) + 1

	// Initialize byte array
	bitset := make([]byte, byteLen)

	// Set bits for each number
	for _, num := range validNumbers {
		byteIndex := num / 8
		if byteIndex < byteLen {
			bitIndex := uint(num % 8)
			// Set the bit (big-endian bit order within byte)
			bitset[byteIndex] |= 1 << (7 - bitIndex)
		}
	}

	return bitset
}

func getBitSet(redisResponse []byte) []bool {
	bitset := make([]bool, len(redisResponse)*8)
	for i := range redisResponse {
		for j := 7; j >= 0; j-- {
			bitPos := uint(i*8 + (7 - j))
			bitset[bitPos] = (redisResponse[i] & (1 << uint(j))) > 0
		}
	}
	return bitset
}

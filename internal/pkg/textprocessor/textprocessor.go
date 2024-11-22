package textprocessor

import (
	"bytes"
	"strings"

	"github.com/jonathanhecl/chunker"
)

func CreateChunks(text string) []string {
	// Maximum characters per chunk
	const maxCharacters = 24576 // 8192 tokens * 3 (1 token = 3 caracters)

	var chunks []string
	var currentChunk strings.Builder

	// Use the chunker library to split text into sentences
	sentences := chunker.ChunkSentences(text)

	for _, sentence := range sentences {
		// Check if adding the sentence exceeds the character limit
		if currentChunk.Len()+len(sentence) <= maxCharacters {
			if currentChunk.Len() > 0 {
				currentChunk.WriteString(" ") // Add a space between sentences
			}
			currentChunk.WriteString(sentence)
		} else {
			// Add the completed chunk to the chunks slice
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()               // Start a new chunk
			currentChunk.WriteString(sentence) // Add the sentence to the new chunk
		}
	}

	// Add the last chunk if it has content
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	// Return the chunks
	return chunks
}

func ConcatenateStrings(strings []string) string {
	var result bytes.Buffer
	for _, str := range strings {
		result.WriteString(str)
	}
	return result.String()
}

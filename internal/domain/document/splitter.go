package document

import (
	"strings"
)

type TextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
}

func NewTextSplitter(chunkSize, chunkOverlap int) *TextSplitter {
	return &TextSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

func (s *TextSplitter) Split(content string) ([]string, error) {
	// 首先按段落分割
	paragraphs := strings.Split(content, "\n\n")

	var chunks []string
	var currentChunk strings.Builder
	currentLength := 0

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		words := strings.Fields(para)
		for _, word := range words {
			wordLength := len(word)

			// 如果添加这个词会超过chunk大小，并且当前chunk不为空
			if currentLength+wordLength > s.ChunkSize && currentLength > 0 {
				chunks = append(chunks, currentChunk.String())

				// 处理重叠部分
				if s.ChunkOverlap > 0 {
					lastChunk := currentChunk.String()
					overlapStart := len(lastChunk) - s.ChunkOverlap
					if overlapStart < 0 {
						overlapStart = 0
					}
					currentChunk.Reset()
					currentChunk.WriteString(lastChunk[overlapStart:])
					currentLength = len(lastChunk) - overlapStart
				} else {
					currentChunk.Reset()
					currentLength = 0
				}
			}

			if currentLength > 0 {
				currentChunk.WriteRune(' ')
				currentLength++
			}

			currentChunk.WriteString(word)
			currentLength += wordLength
		}

		// 段落结束后添加换行
		currentChunk.WriteString("\n\n")
		currentLength += 2
	}

	// 添加最后一个chunk
	if currentLength > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks, nil
}

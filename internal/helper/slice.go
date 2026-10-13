package helper

func Chunk[T any](items []T, size int) [][]T {
	chunks := make([][]T, 0, (len(items)+size-1)/size)

	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}

		chunks = append(chunks, items[i:end])
	}

	return chunks
}

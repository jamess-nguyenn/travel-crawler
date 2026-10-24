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

func Sample[T any](items []T, size int) []T {
	if size <= 0 {
		return []T{}
	}

	if size > len(items) {
		size = len(items)
	}

	return items[:size]
}

func MapSlice[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, 0, len(items))

	for _, item := range items {
		result = append(result, mapper(item))
	}

	return result
}

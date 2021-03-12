package views

func applyEscaping(s string) string {
	var (
		forbidden = ".!()-+#[]{}|"
		escaped   = ""
	)

LOOP:
	for _, char := range s {
		for _, forbChar := range forbidden {
			if char == forbChar {
				escaped += "\\" + string(char)
				continue LOOP
			}
		}
		escaped += string(char)
	}

	return escaped
}

func merge(x []string, y []string) (result []string) {
	copy(x, result)

	for _, elemY := range y {
		for _, elemX := range x {
			if elemX == elemY {
				continue
			}
		}
		result = append(result, elemY)
	}
	return result
}

func limit(slice []string, n int) []string {
	if n >= len(slice) {
		return slice
	}
	return slice[:n]
}

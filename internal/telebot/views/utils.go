package views

func applyEscaping(s string) string {
	var (
		forbidden = ".!()"
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

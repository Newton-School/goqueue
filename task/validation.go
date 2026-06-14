package task

const maxNameLength = 128

func validName(value string) bool {
	if value == "" || len(value) > maxNameLength {
		return false
	}

	for _, char := range value {
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		switch char {
		case '.', '_', ':', '-':
			continue
		default:
			return false
		}
	}

	return true
}

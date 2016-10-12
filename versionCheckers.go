package updater

func VCSimpleDiff(currentVer string) func(string) bool {
	return func(ver string) bool {
		if ver == "" {
			return false
		}

		return currentVer != ver
	}
}

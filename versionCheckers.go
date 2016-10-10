package updater

func VCSimpleDiff(currentVer string) func(string) bool {
	return func(ver string) bool {
		return currentVer != ver
	}
}

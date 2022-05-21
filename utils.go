package godexrcvr

func MmolToMg(mmol float32) int {
	return int(mmol * 18)
}

func MgToMmol(mg int) float32 {
	return float32(mg) / 18.0
}

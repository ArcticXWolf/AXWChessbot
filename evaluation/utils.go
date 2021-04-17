package evaluation

var (
	bitboardFileA uint64 = 0x0101010101010101
	bitboardFileH uint64 = 0x8080808080808080
)

func calculatePawnFileFill(pawnBitboard uint64) uint64 {
	pawnBitboard |= calculatePawnNorthFill(pawnBitboard)
	pawnBitboard |= calculatePawnSouthFill(pawnBitboard)
	return pawnBitboard
}

func calculatePawnNorthFill(pawnBitboard uint64) uint64 {
	pawnBitboard |= (pawnBitboard << 8)
	pawnBitboard |= (pawnBitboard << 16)
	pawnBitboard |= (pawnBitboard << 32)
	return pawnBitboard
}

func calculatePawnSouthFill(pawnBitboard uint64) uint64 {
	pawnBitboard |= (pawnBitboard >> 8)
	pawnBitboard |= (pawnBitboard >> 16)
	pawnBitboard |= (pawnBitboard >> 32)
	return pawnBitboard
}

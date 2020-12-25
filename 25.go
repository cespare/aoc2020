package main

func init() {
	addSolutions(25, problem25)
}

func problem25(ctx *problemContext) {
	scanner := ctx.scanner()
	scanner.scan()
	cardPubKey := scanner.int64()
	scanner.scan()
	doorPubKey := scanner.int64()
	ctx.reportLoad()

	ctx.reportPart1(transform(doorPubKey, findLoopSize(cardPubKey)))
}

func transform(subj, loopSize int64) int64 {
	n := int64(1)
	for i := int64(0); i < loopSize; i++ {
		n = (n * subj) % 20201227
	}
	return n
}

func findLoopSize(pubKey int64) int64 {
	n := int64(1)
	for loopSize := int64(0); ; loopSize++ {
		if n == pubKey {
			return loopSize
		}
		n = (n * 7) % 20201227
	}
}

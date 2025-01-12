package integration

import "strings"

func formatPairForBinance(pair string) string {
	return pairToLowerNoSlash(pair)
}

func formatPairForOKX(pair string) string {
	return pairToUpperNoSlash(pair)
}

func pairToLowerNoSlash(pair string) string {
	return strings.ToLower(strings.ReplaceAll(pair, "/", ""))
}

func pairToUpperNoSlash(pair string) string {
	return strings.ToUpper(strings.ReplaceAll(pair, "/", "-"))
}

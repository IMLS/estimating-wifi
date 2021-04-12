package wordlist

import "errors"

// golang really doesn't have `member`?
func contains(str string) bool {
	for _, s := range Wordlist {
		if str == s {
			return true
		}
	}
	return false
}

func GetPairIndex(str string) (int, error) {
	for ndx, s := range Wordlist {
		if str == s {
			return ndx, nil
		}
	}
	return -1, errors.New("wordpair not found")
}

func CheckWordpair(pair string) error {

	if contains(pair) {
		return nil
	} else {
		return errors.New("not found")
	}
}

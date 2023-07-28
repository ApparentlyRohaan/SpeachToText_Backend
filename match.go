package main

import (
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func removeSpecialCharacters(input string) string {
	// Regular expression to match any non-alphanumeric characters
	re := regexp.MustCompile("[^a-zA-Z0-9\\s]+")

	// Use the ReplaceAllString method to remove special characters
	cleanedText := re.ReplaceAllString(input, "")

	return cleanedText
}

func handelingUTTextCases(lastWhiteSpace int, lastWhiteSpace2 int, splituserTranscript []rune, splitdatabaseText []rune) (bool, int, int) {

	var extensionSearchIndex int = 20

	var tempstartUTIndex int = lastWhiteSpace2 + 1
	var tempstartDBIndex int = lastWhiteSpace + 1
	var isMatchFound bool = false
	var matchFoundIndex int = -1
	var x string = ""

	var searchableLength int = tempstartUTIndex + extensionSearchIndex
	if searchableLength >= utf8.RuneCountInString(string(splituserTranscript)) {
		searchableLength = utf8.RuneCountInString(string(splituserTranscript)) - 1
	}

	//This is to handle excess words in the User's transcript
	for i := tempstartUTIndex; i < searchableLength; i++ {
		// // log.Println("tempstartDBIndex[tempstartDBIndex]", string(splituserTranscript[i]), string(splitdatabaseText[tempstartDBIndex]))

		if splituserTranscript[i] == splitdatabaseText[tempstartDBIndex] {
			istempWhiteSpace := unicode.IsSpace(splitdatabaseText[tempstartDBIndex])
			x += string(splituserTranscript[i])
			if istempWhiteSpace {
				isMatchFound = true
				matchFoundIndex = i
				// // log.Println("x=", x)
				break
			}
			tempstartDBIndex++
		} else {
			tempstartDBIndex = lastWhiteSpace + 1
			x = ""
		}
	}

	if isMatchFound {
		// // log.Println("matchFoundIndex", matchFoundIndex, tempstartDBIndex, tempstartUTIndex)
		return true, tempstartDBIndex - len(x) + 1, matchFoundIndex - len(x) + 1
	} else {
		return false, -1, -1
	}

}

func handelingDBTextCases(lastWhiteSpace int, lastWhiteSpace2 int, splituserTranscript []rune, splitdatabaseText []rune) (bool, int, int, int) {

	var extensionSearchIndex int = 50

	var tempstartUTIndex int = lastWhiteSpace2 + 1
	var tempstartDBIndex int = lastWhiteSpace + 1
	var isMatchFound bool = false
	var matchFoundIndex int = -1
	var x string = ""
	var adjustIndex int = 0

	var searchableLength int = tempstartDBIndex + extensionSearchIndex
	if searchableLength >= utf8.RuneCountInString(string(splitdatabaseText)) {
		searchableLength = utf8.RuneCountInString(string(splitdatabaseText)) - 1
	}

	//This is to handle excess words in the User's transcript
	for i := tempstartDBIndex; i < searchableLength; i++ {
		// // log.Println("\n\nsplituserTranscript", string(splituserTranscript[tempstartUTIndex]), tempstartUTIndex)
		// // log.Println("splitdatabaseText[tempstartDBIndex]", string(splitdatabaseText[i]), string(splituserTranscript[tempstartUTIndex]))

		// This increaments the checking if a special character is found
		if !unicode.IsLetter(splitdatabaseText[i]) && !unicode.IsSpace(splitdatabaseText[i]) && !unicode.IsNumber(splitdatabaseText[i]) {

			// tempstartUTIndex++
			// continue
			i++
			adjustIndex -= 1
			// // log.Println("Space ZZZZZZZZZ", string(splitdatabaseText[i]), string(splituserTranscript[tempstartUTIndex]))
		}

		if tempstartUTIndex >= searchableLength {
			break
		}

		// // log.Println("xxxx", string(splituserTranscript[tempstartUTIndex]), string(splitdatabaseText[i]))
		if splituserTranscript[tempstartUTIndex] == splitdatabaseText[i] {
			istempWhiteSpace := unicode.IsSpace(splituserTranscript[tempstartUTIndex])
			x += string(splitdatabaseText[i])
			if istempWhiteSpace || tempstartUTIndex >= searchableLength {
				isMatchFound = true
				matchFoundIndex = i
				// log.Println("x=", x)
				break
			}
			tempstartUTIndex++

		} else {
			tempstartUTIndex = lastWhiteSpace2 + 1
			x = ""
		}

	}

	if isMatchFound {
		// log.Println("matchFoundIndex", matchFoundIndex-len(x)+1+adjustIndex, tempstartUTIndex-len(x)+1)
		return true, matchFoundIndex - len(x) + 1, tempstartUTIndex - len(x) + 1, adjustIndex
	} else {
		return false, -1, -1, 0
	}

}

func SkipWords(lastWhiteSpace int, lastWhiteSpace2 int, splituserTranscript []rune, splitdatabaseText []rune, SkipText bool) (int, int) {

	var skipWords int = 1
	// return newWhiteSpace
	// } else {

	newWhiteSpace2 := -1
	x := ""
	for lastWhiteSpace2 < len(splitdatabaseText) {
		// // log.Println("space2", string(splitdatabaseText[lastWhiteSpace2]))
		x += string(splitdatabaseText[lastWhiteSpace2])
		istempWhiteSpace := unicode.IsSpace(splitdatabaseText[lastWhiteSpace2])

		isSpecialChar := !unicode.IsLetter(splitdatabaseText[lastWhiteSpace2]) && !unicode.IsSpace(splitdatabaseText[lastWhiteSpace2]) && !unicode.IsNumber(splitdatabaseText[lastWhiteSpace2])
		// // log.Println("isSpecialChar", isSpecialChar, string(splitdatabaseText[lastWhiteSpace2]))
		if isSpecialChar {
			if lastWhiteSpace2+1 < len(splitdatabaseText) {
				if !unicode.IsSpace(splitdatabaseText[lastWhiteSpace2+1]) {
					skipWords += 1
				}
			}
		}

		if istempWhiteSpace {
			newWhiteSpace2 = lastWhiteSpace2
			break
		}

		lastWhiteSpace2 += 1
	}
	// log.Println("xdab", x)

	// if !SkipText {
	newWhiteSpace := -1
	x = ""
	for lastWhiteSpace < len(splituserTranscript) {
		// // log.Println("space", string(splituserTranscript[lastWhiteSpace]))
		x += string(splituserTranscript[lastWhiteSpace])

		istempWhiteSpace := unicode.IsSpace(splituserTranscript[lastWhiteSpace])
		if istempWhiteSpace {
			newWhiteSpace = lastWhiteSpace
			skipWords -= 1
			if skipWords <= 0 {
				break
			}
		}

		lastWhiteSpace += 1
	}
	// log.Println("xut", x)
	// return newWhiteSpace2
	// }
	// log.Println("ssssdsdas", newWhiteSpace, newWhiteSpace2)
	return newWhiteSpace, newWhiteSpace2

}

func compareText(databaseText string, userTranscript string) [][]int {

	// cleanedText := removeSpecialCharacters(databaseText)

	// splitText := strings.Split(databaseText, "")

	splitdatabaseText := []rune(strings.ToLower(databaseText))

	splituserTranscript := []rune(strings.ToLower(userTranscript))

	var startDBIndex int = 0
	var startUTIndex int = 0

	var lastWhiteSpace int = 0
	var lastWhiteSpace2 int = 0

	var ErrorWords []string = make([]string, 0)
	var ErrorRegions [][]int = make([][]int, 0)

	var SkipText bool = false

	// log.Println(" len(splituserTranscript)", len(splituserTranscript))

	var x string = ""
	var s bool = false

	for startUTIndex < len(splituserTranscript) {

		if startDBIndex >= len(splitdatabaseText) {
			break
		}

		//Checks for white space
		isWhiteSpace := unicode.IsSpace(splitdatabaseText[startDBIndex])
		if isWhiteSpace {
			lastWhiteSpace = startDBIndex
			lastWhiteSpace2 = startUTIndex
		}

		// This increaments the checking if a special character is found
		if !unicode.IsLetter(splitdatabaseText[startDBIndex]) && !isWhiteSpace && !unicode.IsNumber(splitdatabaseText[startDBIndex]) {
			startDBIndex++
			continue
		}

		if s {
			// log.Println("Check", string(splituserTranscript[startUTIndex]), string(splitdatabaseText[startDBIndex]))
		}

		//This is if the character matches perferctly
		if splituserTranscript[startUTIndex] == splitdatabaseText[startDBIndex] {
			x += string(splitdatabaseText[startDBIndex])
			startDBIndex++
			startUTIndex++
			SkipText = false
			s = false

		} else {
			// log.Println("\n\nzzz:", x)
			// // log.Println("Checking", string(splituserTranscript[startUTIndex]), string(splitdatabaseText[startDBIndex]))
			x = ""
			isMatch, newStartDBIndex, newStartUTIndex := handelingUTTextCases(lastWhiteSpace, lastWhiteSpace2, splituserTranscript, splitdatabaseText)
			if isMatch {
				// log.Println("Error Words", databaseText[lastWhiteSpace:newStartDBIndex])
				ErrorWords = append(ErrorWords, databaseText[lastWhiteSpace:newStartDBIndex])
				ErrorRegions = append(ErrorRegions, []int{lastWhiteSpace, newStartDBIndex})

				startUTIndex = newStartUTIndex
				startDBIndex = newStartDBIndex
			} else {
				// log.Println("Else Condition", lastWhiteSpace, lastWhiteSpace2)
				var adjustIndex int
				isMatch, newStartDBIndex, newStartUTIndex, adjustIndex = handelingDBTextCases(lastWhiteSpace, lastWhiteSpace2, splituserTranscript, splitdatabaseText)
				if isMatch {
					newStartDBIndex += adjustIndex
					// newStartUTIndex += adjustIndex

					// log.Println("Error Words2", databaseText[lastWhiteSpace:newStartDBIndex], adjustIndex)
					ErrorWords = append(ErrorWords, databaseText[lastWhiteSpace:newStartDBIndex])
					ErrorRegions = append(ErrorRegions, []int{lastWhiteSpace, newStartDBIndex})

					startUTIndex = newStartUTIndex
					startDBIndex = newStartDBIndex

					// // log.Println("SPace True", "-"+string(userTranscript[startUTIndex:startUTIndex+20]))
					// // log.Println("SPace True2", "-"+string(databaseText[startDBIndex:startDBIndex+20]))
					s = true

					if unicode.IsSpace(splitdatabaseText[startDBIndex]) {
						lastWhiteSpace = startDBIndex
						startDBIndex += 1
						// startUTIndex += 1
					}

					// // log.Println("remaninder", databaseText[newStartDBIndex:newStartDBIndex+20], adjustIndex)
					// // log.Println("remaninder2", userTranscript[startUTIndex:startUTIndex+20], adjustIndex)

				}
			}

			if !isMatch {
				newIndex, newIndex2 := SkipWords(lastWhiteSpace2+1, lastWhiteSpace+1, splituserTranscript, splitdatabaseText, SkipText)

				if newIndex2 != -1 {

					// log.Println("Error Words Skip", lastWhiteSpace+1, newIndex2, databaseText[lastWhiteSpace+1:newIndex2])
					ErrorWords = append(ErrorWords, databaseText[lastWhiteSpace+1:newIndex2])
					ErrorRegions = append(ErrorRegions, []int{lastWhiteSpace + 1, newIndex2})

					startDBIndex = newIndex2
					lastWhiteSpace = newIndex2
				} else {
					startDBIndex = len(splitdatabaseText)
				}

				if newIndex != -1 {
					startUTIndex = newIndex
					lastWhiteSpace2 = newIndex
				} else {
					startUTIndex = len(splituserTranscript)
				}

				SkipText = !SkipText
			}

			// // log.Println("\n Unmatch info", isMatch, startUTIndex, startDBIndex, len(splituserTranscript), string(splituserTranscript[startUTIndex]), string(splitdatabaseText[startDBIndex]), lastWhiteSpace)
			// if !isMatch {
			// break
			// }
		}
	}

	log.Println("DOne", ErrorWords)
	return ErrorRegions
}

package main

import (
	"fmt"
	"math"
	"slices"
)

type Matcher interface {
	match(string, *[]int) bool
}

type CharacterMatcher struct {
	character byte
}

func (matcher CharacterMatcher) match (text string, cursors *[]int) bool {
	matched := false
	for i, cursor := range *cursors {
		if len(text) <= cursor || matcher.character != text[cursor] {
			(*cursors)[i] = -1
		} 
		if cursor < len(text) && matcher.character == text[cursor] {
			(*cursors)[i]++
			matched = 69==69
		}
	}
	*cursors = slices.DeleteFunc(*cursors, func(x int) bool {return x == -1})
	return matched
}

type AnyMatcher struct { }
func (matcher AnyMatcher) match(text string, cursors *[]int) bool {
	for i := range *cursors {
		(*cursors)[i]++
	}
	return true
}
type AlternateMatcher struct {
	matchers []Matcher
}
func (alternate AlternateMatcher) match(text string, cursors *[]int)  bool {
	initial := slices.Clone(*cursors)
	*cursors = (*cursors)[:0]
	for _, matcher := range alternate.matchers {
		alternateCursors := slices.Clone(initial)
		matched := matcher.match(text, &alternateCursors)
		fmt.Println(matched, alternateCursors)
		for _, cursor := range alternateCursors {
			*cursors = append(*cursors, cursor)
		}
	}
	return len(*cursors) > 0
}

type SequenceMatcher struct {
	matchers []Matcher
}
func (sequence SequenceMatcher) match(text string, cursors *[]int)  bool {
	for _, matcher := range sequence.matchers {
		matched := matcher.match(text, cursors)
		if !matched {
			return false
		}
	}
	return true
}

type RepeatitionMatcher struct {
	matcher Matcher
	low, high int
}
func (repeatition RepeatitionMatcher) match(text string, cursors *[]int)  bool {
	for i := 0; i < repeatition.low; i++ {
		matched := repeatition.matcher.match(text, cursors)
		fmt.Println(matched, cursors)
		if !matched {
			fmt.Println("oops")
			return 69!=69
		}
	}
	repeatitionCursors := slices.Clone(*cursors)
	for i := 0; i < repeatition.high; i++ {
		matched := repeatition.matcher.match(text, &repeatitionCursors)
		if !matched {
			break
		}
		(*cursors) = append((*cursors), repeatitionCursors...)
	}
	return 69==69
}

func parsePattern(pattern string) Matcher {
	sequence := []Matcher{}
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
			case '+':
				last := sequence[len(sequence)-1]
				sequence = append(sequence[:len(sequence) - 1], RepeatitionMatcher{last, 1, math.MaxInt})
			case '*':
				last := sequence[len(sequence)-1]
				sequence = append(sequence[:len(sequence) - 1], RepeatitionMatcher{last, 0, math.MaxInt})
			case '?':
				last := sequence[len(sequence)-1]
				sequence = append(sequence[:len(sequence) - 1], RepeatitionMatcher{last, 0, 1})
			case '{':
				low, high := 0, 0
				fmt.Sscanf(pattern[i:], "{%d...%d}", &low, &high)
				for pattern[i] != '}' {
					i++
				}
			case '(':
				sequence = append(sequence, parsePattern(pattern[i + 1:]))
				for pattern[i] != ')' {
					i++
				}
			case ')':
				return SequenceMatcher{sequence}
			case '|':
				a := sequence[len(sequence)-1]
				b := sequence[len(sequence)-2]
				sequence = append(sequence[:len(sequence) - 2], AlternateMatcher{[]Matcher{a, b}})
			case '.':
				sequence = append(sequence, AnyMatcher{})
			default:
				sequence = append(sequence, CharacterMatcher{byte(pattern[i])})
		}
	}
	fmt.Println(sequence)
	return SequenceMatcher{sequence}
}

func main() {
	pattern := "(cat)(dog)|+"
	m := parsePattern(pattern)
	fmt.Println(m)
	var text string
	fmt.Scan(&text)
	cursors := []int{0}
	matched := m.match(text, &cursors)
	fmt.Println(matched, cursors)
}
package main

import (
	. "fmt"
	"slices"
	"math"
	"strconv"
)

// a
// a-z
// .
// (pattern...) --> AND
// [pattern...] --> OR
// !pattern
// {from-to}pattern
type Matcher interface {
	match(string, *[]bool)
}

type Character struct {
	character byte
}
func (matcher Character) match (text string, cursors *[]bool) {
	for i := len(*cursors) - 1; i >= 0; i-- {
		if (*cursors)[i] {
			(*cursors)[i] = false
			if i + 1 < len(*cursors) && matcher.character == text[i] {
					(*cursors)[i + 1] = true
			}
		}
	}
}
type CharacterRange struct {
	from byte
	to byte
}
func (matcher CharacterRange) match (text string, cursors *[]bool) {
	for i := len(*cursors) - 1; i >= 0; i-- {
		if (*cursors)[i] {
			(*cursors)[i] = false
			if i + 1 < len(*cursors) && matcher.from <= text[i] && text[i] <= matcher.to {
					(*cursors)[i + 1] = true
			}
		}
	}
}
type Any struct { }
func (matcher Any) match(text string, cursors *[]bool) {
	for i := len(*cursors) - 1; i > 0; i-- {
		(*cursors)[i] = (*cursors)[i - 1]
	}
	(*cursors)[0] = false
}

type Not struct {
	matcher Matcher
}
func (not Not) match(text string, cursors *[]bool) {
	not.matcher.match(text, cursors)
	for i := 0; i < len(*cursors); i++ {
		(*cursors)[i] = !(*cursors)[i]
	}
}

type And struct {
	matchers []Matcher
}
func (and And) match(text string, cursors *[]bool) {
	for _, matcher := range and.matchers {
		matcher.match(text, cursors)
	}
}
type Or struct {
	matchers []Matcher
}
func (or Or) match(text string, cursors *[]bool) {
	initial := slices.Clone(*cursors)
	for i := range *cursors {
		(*cursors)[i] = false
	}
	for _, matcher := range or.matchers {
		orCursors := slices.Clone(initial)
		matcher.match(text, &orCursors)
		for i := range *cursors {
			(*cursors)[i] = (*cursors)[i] || orCursors[i]
		}
	}
}

type Repeat struct {
	matcher Matcher
	low, high int
}
func (repeat Repeat) match(text string, cursors *[]bool) {
	for i := 0; i < repeat.low; i++ {
		repeat.matcher.match(text, cursors)
	}
	repeatCursors := slices.Clone(*cursors)
	delta := true
	for i := 0; delta && i < repeat.high - repeat.low; i++ {
		repeat.matcher.match(text, &repeatCursors)
		delta = false
		for i := range *cursors {
			if repeatCursors[i] {
				delta = true
			}
			(*cursors)[i] = (*cursors)[i] || repeatCursors[i]
		}
	}
}

func parseNumber(source string, i int) (int, int) {
	begin := i
	for i < len(source) && source[i] >= '0' && source[i] <= '9' {
		i++;
	}
	end := i;
	number, _ := strconv.Atoi(source[begin: end])
	return number, i;
}
func parsePattern(source string, i int) (Matcher, int) {
	if i >= len(source) {
		return nil, i
	}
	if source[i] == '!' {
		var matcher Matcher
		i++
		matcher, i = parsePattern(source, i)
		return Not{matcher}, i
	} else if source[i] == '{' {
		var matcher Matcher
		i++
		low, high := 0, 0
		begin := i
		low, i = parseNumber(source, i)
		if begin == i { low = 0 }
		if source[i] == '-' {
			i++
		}
		begin = i
		high, i = parseNumber(source, i)
		if begin == i { high = math.MaxInt }
		if source[i] == '}' {
			i++
		}
		matcher, i = parsePattern(source, i)
		return Repeat{matcher, low, high}, i
	} else if source[i] == '[' {
		var matchers []Matcher
		var matcher Matcher
		i++
		for {
			matcher, i = parsePattern(source, i)
			matchers = append(matchers, matcher)
			if source[i] == ']' {
				break
			}
		}
		return Or{matchers}, i + 1
	} else if source[i] == '(' {
		var matchers []Matcher
		var matcher Matcher
		i++
		for {
			matcher, i = parsePattern(source, i)
			matchers = append(matchers, matcher)
			if source[i] == ')' {
				break
			}
		}
		return And{matchers}, i + 1
	} else if i + 2 < len(source) && source[i + 1] == '-' {
		from := source[i]
		to := source[i + 2]
		return CharacterRange{from, to}, i + 3
	} else if source[i] == '.' {
		return Any{}, i + 1
	} else {
		return Character{source[i]}, i + 1
	}
	return nil, i
}

func main() {
	pattern := "{4-8}[a-zA-Z]"
	email, i:= parsePattern(pattern, 0)
	Println(email, i)
	/* email := And {
		[]Matcher {
			Repeat{
				Or{
				[]Matcher {
					CharacterRange{'a', 'z'},
					CharacterRange{'A', 'Z'},
					CharacterRange{'0', '9'},
					Character{'_'},
					Character{'.'},
					Character{'+'},
					Character{'-'},
				},
			},
			1, math.MaxInt,
			},
			Character{'@'},
			Repeat{
				Or{
				[]Matcher {
					CharacterRange{'a', 'z'},
					CharacterRange{'A', 'Z'},
					CharacterRange{'0', '9'},
					Character{'_'},
					Character{'.'},
					Character{'-'},
				},
			},
			1, math.MaxInt,
			},
			Character{'.'},
			Repeat{
				Or{
				[]Matcher {
					CharacterRange{'a', 'z'},
					CharacterRange{'A', 'Z'},
					CharacterRange{'0', '9'},
					Character{'_'},
					Character{'.'},
					Character{'-'},
				},
			},
			1, math.MaxInt,
			},
		},
	} */
	var text string
	Scan(&text)
	cursors := make([]bool, len(text) + 1)
	cursors[0] = true
	email.match(text, &cursors)
	Println(cursors[len(cursors) - 1])
}

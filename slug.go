// Copyright 2013 by Dobrosław Żybort. All rights reserved.
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package slug

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/gosimple/unidecode"
)

var (
	// CustomSub stores custom substitution map
	CustomSub map[string]string
	// CustomRuneSub stores custom rune substitution map
	CustomRuneSub map[rune]string

	// MaxLength stores maximum slug length.
	// It's smart so it will cat slug after full word.
	// By default slugs aren't shortened.
	// If MaxLength is smaller than length of the first word, then returned
	// slug will contain only substring from the first word truncated
	// after MaxLength.
	MaxLength int

	// Lowercase defines if the resulting slug is transformed to lowercase.
	// Default is true.
	Lowercase = true

	regexpNonAuthorizedChars = regexp.MustCompile("[^a-zA-Z0-9-_]")
	regexpMultipleDashes     = regexp.MustCompile("-+")
)

//=============================================================================

// Make returns slug generated from provided string. Will use "en" as language
// substitution.
func Make(s string) (slug string) {
	return MakeLang(s, "en")
}

// MakeLang returns slug generated from provided string and will use provided
// language for chars substitution.
func MakeLang(s string, lang string) (slug string) {
	slug = strings.TrimSpace(s)

	// Custom substitutions
	// Always substitute runes first
	slug = SubstituteRune(slug, CustomRuneSub)
	slug = Substitute(slug, CustomSub)

	// Process string with selected substitution language.
	// Catch ISO 3166-1, ISO 639-1:2002 and ISO 639-3:2007.
	switch strings.ToLower(lang) {
	case "de", "deu":
		slug = SubstituteRune(slug, deSub)
	case "en", "eng":
		slug = SubstituteRune(slug, enSub)
	case "es", "spa":
		slug = SubstituteRune(slug, esSub)
	case "fi", "fin":
		slug = SubstituteRune(slug, fiSub)
	case "fr", "fra":
		slug = SubstituteRune(slug, frSub)
	case "gr", "el", "ell":
		slug = SubstituteRune(slug, grSub)
	case "kz", "kk", "kaz":
		slug = SubstituteRune(slug, kkSub)
	case "nl", "nld":
		slug = SubstituteRune(slug, nlSub)
	case "pl", "pol":
		slug = SubstituteRune(slug, plSub)
	case "sv", "swe":
		slug = SubstituteRune(slug, svSub)
	case "nn", "nno":
		slug = SubstituteRune(slug, nnSub)
	case "nb", "nob":
		slug = SubstituteRune(slug, nbSub)
	case "sl", "slv":
		slug = SubstituteRune(slug, slSub)
	case "tr", "tur":
		slug = SubstituteRune(slug, trSub)
	default: // fallback to "en" if lang not found
		slug = SubstituteRune(slug, enSub)
	}

	// Process all non ASCII symbols
	slug = unidecode.Unidecode(slug)

	if Lowercase {
		slug = strings.ToLower(slug)
	}

	// Process all remaining symbols
	slug = regexpNonAuthorizedChars.ReplaceAllString(slug, "-")
	slug = regexpMultipleDashes.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-_")

	if MaxLength > 0 {
		slug = smartTruncate(slug)
	}

	return slug
}

// Substitute returns string with superseded all substrings from
// provided substitution map. Substitution map will be applied in alphabetic
// order. Many passes, on one substitution another one could apply.
func Substitute(s string, sub map[string]string) (buf string) {
	buf = s
	var keys []string
	for k := range sub {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		buf = strings.Replace(buf, key, sub[key], -1)
	}
	return
}

// SubstituteRune substitutes string chars with provided rune
// substitution map. One pass.
func SubstituteRune(s string, sub map[rune]string) string {
	var buf bytes.Buffer
	for _, c := range s {
		if d, ok := sub[c]; ok {
			buf.WriteString(d)
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

func smartTruncate(text string) string {
	if len(text) < MaxLength {
		return text
	}

	var truncated string
	words := strings.SplitAfter(text, "-")
	// If MaxLength is smaller than length of the first word return word
	// truncated after MaxLength.
	if len(words[0]) > MaxLength {
		return words[0][:MaxLength]
	}
	for _, word := range words {
		if len(truncated)+len(word)-1 <= MaxLength {
			truncated = truncated + word
		} else {
			break
		}
	}
	return strings.Trim(truncated, "-")
}

// IsSlug returns True if provided text does not contain white characters,
// punctuation, all letters are lower case and only from ASCII range.
// It could contain `-` and `_` but not at the beginning or end of the text.
// It should be in range of the MaxLength var if specified.
// All output from slug.Make(text) should pass this test.
func IsSlug(text string) bool {
	if text == "" ||
		(MaxLength > 0 && len(text) > MaxLength) ||
		text[0] == '-' || text[0] == '_' ||
		text[len(text)-1] == '-' || text[len(text)-1] == '_' {
		return false
	}
	for _, c := range text {
		if (c < 'a' || c > 'z') && c != '-' && c != '_' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

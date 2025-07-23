package main

import (
	"encoding/json"
	"fmt"
)

func parseWordJSON(inputWord string, rawData []byte) (*Word, error) {

	var raw []interface{}
	err := json.Unmarshal(rawData, &raw)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON - Structure not matched: %v", err)
	}

	if len(raw) == 9 {
		return nil, fmt.Errorf("No Definition for %v", inputWord)
	}

	word := &Word{
		Word:     inputWord,
		SynGroup: []SynonymGroup{},
		DefGroup: []DefinitionGroup{},
		Examples: []Example{},
	}

	// parse phonetic @[0][1][3]
	parseIPA(word, raw[0])

	// parse definitions @[12]
	if len(raw) > 12 {
		parseDefGroup(word, raw[12])
	}

	// parse synonyms @[11]
	if len(raw) > 11 && raw[11] != nil {
		parseSynGroup(word, raw[11])
	}

	// parse examples @[13]
	if len(raw) > 13 && raw[13] != nil {
		parseExamples(word, raw[13])
	}

	// at [14] base word is represented as [ [base-word ]] but who gives a fuck
	return word, nil
}

func parseIPA(word *Word, raw interface{}) {

	if rawList, ok := raw.([]interface{}); ok {
		if len(rawList) > 1 {
			secondList, ok := rawList[1].([]interface{})
			if ok && len(secondList) > 0 {
				if ipa, ok := secondList[len(secondList)-1].(string); ok {
					word.IPA = "/" + ipa + "/"
				}
			}
		}
	}
}

func parseDefGroup(word *Word, raw interface{}) {
	if rawGroups, ok := raw.([]interface{}); ok {
		for _, groups := range rawGroups {
			if rawGroup, ok := groups.([]interface{}); ok {

				defGroup := &DefinitionGroup{}
				defGroup.POS = rawGroup[0].(string)
				parseDefinitions(defGroup, rawGroup[1])
				word.DefGroup = append(word.DefGroup, *defGroup)
			}
		}
	}
}

func parseDefinitions(defGroup *DefinitionGroup, raw interface{}) {
	if rawDefs, ok := raw.([]interface{}); ok {
		for _, rawDef := range rawDefs {
			if def, ok := rawDef.([]interface{}); ok {

				d := &Definition{}

				// definition and reference always exist(based on observation)
				d.Meaning = def[0].(string)
				d.Ref = def[1].(string)

				if len(def) >= 3 && def[2] != nil {
					d.Example = def[2].(string)
				}
				// process if any extra contextual category/domain exist to each definition
				if len(def) == 4 {
					parseDomainAndCategoryInDef(d, def[3])
				}
				defGroup.Definition = append(defGroup.Definition, *d)
			}
		}
	}
}

func parseDomainAndCategoryInDef(d *Definition, raw interface{}) {
	unknown := raw.([]interface{})
	if len(unknown) == 1 {
		for _, c := range unknown[0].([]interface{}) {
			d.Category = append(d.Category, c.(string))
		}
	} else if len(unknown) == 3 {
		for _, c := range unknown[2].([]interface{}) {
			d.Domains = append(d.Domains, c.(string))
		}
	} else {
		panic("Unknown Number of word Category! Nor linguistic or domain specific")
	}
}

func parseSynGroup(word *Word, raw interface{}) {
	rawGroups := raw.([]interface{})
	for _, rawGroups := range rawGroups {

		group := rawGroups.([]interface{})
		synGroup := &SynonymGroup{}
		synGroup.POS = group[0].(string)
		parseSynonyms(synGroup, group[1])
		word.SynGroup = append(word.SynGroup, *synGroup)

	}
}
func parseSynonyms(synGroup *SynonymGroup, raw interface{}) {
	rawSyns := raw.([]interface{})
	for _, rawSyn := range rawSyns {

		syns := rawSyn.([]interface{})
		syn := &Synonyms{}
		syn.Ref = syns[1].(string)

		for _, s := range syns[0].([]interface{}) {
			syn.Synonyms = append(syn.Synonyms, s.(string))
		}

		if len(syns) == 3 {
			parseDomainAndCategoryInSyn(syn, syns[2])
		}
		synGroup.Synonyms = append(synGroup.Synonyms, *syn)
	}

}

func parseDomainAndCategoryInSyn(s *Synonyms, raw interface{}) {
	unknown := raw.([]interface{})

	if len(unknown) == 1 {
		for _, c := range unknown[0].([]interface{}) {
			s.Category = append(s.Category, c.(string))
		}
	} else if len(unknown) == 3 {
		for _, c := range unknown[2].([]interface{}) {
			s.Domains = append(s.Domains, c.(string))
		}
	} else {
		panic("Unknown Number of word Category! Nor linguistic or domain specific")
	}
}

func parseExamples(word *Word, raw interface{}) {

	outer := raw.([]interface{})
	inner := outer[0].([]interface{})
	for _, example := range inner {
		ex := example.([]interface{})

		example := &Example{}
		example.Text = ex[0].(string)
		example.Ref = ex[5].(string)

		word.Examples = append(word.Examples, *example)
	}
}

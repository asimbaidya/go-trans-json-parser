package main

type Word struct {
	Word     string
	IPA      string
	SynGroup []SynonymGroup
	DefGroup []DefinitionGroup
	Examples []Example
}

type SynonymGroup struct {
	POS      string // index 0
	Tags     []string
	Synonyms []Synonyms
}

type Synonyms struct {
	Ref      string   `json:"-"` // index 1,1
	Synonyms []string // index 1,0
	Domains  []string
	Category []string
}

type DefinitionGroup struct {
	POS        string       // index 0
	Definition []Definition // index 1,0
}

type Definition struct {
	Ref     string `json:"-"` // index 1,1
	Meaning string
	Example string

	Domains  []string
	Category []string
}

type Example struct {
	Text string // index 0,0
	Ref  string `json:"-"` // index 0,5
}

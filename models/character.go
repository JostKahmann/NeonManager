package models

import (
	"html/template"
	"strconv"
)

type ChoiceType int

const (
	ChoiceAffinity ChoiceType = iota
	ChoiceAbility
	ChoiceSkill
)

type ModType int

const (
	ModAnyStat ModType = iota
	ModStat
)

type DbItem interface {
	Pk() string
	SetName(string)
}

type LevelItem interface {
	Pk() string
	SetName(string)
	GetLevel() int
	SetLevel(int)
}

type Character struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	GP          int         `json:"gp"`
	XP          int         `json:"xp"`
	Stats       Stats       `json:"stats"`
	Race        Race        `json:"race"`
	Background  Background  `json:"background"`
	Boons       []*Affinity `json:"boons"`
	Banes       []*Affinity `json:"banes"`
	Abilities   []*Ability  `json:"abilities"`
	Skills      []*Skill    `json:"skills"`
	Description string      `json:"description"`
	Extensions  []string    `json:"extensions"`
}

type Stats struct {
	Id  int `json:"id"`
	Cr  int `json:"cr"`
	Int int `json:"int"`
	Ins int `json:"ins"`
	Ch  int `json:"ch"`
	Ag  int `json:"ag"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Str int `json:"str"`
}

type Race struct {
	Name           string        `json:"name"`
	Cost           int           `json:"cost"`
	Extension      string        `json:"extension"`
	Description    string        `json:"description"`
	Stats          Stats         `json:"stats"`
	Backgrounds    []*Background `json:"backgrounds"`    // usual backgrounds
	NotBackgrounds []*Background `json:"notBackgrounds"` // disallowed backgrounds
	Boons          []*Affinity   `json:"boons"`
	Banes          []*Affinity   `json:"banes"`
	Abilities      []*Ability    `json:"abilities"`
	Skills         []*Skill      `json:"skills"`
	Choices        []ChoiceGroup `json:"choices"`
}

type Background struct {
	Name        string        `json:"name"`
	Cost        int           `json:"cost"`
	Extension   string        `json:"extension"`
	Description string        `json:"description"`
	Stats       Stats         `json:"stats"`
	Boons       []*Affinity   `json:"boons"`
	Banes       []*Affinity   `json:"banes"`
	Abilities   []*Ability    `json:"abilities"`
	Skills      []*Skill      `json:"skills"`
	Choices     []ChoiceGroup `json:"choices"`
}

type ChoiceGroup struct {
	ChoiceType ChoiceType  `json:"type"`
	Count      int         `json:"count"`
	Affinities []*Affinity `json:"affinities"`
	Abilities  []*Ability  `json:"abilities"`
	Skills     []*Skill    `json:"skills"`
}

type Affinity struct {
	Name        string        `json:"name"`
	Cost        int           `json:"cost"`
	Extension   string        `json:"extension"`
	Description string        `json:"description"`
	IsBoon      bool          `json:"isBoon"`
	Level       int           `json:"level"`
	Modifiers   []AffinityMod `json:"modifiers"`
}

type AffinityMod struct {
	ModType ModType `json:"type"`
	Stat    string  `json:"stat"`
	StatMod int     `json:"statMod"`
}

type Ability struct {
	Name      string       `json:"name"`
	Cost      int          `json:"cost"`
	Extension string       `json:"extension"`
	Effect    string       `json:"effect"`
	Requires  []AbilityReq `json:"requires"`
}

type AbilityReq struct {
	Id           int        `json:"id"`
	Ability      string     `json:"ability"`
	RequiredType ChoiceType `json:"type"`
	RAffinity    string     `json:"rAffinity"`
	RAbility     string     `json:"rAbility"`
	RSkill       string     `json:"rSkill"`
}

type Skill struct {
	Name        string `json:"name"`
	Cost        int    `json:"cost"`
	Extension   string `json:"extension"`
	Description string `json:"description"`
	Stat        string `json:"stat"`
	Level       int    `json:"level"`
}

type Article struct {
	Id    int           `json:"id"`
	Title string        `json:"title"`
	Text  template.HTML `json:"text"`
	Table string        `json:"table"`
	Tags  []string      `json:"tags"`
}

func (c Character) Pk() string {
	return strconv.Itoa(c.Id)
}

func (c Character) SetName(name string) {
	c.Name = name
}

func (s Stats) Pk() string {
	return strconv.Itoa(s.Id)
}

func (s Stats) SetName(_ string) {

}

func (r Race) Pk() string {
	return r.Name
}

func (r Race) SetName(name string) {
	r.Name = name
}

func (b Background) Pk() string {
	return b.Name
}

func (b Background) SetName(name string) {
	b.Name = name
}

func (a Affinity) Pk() string {
	return a.Name
}

func (a Affinity) SetName(name string) {
	a.Name = name
}

func (a Affinity) GetLevel() int {
	return a.Level
}

func (a Affinity) SetLevel(level int) {
	a.Level = level
}

func (a Ability) Pk() string {
	return a.Name
}

func (a Ability) SetName(name string) {
	a.Name = name
}

func (a Ability) GetLevel() int {
	return 0
}

func (a Ability) SetLevel(_ int) {

}

func (a AbilityReq) Pk() string {
	return strconv.Itoa(a.Id)
}

func (a AbilityReq) SetName(_ string) {

}

func (a AbilityReq) SetLevel(_ int) {

}

func (s Skill) Pk() string {
	return s.Name
}

func (s Skill) SetName(name string) {
	s.Name = name
}

func (s Skill) GetLevel() int {
	return s.Level
}

func (s Skill) SetLevel(level int) {
	s.Level = level
}

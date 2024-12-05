package data

import (
	"NeonManager/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func parseCsv[T any](filePath string, parse func([]string) T) ([]T, error) {
	var file *os.File
	if f, err := os.Open(filePath); err != nil {
		return nil, err
	} else {
		file = f
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'
	items := make([]T, 0)
	for record, err := reader.Read(); err == nil && record != nil; record, err = reader.Read() {
		items = append(items, parse(record))
	}
	return items, nil
}

func writeCsv[T any](filePath string, stringify func(T) []string, items []T) error {
	var file *os.File
	if f, err := os.Open(filePath); err != nil {
		return err
	} else {
		file = f
	}
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	for _, item := range items {
		itemStr := stringify(item)
		err := writer.Write(itemStr)
		if err != nil {
			return fmt.Errorf("error writing item %v csv: %w", itemStr, err)
		}
	}
	return nil
}

func parseSkill(line []string) models.Skill {
	skill := models.Skill{}
	for i, entry := range line {
		if entry == "" {
			continue
		}
		switch i {
		case 0:
			skill.Name = entry
		case 1:
			if cost, err := strconv.Atoi(entry); err == nil {
				skill.Cost = cost
			}
		case 2:
			skill.Extension = entry
		case 3:
			skill.Description = entry
		case 4:
			skill.Stat = entry
		}
	}
	return skill
}

func parseAbility(line []string) models.Ability {
	ability := models.Ability{}
	for i, entry := range line {
		if entry == "" {
			continue
		}
		switch i {
		case 0:
			ability.Name = entry
		case 1:
			if cost, err := strconv.Atoi(entry); err == nil {
				ability.Cost = cost
			}
		case 2:
			ability.Extension = entry
		case 3:
			ability.Effect = entry
		case 4:
			reqs := strings.Split(entry, ",")
			ability.Requires = make([]models.AbilityReq, len(reqs))
			for j, req := range reqs {
				reqSpreads := strings.Split(req, "~")
				var choice models.ChoiceType

				if c, err := strconv.Atoi(reqSpreads[0]); err != nil || c < 0 || c > 2 {
					choice = models.ChoiceAffinity
				} else {
					choice = models.ChoiceType(c)
				}
				ability.Requires[j].RequiredType = choice
				switch choice {
				case models.ChoiceAffinity:
					ability.Requires[j].RAffinity = reqSpreads[1]
				case models.ChoiceAbility:
					ability.Requires[j].RAbility = reqSpreads[1]
				case models.ChoiceSkill:
					ability.Requires[j].RSkill = reqSpreads[1]
				}
			}
		}
	}
	return ability
}

func parseAffinity(line []string) models.Affinity {
	affinity := models.Affinity{}
	for i, entry := range line {
		if entry == "" {
			continue
		}
		switch i {
		case 0:
			affinity.Name = entry
		case 1:
			if cost, err := strconv.Atoi(entry); err == nil {
				affinity.Cost = cost
			}
		case 2:
			affinity.Extension = entry
		case 3:
			affinity.Description = entry
		case 4:
			affinity.IsBoon = entry == "1"
		case 5:
			_ = strings.Split(entry, ",")
			// todo
		}
	}
	return affinity
}

func parseStats(line []string) models.Stats {
	statVals := make([]int, 8)
	for j, statStr := range line {
		if stat, err := strconv.Atoi(statStr); err == nil {
			statVals[j] = stat
		}
	}
	return models.Stats{
		Cr:  statVals[0],
		Int: statVals[1],
		Ins: statVals[2],
		Ch:  statVals[3],
		Ag:  statVals[4],
		Dex: statVals[5],
		Con: statVals[6],
		Str: statVals[7],
	}
}

func parseLevelList[T models.LevelItem](line []string, _ T) []*T {
	items := make([]*T, len(line))
	for i, entry := range line {
		items[i] = new(T)
		split := strings.Split(entry, "@")
		(*items[i]).SetName(split[0])
		if len(split) == 2 {
			if level, err := strconv.Atoi(split[1]); err == nil {
				(*items[i]).SetLevel(level)
			}
		}
	}
	return items
}

func parseList[T models.DbItem](line []string, _ T) []*T {
	items := make([]*T, len(line))
	for i, entry := range line {
		items[i] = new(T)
		(*items[i]).SetName(entry)
	}
	return items
}

func parseChoices(line []string) []models.ChoiceGroup {
	choices := make([]models.ChoiceGroup, len(line))
	for i, entry := range line {
		if entry == "" {
			continue
		}
		items := strings.Split(entry, "~")
		var choiceType models.ChoiceType
		if c, err := strconv.Atoi(items[0]); err != nil || c < 0 || c > 2 {
			choiceType = models.ChoiceAffinity
		} else {
			choiceType = models.ChoiceType(c)
		}
		var count int
		if c, err := strconv.Atoi(items[1]); err != nil {
			count = 0
		} else {
			count = c
		}
		group := models.ChoiceGroup{ChoiceType: choiceType, Count: count}
		switch choiceType {
		case models.ChoiceAffinity:
			group.Affinities = parseLevelList(items[2:], models.Affinity{})
		case models.ChoiceAbility:
			group.Abilities = parseList(items[2:], models.Ability{})
		case models.ChoiceSkill:
			group.Skills = parseLevelList(items[2:], models.Skill{})
		}
		choices[i] = group
	}
	return choices
}

func parseBackground(line []string) models.Background {
	background := models.Background{}
	for i, entry := range line {
		if entry == "" {
			continue
		}
		switch i {
		case 0:
			background.Name = entry
		case 1:
			if cost, err := strconv.Atoi(entry); err == nil {
				background.Cost = cost
			}
		case 2:
			background.Extension = entry
		case 3:
			background.Description = entry
		case 4:
			background.Stats = parseStats(strings.Split(entry, ","))
		case 5:
			background.Boons = parseLevelList(strings.Split(entry, ","), models.Affinity{})
		case 6:
			background.Banes = parseLevelList(strings.Split(entry, ","), models.Affinity{})
		case 7:
			background.Abilities = parseList(strings.Split(entry, ","), models.Ability{})
		case 8:
			background.Skills = parseLevelList(strings.Split(entry, ","), models.Skill{})
		case 9:
			background.Choices = parseChoices(strings.Split(entry, ","))
		}
	}
	return background
}

func parseRace(line []string) models.Race {
	race := models.Race{}
	for i, entry := range line {
		if entry == "" {
			continue
		}
		switch i {
		case 0:
			race.Name = entry
		case 1:
			if cost, err := strconv.Atoi(entry); err == nil {
				race.Cost = cost
			}
		case 2:
			race.Extension = entry
		case 3:
			race.Description = entry
		case 4:
			race.Stats = parseStats(strings.Split(entry, ","))
		case 5:
			race.Backgrounds = parseList(strings.Split(entry, ","), models.Background{})
		case 6:
			race.NotBackgrounds = parseList(strings.Split(entry, ","), models.Background{})
		case 7:
			race.Boons = parseLevelList(strings.Split(entry, ","), models.Affinity{})
		case 8:
			race.Banes = parseLevelList(strings.Split(entry, ","), models.Affinity{})
		case 9:
			race.Abilities = parseList(strings.Split(entry, ","), models.Ability{})
		case 10:
			race.Skills = parseLevelList(strings.Split(entry, ","), models.Skill{})
		case 11:
			race.Choices = parseChoices(strings.Split(entry, ","))
		}
	}
	return race
}

func flatten(values []string, sep rune) string {
	var builder strings.Builder
	for i, value := range values {
		if i != 0 {
			builder.WriteRune(sep)
		}
		builder.WriteString(value)
	}
	return builder.String()
}

func flattenFunc[T any](values []T, toVal func(T) string, sep rune) string {
	var builder strings.Builder
	for i, value := range values {
		if i != 0 {
			builder.WriteRune(sep)
		}
		builder.WriteString(toVal(value))
	}
	return builder.String()
}

func writeSkill(skill models.Skill) []string {
	return []string{
		skill.Name,
		strconv.Itoa(skill.Cost),
		skill.Extension,
		skill.Description,
		skill.Stat}
}

func writeAbility(ability models.Ability) []string {
	reqs := make([]string, len(ability.Requires))
	for i, req := range ability.Requires {
		switch req.RequiredType {
		case models.ChoiceAffinity:
			reqs[i] = strconv.Itoa(int(req.RequiredType)) + "~" + req.RAffinity
		case models.ChoiceAbility:
			reqs[i] = strconv.Itoa(int(req.RequiredType)) + "~" + req.RAbility
		case models.ChoiceSkill:
			reqs[i] = strconv.Itoa(int(req.RequiredType)) + "~" + req.RSkill
		}
	}
	return []string{
		ability.Name,
		strconv.Itoa(ability.Cost),
		ability.Extension,
		ability.Effect,
		flatten(reqs, ',')}
}

func writeAffinity(affinity models.Affinity) []string {
	var isBoon string
	if affinity.IsBoon {
		isBoon = "1"
	} else {
		isBoon = "0"
	}
	return []string{
		affinity.Name,
		strconv.Itoa(affinity.Cost),
		affinity.Extension,
		affinity.Description,
		isBoon,
	}
}

func writeStats(stats models.Stats, sep rune) string {
	return flattenFunc([]int{
		stats.Cr,
		stats.Int,
		stats.Ins,
		stats.Ch,
		stats.Ag,
		stats.Dex,
		stats.Con,
		stats.Str,
	}, strconv.Itoa, sep)
}

func wAff(item *models.Affinity) string {
	return (*item).Name + "@" + strconv.Itoa((*item).Level)
}

func wAbl(item *models.Ability) string {
	return (*item).Name
}

func wSkl(item *models.Skill) string {
	return (*item).Name + "@" + strconv.Itoa((*item).Level)
}

func writeChoices(choices []models.ChoiceGroup) string {
	items := make([]string, len(choices))
	for i, item := range choices {
		items[i] = strconv.Itoa(int(item.ChoiceType)) + "~" + strconv.Itoa(item.Count) + "~"
		switch item.ChoiceType {
		case models.ChoiceAffinity:
			items[i] += flattenFunc(item.Affinities, wAff, '~')
		case models.ChoiceAbility:
			items[i] += flattenFunc(item.Abilities, wAbl, '~')
		case models.ChoiceSkill:
			items[i] += flattenFunc(item.Skills, wSkl, '~')
		}
	}
	return flatten(items, ',')
}

func writeBackground(background models.Background) []string {
	return []string{
		background.Name,
		strconv.Itoa(background.Cost),
		background.Extension,
		background.Description,
		writeStats(background.Stats, ','),
		flattenFunc(background.Boons, wAff, ','),
		flattenFunc(background.Banes, wAff, ','),
		flattenFunc(background.Abilities, wAbl, ','),
		flattenFunc(background.Skills, wSkl, ','),
		writeChoices(background.Choices),
	}
}

func wBg(item *models.Background) string {
	return (*item).Name
}

func writeRace(race models.Race) []string {
	return []string{
		race.Name,
		strconv.Itoa(race.Cost),
		race.Extension,
		race.Description,
		writeStats(race.Stats, ','),
		flattenFunc(race.Backgrounds, wBg, ','),
		flattenFunc(race.NotBackgrounds, wBg, ','),
		flattenFunc(race.Boons, wAff, ','),
		flattenFunc(race.Banes, wAff, ','),
		flattenFunc(race.Abilities, wAbl, ','),
		flattenFunc(race.Skills, wSkl, ','),
		writeChoices(race.Choices),
	}
}

func readAndInsert[T models.DbItem](filePath string, parse func([]string) T, insert func(*T) error) {
	if items, err := parseCsv(filePath, parse); err != nil {
		log.Printf("Failed to parse csv \"%s\": %v", filePath, err)
		return
	} else {
		for _, item := range items {
			if err = insert(&item); err != nil {
				log.Printf("Failed to insert item \"%s\" from \"%s\": %v", item.Pk(), filePath, err)
			}
		}
	}
}

func ReadAll(dir string) {
	readAndInsert(path.Join(dir, "skills.csv"), parseSkill, InsertSkill)
	readAndInsert(path.Join(dir, "abilities.csv"), parseAbility, InsertAbility)
	readAndInsert(path.Join(dir, "affinities.csv"), parseAffinity, InsertAffinity)
	readAndInsert(path.Join(dir, "backgrounds.csv"), parseBackground, InsertBackground)
	readAndInsert(path.Join(dir, "races.csv"), parseRace, InsertRace)
}

func fetchAndWrite[T any](fetch func() ([]T, error), location string, write func(T) []string) {
	if items, err := fetch(); err != nil {
		log.Printf("Failed to fetch items for csv \"%s\": %v", location, err)
	} else {
		if err = writeCsv(location, write, items); err != nil {
			log.Printf("Failed to write csv \"%s\": %v", location, err)
		}
	}
}

func SaveAll(targetDir string) {
	fetchAndWrite(FetchSkills, path.Join(targetDir, "skills.csv"), writeSkill)
	fetchAndWrite(FetchAbilities, path.Join(targetDir, "abilities.csv"), writeAbility)
	fetchAndWrite(func() ([]models.Affinity, error) { return FetchAffinities(nil) }, path.Join(targetDir, "affinities.csv"), writeAffinity)
	fetchAndWrite(func() ([]models.Background, error) { return FetchBackgrounds(false) }, path.Join(targetDir, "backgrounds.csv"), writeBackground)
	fetchAndWrite(func() ([]models.Race, error) { return FetchRaces(false) }, path.Join(targetDir, "races.csv"), writeRace)
}

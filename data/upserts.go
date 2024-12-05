package data

import (
	"NeonManager/models"
	"database/sql"
	"fmt"
	"strings"
)

// stat
var statInsert = `
INSERT INTO stats(cr, int, ins, ch, ag, dex, con, str) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
var statUpdate = `
UPDATE stats SET cr=?, int=?, ins=?, ch=?, ag=?, dex=?, con=?, str=?
WHERE id = ?`

// skill
var skillInsert = `
INSERT INTO skill(name, cost, extension, stat, description) 
SELECT ?, ?, (SELECT id FROM extension WHERE name = ?), ?, ?
WHERE NOT EXISTS (SELECT 1 FROM skill WHERE name = ?)`

// ability
var abilityInsert = `
INSERT INTO ability(name, cost, extension, effect)
SELECT ?, ?, (SELECT id FROM extension WHERE name = ?), ?
WHERE NOT EXISTS (SELECT 1 FROM ability WHERE name = ?)`
var abilityRequiresInsert = "INSERT INTO _abilities_requires(ability, reqType, rAffinity, rAbility, rSkill) VALUES "

// affinity
var affinityInsert = `
INSERT INTO affinity(name, cost, extension, description, isBoon)
SELECT ?, ?, (SELECT id FROM extension WHERE name = ?), ?, ?
WHERE NOT EXISTS (SELECT 1 FROM affinity WHERE name = ?)`

// background
var bgInsert = `
INSERT INTO background(name, cost, extension, description, stats)
SELECT ?, ?, (SELECT id FROM extension WHERE name = ?), ?, ?
WHERE NOT EXISTS (SELECT 1 FROM background WHERE name = ?)`
var bgAffinitiesInsert = "INSERT INTO _backgrounds_affinities(background, affinity, level, grp, gcount) VALUES "
var bgAbilitiesInsert = "INSERT INTO _backgrounds_abilities(background, ability, grp, gcount) VALUES "
var bgSkillsInsert = "INSERT INTO _backgrounds_skills(background, skill, level, grp, gcount) VALUES "

// race
var raceInsert = `
INSERT INTO race(name, cost, extension, description, stats)
SELECT ?, ?, (SELECT id FROM extension WHERE name = ?), ?, ?
WHERE NOT EXISTS (SELECT 1 FROM race WHERE name = ?)`
var raceBgInsert = "INSERT INTO _races_backgrounds(race, background, allowed) VALUES "
var raceAffinitiesInsert = "INSERT INTO _races_affinities(race, affinity, level, grp, gcount) VALUES "
var raceAbilitiesInsert = "INSERT INTO _races_abilities(race, ability, grp, gcount) VALUES "
var raceSkillsInsert = "INSERT INTO _races_skills(race, skill, level, grp, gcount) VALUES "

// character
var characterInsert = `
INSERT INTO character(name, gp, xp, description, stats, race, background, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
var characterUpdate = `
UPDATE character SET name = ?, gp = ?, xp = ?, stats = ?, race = ?, background = ?, description = ?
WHERE id = ?`
var characterClearLinks = `
DELETE FROM _characters_affinities WHERE character = ?;
DELETE FROM _characters_abilities WHERE character = ?;
DELETE FROM _characters_skills WHERE character = ?;
DELETE FROM _characters_extensions WHERE character = ?;`
var charExtInsert = "INSERT INTO _characters_extensions(character, extension) VALUES "
var charAffInsert = "INSERT INTO _characters_affinities(character, affinity, level) VALUES "
var charAblInsert = "INSERT INTO _characters_abilities(character, ability) VALUES "
var charSklInsert = "INSERT INTO _characters_skills(character, skill, level) VALUES "

func InsertStats(stats *models.Stats) (err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(statInsert); err != nil {
		return fmt.Errorf("failed to prepare stats insert statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	var res sql.Result
	if res, err = stmt.Exec((*stats).Cr, (*stats).Int, (*stats).Ins, (*stats).Ch, (*stats).Ag, (*stats).Dex,
		(*stats).Con, (*stats).Str); err != nil {
		return fmt.Errorf("failed to insert stats: %w", err)
	}
	var id int64
	if id, err = res.LastInsertId(); err == nil {
		(*stats).Id = int(id)
	}
	return
}

func UpdateStats(stats *models.Stats) (err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(statUpdate); err != nil {
		return fmt.Errorf("failed to prepare stats update statement: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec((*stats).Cr, (*stats).Int, (*stats).Ins, (*stats).Ch, (*stats).Ag, (*stats).Dex, (*stats).Con,
		(*stats).Str, (*stats).Id)
	return
}

func InsertSkill(skill *models.Skill) error {
	if stmt, err := db.Prepare(skillInsert); err != nil {
		return fmt.Errorf("failed to prepare skill insert statement: %w", err)
	} else {
		if _, err = stmt.Exec((*skill).Name, (*skill).Cost, (*skill).Extension, (*skill).Stat, (*skill).Description,
			(*skill).Name); err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
		_ = stmt.Close()
	}
	return nil
}

// ns transforms a string ptr to a NullString where nil and &"" are NULL values
func ns(t *string) sql.NullString {
	if t == nil || *t == "" {
		return sql.NullString{
			String: "NULL",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *t,
		Valid:  true,
	}
}

// valueStr creates the VALUES part of an INSERT statement, producing itemCount parentheses with propertyCount "?" each
func valuesStr(itemCount int, propertyCount int) string {
	str := "(" + strings.Repeat("?,", propertyCount)[:propertyCount*2-1] + "),"
	str = strings.Repeat(str, itemCount)
	return str[:len(str)-1]
}

func InsertAbility(ability *models.Ability) error {
	if stmt, err := db.Prepare(abilityInsert); err != nil {
		return fmt.Errorf("failed to prepare ability insert statement: %w", err)
	} else {
		if _, err = stmt.Exec((*ability).Name, (*ability).Cost, (*ability).Extension, (*ability).Effect,
			(*ability).Name); err != nil {
			return fmt.Errorf("failed to insert ability: %w", err)
		}
		_ = stmt.Close()
	}

	ln := len((*ability).Requires)
	if ln == 0 {
		return nil
	}
	values := make([]any, ln*5, 0)
	for _, req := range (*ability).Requires {
		switch req.RequiredType {
		case 0:
			values = append(values, (*ability).Name, req.RequiredType, nil, req.RAbility, nil)
		case 1:
			values = append(values, (*ability).Name, req.RequiredType, req.RAffinity, nil, nil)
		case 2:
			values = append(values, (*ability).Name, req.RequiredType, nil, nil, req.RSkill)
		}
	}
	if stmt, err := db.Prepare(abilityRequiresInsert + valuesStr(ln, 5)); err != nil {
		return fmt.Errorf("failed to prepare ability requirement insert statement for %s: %w", (*ability).Name, err)
	} else {
		if _, err = stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to insert ability requirements for %s: %w", (*ability).Name, err)
		}
		_ = stmt.Close()
	}
	return nil
}

func InsertAffinity(affinity *models.Affinity) error {
	if stmt, err := db.Prepare(affinityInsert); err != nil {
		return fmt.Errorf("failed to prepare affinity insert statement: %w", err)
	} else {
		if _, err = stmt.Exec((*affinity).Name, (*affinity).Cost, (*affinity).Extension, (*affinity).Description,
			(*affinity).IsBoon); err != nil {
			return fmt.Errorf("failed to insert affinity: %w", err)
		}
		_ = stmt.Close()
	}
	return nil
}

// links affinities using statement which should be like "INSERT INTO %(pk%, affinity, level, grp, gcount) VALUES"
func linkAffinities(statement string, pk string, affinities []*models.Affinity) error {
	ln := len(affinities)
	if ln == 0 {
		return nil
	}
	values := make([]any, ln*5, 0)
	for _, aff := range affinities {
		values = append(values, pk, (*aff).Name, (*aff).Level, 0, 0)
	}
	if stmt, err := db.Prepare(statement + valuesStr(ln, 5)); err != nil {
		return fmt.Errorf("failed to prepare affinities insert statement: %w", err)
	} else {
		if _, err = stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to link affinities: %w", err)
		}
		_ = stmt.Close()
	}
	return nil
}

// links abilities using statement which should be like "INSERT INTO %(pk%, ability, grp, gcount) VALUES"
func linkAbilities(statement string, pk string, abilities []*models.Ability) error {
	ln := len(abilities)
	if ln == 0 {
		return nil
	}
	values := make([]any, ln*4, 0)
	for _, abl := range abilities {
		values = append(values, pk, (*abl).Name, 0, 0)
	}
	if stmt, err := db.Prepare(statement + valuesStr(ln, 4)); err != nil {
		return fmt.Errorf("failed to prepare affinities insert statement: %w", err)
	} else {
		if _, err = stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to link affinities: %w", err)
		}
		_ = stmt.Close()
	}
	return nil
}

// links skills using statement which should be like "INSERT INTO %(%, skill, level, grp, gcount) VALUES"
func linkSkills(statement string, pk string, skills []*models.Skill) error {
	ln := len(skills)
	if ln == 0 {
		return nil
	}
	values := make([]any, ln*5, 0)
	for _, skl := range skills {
		values = append(values, pk, (*skl).Name, (*skl).Level, 0, 0)
	}
	if stmt, err := db.Prepare(statement + valuesStr(ln, 5)); err != nil {
		return fmt.Errorf("failed to prepare affinities insert statement: %w", err)
	} else {
		if _, err = stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to link affinities: %w", err)
		}
		_ = stmt.Close()
	}
	return nil
}

// links choices using the statements, see linkAffinities, linkAbilities, linkSkills
func linkChoices(affStmt string, ablStmt string, sklStmt string, pk string, choices []models.ChoiceGroup) error {
	if len(choices) == 0 {
		return nil
	}
	aff := make([]any, 0)
	abl := make([]any, 0)
	skl := make([]any, 0)
	for i, c := range choices {
		switch c.ChoiceType {
		case 0:
			for _, a := range c.Affinities {
				aff = append(aff, pk, (*a).Name, (*a).Level, i+1, c.Count)
			}
		case 1:
			for _, a := range c.Abilities {
				abl = append(abl, pk, (*a).Name, i+1, c.Count)
			}
		case 2:
			for _, a := range c.Skills {
				skl = append(skl, pk, (*a).Name, i+1, c.Count)
			}
		}
	}
	if len(aff) > 0 {
		if stmt, err := db.Prepare(affStmt + valuesStr(len(aff)/5, 5)); err != nil {
			return fmt.Errorf("failed to prepare affinity choice link statement: %w", err)
		} else {
			if _, err = stmt.Exec(aff...); err != nil {
				return fmt.Errorf("failed to link affinity choices for %s: %w", pk, err)
			}
			_ = stmt.Close()
		}
	}
	if len(abl) > 0 {
		if stmt, err := db.Prepare(ablStmt + valuesStr(len(abl)/4, 4)); err != nil {
			return fmt.Errorf("failed to prepare affinity choice link statement: %w", err)
		} else {
			if _, err = stmt.Exec(abl...); err != nil {
				return fmt.Errorf("failed to link affinity choices for %s: %w", pk, err)
			}
			_ = stmt.Close()
		}
	}
	if len(skl) > 0 {
		if stmt, err := db.Prepare(sklStmt + valuesStr(len(skl)/5, 5)); err != nil {
			return fmt.Errorf("failed to prepare affinity choice link statement: %w", err)
		} else {
			if _, err = stmt.Exec(skl...); err != nil {
				return fmt.Errorf("failed to link affinity choices for %s: %w", pk, err)
			}
			_ = stmt.Close()
		}
	}
	return nil
}

func InsertBackground(bg *models.Background) error {
	if (*bg).Stats.Id == 0 {
		if err := InsertStats(&bg.Stats); err != nil {
			return fmt.Errorf("failed to insert stats for background %s: %w", (*bg).Name, err)
		}
	}
	if stmt, err := db.Prepare(bgInsert); err != nil {
		return fmt.Errorf("failed to prepare background insert statement: %w", err)
	} else {
		if _, err = stmt.Exec((*bg).Name, (*bg).Cost, (*bg).Extension, (*bg).Description, (*bg).Stats); err != nil {
			return fmt.Errorf("failed to insert background: %w", err)
		}
		_ = stmt.Close()
	}
	if err := linkAffinities(bgAffinitiesInsert, (*bg).Name, append((*bg).Boons, (*bg).Banes...)); err != nil {
		return fmt.Errorf("failed to link affinities to background %s: %w", (*bg).Name, err)
	}
	if err := linkAbilities(bgAbilitiesInsert, (*bg).Name, (*bg).Abilities); err != nil {
		return fmt.Errorf("failed to link abilities to background %s: %w", (*bg).Name, err)
	}
	if err := linkSkills(bgSkillsInsert, (*bg).Name, (*bg).Skills); err != nil {
		return fmt.Errorf("failed to link skills to background %s: %w", (*bg).Name, err)
	}
	if err := linkChoices(bgAffinitiesInsert, bgAbilitiesInsert, bgSkillsInsert, (*bg).Name, (*bg).Choices); err != nil {
		return fmt.Errorf("failed to link choices to background %s: %w", (*bg).Name, err)
	}
	return nil
}

func InsertRace(race *models.Race) error {
	if (*race).Stats.Id == 0 {
		if err := InsertStats(&race.Stats); err != nil {
			return fmt.Errorf("failed to insert stats for race %s: %w", (*race).Name, err)
		}
	}
	if stmt, err := db.Prepare(raceInsert); err != nil {
		return fmt.Errorf("failed to prepare race insert statement: %w", err)
	} else {
		if _, err = stmt.Exec((*race).Name, (*race).Cost, (*race).Extension, (*race).Description, (*race).Stats); err != nil {
			return fmt.Errorf("failed to insert race: %w", err)
		}
		_ = stmt.Close()
	}
	if ln := len((*race).Backgrounds) + len((*race).NotBackgrounds); ln > 0 {
		values := make([]any, ln*3, 0)
		for _, bg := range (*race).Backgrounds {
			values = append(values, (*race).Name, bg.Name, 1)
		}
		for _, bg := range (*race).NotBackgrounds {
			values = append(values, (*race).Name, bg.Name, 0)
		}
		if stmt, err := db.Prepare(raceBgInsert + valuesStr(ln, 3)); err != nil {
			return fmt.Errorf("failed to prepare race link backgrounds statement for race %s: %w", (*race).Name, err)
		} else {
			if _, err = stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to link backgrounds to race %s: %w", (*race).Name, err)
			}
			_ = stmt.Close()
		}
	}
	if err := linkAffinities(raceAffinitiesInsert, (*race).Name, append((*race).Boons, (*race).Banes...)); err != nil {
		return fmt.Errorf("failed to link affinities to race %s: %w", (*race).Name, err)
	}
	if err := linkAbilities(raceAbilitiesInsert, (*race).Name, (*race).Abilities); err != nil {
		return fmt.Errorf("failed to link abilities to race %s: %w", (*race).Name, err)
	}
	if err := linkSkills(raceSkillsInsert, (*race).Name, (*race).Skills); err != nil {
		return fmt.Errorf("failed to link skills to race %s: %w", (*race).Name, err)
	}
	if err := linkChoices(raceAffinitiesInsert, raceAbilitiesInsert, raceSkillsInsert, (*race).Name, (*race).Choices); err != nil {
		return fmt.Errorf("failed to link choices to race %s: %w", (*race).Name, err)
	}
	return nil
}

func InsertUpdateCharacter(character *models.Character) error {
	if (*character).Stats.Id == 0 {
		if err := InsertStats(&character.Stats); err != nil {
			return fmt.Errorf("failed to insert stats for character %s: %w", (*character).Name, err)
		}
	} else {
		if err := UpdateStats(&character.Stats); err != nil {
			return fmt.Errorf("failed to update stats for character %s: %w", (*character).Name, err)
		}
	}
	if (*character).Id == 0 {
		if stmt, err := db.Prepare(characterInsert); err != nil {
			return fmt.Errorf("failed to prepare character insert statement: %w", err)
		} else {
			var res sql.Result
			if res, err = stmt.Exec((*character).Name, (*character).GP, (*character).XP, (*character).Stats.Id,
				ns(&(*character).Race.Name), ns(&(*character).Background.Name), (*character).Description); err != nil {
				return fmt.Errorf("failed to insert character %s: %w", (*character).Name, err)
			}
			var id int64
			if id, err = res.LastInsertId(); err != nil {
				return fmt.Errorf("failed to get id from inserted character %s: %w", (*character).Name, err)
			}
			(*character).Id = int(id)
		}
	} else {
		if stmt, err := db.Prepare(characterUpdate); err != nil {
			return fmt.Errorf("failed to prepare character update statement: %w", err)
		} else {
			if _, err = stmt.Exec((*character).Name, (*character).GP, (*character).XP, (*character).Stats.Id,
				ns(&(*character).Race.Name), ns(&(*character).Background.Name), (*character).Description,
				(*character).Id); err != nil {
				return fmt.Errorf("failed to update character %s: %w", (*character).Name, err)
			}
		}
	}
	id := (*character).Id
	if stmt, err := db.Prepare(characterClearLinks); err != nil {
		return fmt.Errorf("failed to prepare character clear links statement: %w", err)
	} else {
		if _, err = stmt.Exec(id, id, id, id); err != nil {
			return fmt.Errorf("failed to clear character links %s: %w", (*character).Name, err)
		}
	}
	ln := len((*character).Extensions)
	if ln > 0 {
		chExt := strings.Repeat("(?, (SELECT id FROM extension WHERE name = ?)),", ln)
		if stmt, err := db.Prepare(charExtInsert + chExt[:len(chExt)-1]); err != nil {
			return fmt.Errorf("failed to prepare character extensions insert statement: %w", err)
		} else {
			values := make([]any, len((*character).Extensions), 0)
			for _, v := range (*character).Extensions {
				values = append(values, id, v)
			}
			if _, err = stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to insert extensions for character %s: %w", (*character).Name, err)
			}
		}
	}
	ln = len((*character).Boons) + len((*character).Banes)
	if ln > 0 {
		if stmt, err := db.Prepare(charAffInsert + valuesStr(ln, 3)); err != nil {
			return fmt.Errorf("failed to prepare affinities link statement for character %s: %w", (*character).Name, err)
		} else {
			values := make([]any, ln*3, 0)
			for _, a := range append((*character).Boons, (*character).Banes...) {
				values = append(values, id, a.Name, a.Level)
			}
			if _, err = stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to link affinities for character %s: %w", (*character).Name, err)
			}
		}
	}
	ln = len((*character).Abilities)
	if ln > 0 {
		if stmt, err := db.Prepare(charAblInsert + valuesStr(ln, 2)); err != nil {
			return fmt.Errorf("failed to prepare abilities link statement for %s: %w", (*character).Name, err)
		} else {
			values := make([]any, ln*2, 0)
			for _, a := range (*character).Abilities {
				values = append(values, id, a.Name)
			}
			if _, err = stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to link abilities for character %s: %w", (*character).Name, err)
			}
		}
	}
	ln = len((*character).Skills)
	if ln > 0 {
		if stmt, err := db.Prepare(charSklInsert + valuesStr(ln, 3)); err != nil {
			return fmt.Errorf("failed to prepare skills link statement for %s: %w", (*character).Name, err)
		} else {
			values := make([]any, ln*3, 0)
			for _, s := range (*character).Skills {
				values = append(values, id, s.Name, s.Level)
			}
			if _, err = stmt.Exec(values...); err != nil {
				return fmt.Errorf("failed to link skills for character %s: %w", (*character).Name, err)
			}
		}
	}
	return nil
}

func InsertArticle(title string, text string, tags []string) error {
	var id int64
	if stmt, err := db.Prepare("INSERT INTO article(title, document) VALUES(?, ?)"); err != nil {
		return fmt.Errorf("failed to prepare article statement: %w", err)
	} else {
		var res sql.Result
		if res, err = stmt.Exec(title, text); err != nil {
			return fmt.Errorf("failed to insert article %s: %w", title, err)
		}
		if id, err = res.LastInsertId(); err != nil {
			return fmt.Errorf("failed to get id from article %s: %w", title, err)
		}
	}
	if tags == nil || len(tags) == 0 {
		return nil
	}
	values := make([]any, 0)
	for _, tag := range tags {
		values = append(values, id, tag)
	}
	if stmt, err := db.Prepare("INSERT INTO _articles_tags(article, tag) VALUES " + valuesStr(len(tags), 2)); err != nil {
		return fmt.Errorf("failed to prepare article tags statement: %w", err)
	} else {
		if _, err = stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to insert tags for article %s: %w", title, err)
		}
	}
	return nil
}

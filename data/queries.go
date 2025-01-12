package data

import (
	"NeonManager/models"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"
)

var db *sql.DB

var raceQuery = `
SELECT * FROM (SELECT r.name, r.cost, e.name, r.description, s.id, s.cr, s.int, s.ins, s.ch, s.ag, s.dex, s.con, s.str,
   COALESCE(GROUP_CONCAT(CASE WHEN b.allowed = 0 THEN b.background END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN b.allowed = 1 THEN b.background END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 1 AND af.grp = 0 THEN bb.name || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 0 AND af.grp = 0 THEN bb.name || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN a.grp = 0 THEN a.ability END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN sk.grp = 0 THEN sk.skill || '@' || sk.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN af.grp <> 0 THEN af.grp || ':' || af.gcount || ':' || af.affinity || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN a.grp <> 0 THEN a.grp || ':' || a.gcount || ':' || a.ability END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN sk.grp <> 0 THEN sk.grp || ':' || sk.gcount || ':' || sk.skill || '@' || sk.level END, ','), '')
FROM race r
   JOIN extension e ON r.extension = e.id
   JOIN stats s ON r.stats = s.id
   JOIN _races_backgrounds b ON r.name = b.race
   JOIN _races_affinities af ON r.name = af.race
   JOIN affinity bb ON af.affinity = bb.name
   JOIN _races_abilities a ON r.name = a.race
   JOIN _races_skills sk ON r.name = sk.race
GROUP BY r.name
HAVING COUNT(*) > 0
)
`
var bgQuery = `
SELECT * FROM (SELECT b.name, b.cost, e.name, b.description, s.id, s.cr, s.int, s.ins, s.ch, s.ag, s.dex, s.con, s.str,
   COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 1 AND af.grp = 0 THEN bb.name || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 0 AND af.grp = 0 THEN bb.name || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN a.grp = 0 THEN a.ability END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN sk.grp = 0 THEN sk.skill || '@' || sk.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN af.grp <> 0 THEN af.grp || ':' || af.gcount || ':' || af.affinity || '@' || af.level END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN a.grp <> 0 THEN a.grp || ':' || a.gcount || ':' || a.ability END, ','), ''),
   COALESCE(GROUP_CONCAT(CASE WHEN sk.grp <> 0 THEN sk.grp || ':' || sk.gcount || ':' || sk.skill || '@' || sk.level END, ','), '')
FROM background b
    JOIN extension e ON e.id = b.extension 
    JOIN stats s ON s.id = b.stats 
    JOIN _backgrounds_affinities af ON b.name = af.background
    JOIN affinity bb ON af.affinity = bb.name
    JOIN _backgrounds_abilities a ON b.name = a.background
    JOIN _backgrounds_skills sk ON b.name = sk.background
GROUP BY b.name
HAVING COUNT(*) > 0
)
`
var affinityQuery = `
SELECT a.name, a.cost, e.name, a.description, a.isBoon
FROM affinity a
    JOIN extension e ON a.extension = e.id
`
var abilityQuery = `
SELECT * FROM (SELECT a.name, a.cost, e.name, a.effect,
       COALESCE(GROUP_CONCAT(CASE WHEN r.reqType = 0 THEN r.rAbility || '@' || r.id END, ','), ''),
       COALESCE(GROUP_CONCAT(CASE WHEN r.reqType = 1 THEN r.rAffinity || '@' || r.id END, ','), ''),
       COALESCE(GROUP_CONCAT(CASE WHEN r.reqType = 2 THEN r.rSkill || '@' || r.id END, ','), '') /* might also need skill level */
FROM ability a
    JOIN extension e ON a.extension = e.id
    JOIN _abilities_requires r ON a.name = r.ability
GROUP BY a.name
HAVING COUNT(*) > 0
)
`
var skillQuery = `
SELECT s.name, s.cost, e.name, s.stat, s.description
FROM skill s
    JOIN extension e ON e.id = s.extension
`
var characterQuery = `
SELECT * FROM (SELECT c.id, c.name, c.gp, c.xp, s.id, s.cr, s.int, s.ins, s.ch, s.ag, s.dex, s.con, s.str, c.race, 
       c.background, c.description, COALESCE(GROUP_CONCAT(e.name, ','), ''), 
       COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 1 THEN bb.name || '@' || af.level END, ','), ''),
       COALESCE(GROUP_CONCAT(CASE WHEN bb.isBoon = 0 THEN bb.name || '@' || af.level END, ','), ''),
       COALESCE(GROUP_CONCAT(ca.ability, ','), ''), GROUP_CONCAT(sk.skill || '@' || sk.level, ',')
FROM character c
    JOIN stats s ON s.id = c.stats
    JOIN _characters_extensions ce ON c.id = ce.character
    JOIN extension e ON ce.extension = e.id
    JOIN _characters_affinities af ON c.id = af.character
    JOIN affinity bb ON af.affinity = bb.name
    JOIN _characters_abilities ca ON c.id = ca.character
    JOIN _characters_skills sk ON c.id = sk.character
GROUP BY c.id
HAVING COUNT(*) > 0
)
`
var articleParentQuery = `
WITH RECURSIVE topic AS (
    SELECT id, parent, title
    FROM article
    WHERE parent IS NULL

    UNION ALL

    SELECT a.id, a.parent, a.title
    FROM article a
    INNER JOIN topic t ON a.parent = t.id
)
SELECT id, parent, title FROM topic
ORDER BY parent, id
`
var articleTitleQuery = `
WITH RECURSIVE topic AS (
    SELECT id, parent, title
    FROM article
    WHERE title = ?

    UNION ALL

    SELECT a.id, a.parent, a.title
    FROM article a
    INNER JOIN topic t ON a.parent = t.id
)
SELECT id, parent, title FROM topic
ORDER BY parent, id /* parent id should always be lower than child ids */
`

// Init initializes db connection
func Init() error {
	if db != nil {
		log.Println("Already initialized")
		return nil
	}

	if database, err := sql.Open("sqlite3", "./data.db"); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	} else {
		db = database
	}

	// enable foreign key validation
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign key validation: %w", err)
	}

	if res, err := db.Query("SELECT * FROM sqlite_master"); err != nil {
		return fmt.Errorf("failed to retrieve sqlite_master: %w", err)
	} else {
		if !res.Next() {
			_ = res.Close()
			return createSchema()
		}
		_ = res.Close()
	}
	return nil
}

// createSchema creates the schema for the db, may also be called if the schema already exists (no-op)
func createSchema() error {
	var sqlStr string
	if bytes, err := os.ReadFile("./sql/data.sql"); err != nil {
		return fmt.Errorf("failed to retrieve schema from file: %w", err)
	} else {
		sqlStr = string(bytes)
	}

	if res, err := db.Exec(sqlStr); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	} else {
		rows, _ := res.RowsAffected()
		log.Printf("Schema created (%d rows affected)", rows)
	}
	return nil
}

// fetchItem populates the item via the scan function from the database using the query with item.Pk() as argument
func fetchItem[T models.DbItem](query string, item *T, scan func(*sql.Rows, *T) error) error {
	if stmt, err := db.Prepare(query); err != nil {
		return fmt.Errorf("failed to prepare fetchItem query: %w", err)
	} else {
		defer func() {
			_ = stmt.Close()
		}()
		rows, _ := stmt.Query((*item).Pk())
		if rows.Next() {
			if err = scan(rows, item); err != nil {
				return fmt.Errorf("failed to scan fetchItem from row: %w", err)
			}
		} else {
			return fmt.Errorf("item with pk \"%s\" not found (query: \"%s\")", (*item).Pk(), query)
		}
	}
	return nil
}

// fetchItems fetches all items via the scan function from the result of the query
func fetchItems[T any](stmt *sql.Stmt, args any, scan func(*sql.Rows, *T) error) ([]T, error) {
	var rows *sql.Rows
	var err error
	if args != nil {
		rows, err = stmt.Query(args)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to prepare fetchItems query: %w", err)
	} else {
		defer func() {
			_ = rows.Close()
		}()
		items := make([]T, 0)
		for rows.Next() {
			var item T
			if err = scan(rows, &item); err != nil {
				return items, fmt.Errorf("failed to scan fetchItems row: %w", err)
			} else {
				items = append(items, item)
			}
		}
		return items, nil
	}
}

/* splitNumber splits a string into text as prefix and number as postfix on the given sep */
func splitNumber(s string, sep string) (string, int) {
	split := strings.Split(s, sep)
	if len(split) != 2 {
		return s, 0
	}
	if i, err := strconv.Atoi(split[1]); err != nil {
		return split[0], 0
	} else {
		return split[0], i
	}
}

func resolveLevelList[T models.LevelItem](csv string, _ T) []*T {
	names := strings.Split(csv, ",")
	items := make([]*T, len(names))
	for i, itemStr := range names {
		items[i] = new(T)
		name, level := splitNumber(itemStr, "@")
		(*items[i]).SetName(name)
		(*items[i]).SetLevel(level)
	}
	return items
}

func resolveList[T models.DbItem](csv string, _ T) []*T {
	names := strings.Split(csv, ",")
	items := make([]*T, len(names))
	for i, name := range names {
		items[i] = new(T)
		(*items[i]).SetName(name)
	}
	return items
}

func resolveGroups(csv string, cType int) []models.ChoiceGroup {
	groupSplits := make([][]string, 0)
	groupCount := 0
	for i, group := range strings.Split(csv, ",") {
		groupSplits = append(groupSplits, strings.Split(group, ":"))
		if n, err := strconv.Atoi(groupSplits[i][0]); err == nil {
			if n > groupCount {
				groupCount = n
			}
		}
	}
	groups := make([]models.ChoiceGroup, groupCount)
	for _, group := range groupSplits {
		groupNumber, _ := strconv.Atoi(group[0])
		itemCount, _ := strconv.Atoi(group[1])
		switch cType {
		case 0:
			if groups[groupNumber].Affinities == nil {
				groups[groupNumber].Count = itemCount
				groups[groupNumber].Affinities = make([]*models.Affinity, 0)
			}
			name, level := splitNumber(group[1], "@")
			groups[groupCount].Affinities = append(groups[groupNumber].Affinities,
				&models.Affinity{Name: name, Level: level})
		case 1:
			if groups[groupNumber].Abilities == nil {
				groups[groupNumber].Count = itemCount
				groups[groupNumber].Abilities = make([]*models.Ability, 0)
			}
			groups[groupNumber].Abilities = append(groups[groupNumber].Abilities, &models.Ability{Name: group[2]})
		case 2:
			if groups[groupNumber].Skills == nil {
				groups[groupNumber].Count = itemCount
				groups[groupNumber].Skills = make([]*models.Skill, 0)
			}
			name, level := splitNumber(group[1], "@")
			groups[groupCount].Skills = append(groups[groupNumber].Skills,
				&models.Skill{Name: name, Level: level})
		}
	}
	return groups
}

func scanRace(rows *sql.Rows, item *models.Race) (err error) {
	var bgStr string
	var notBgStr string
	var boonStr string
	var baneStr string
	var abilityStr string
	var skillStr string
	var affGroups string
	var ablGroups string
	var sklGroups string
	if err = rows.Scan(&(*item).Name, &(*item).Cost, &(*item).Extension, &(*item).Description, &(*item).Stats.Id,
		&(*item).Stats.Cr, &(*item).Stats.Int, &(*item).Stats.Ins, &(*item).Stats.Ch, &(*item).Stats.Ag,
		&(*item).Stats.Dex, &(*item).Stats.Con, &(*item).Stats.Str, &notBgStr, &bgStr, &boonStr, &baneStr, &abilityStr,
		&skillStr, &affGroups, &ablGroups, &sklGroups); err == nil {

		(*item).Backgrounds = resolveList(bgStr, models.Background{})
		(*item).NotBackgrounds = resolveList(notBgStr, models.Background{})

		(*item).Boons = resolveLevelList(boonStr, models.Affinity{})
		(*item).Banes = resolveLevelList(baneStr, models.Affinity{})

		(*item).Abilities = resolveList(abilityStr, models.Ability{})

		(*item).Skills = resolveLevelList(skillStr, models.Skill{})

		groups := resolveGroups(affGroups, 0)
		groups = append(groups, resolveGroups(ablGroups, 1)...)
		groups = append(groups, resolveGroups(sklGroups, 2)...)
		if len(groups) > 0 {
			(*item).Choices = groups
		}
	}
	return err
}

func scanBg(rows *sql.Rows, item *models.Background) (err error) {
	var boonStr string
	var baneStr string
	var abilityStr string
	var skillStr string
	var affGroups string
	var ablGroups string
	var sklGroups string
	if err = rows.Scan(&(*item).Name, &(*item).Cost, &(*item).Extension, &(*item).Description, &(*item).Stats.Id,
		&(*item).Stats.Cr, &(*item).Stats.Int, &(*item).Stats.Ins, &(*item).Stats.Ch, &(*item).Stats.Ag,
		&(*item).Stats.Dex, &(*item).Stats.Con, &(*item).Stats.Str, &boonStr, &baneStr, &abilityStr, &skillStr,
		&affGroups, &ablGroups, &sklGroups); err == nil {

		(*item).Boons = resolveLevelList(boonStr, models.Affinity{})
		(*item).Banes = resolveLevelList(baneStr, models.Affinity{})

		(*item).Abilities = resolveList(abilityStr, models.Ability{})

		(*item).Skills = resolveLevelList(skillStr, models.Skill{})

		groups := resolveGroups(affGroups, 0)
		groups = append(groups, resolveGroups(ablGroups, 1)...)
		groups = append(groups, resolveGroups(sklGroups, 2)...)
		if len(groups) > 0 {
			(*item).Choices = groups
		}
	}
	return
}

func scanAffinity(rows *sql.Rows, item *models.Affinity) error {
	return rows.Scan(&(*item).Name, &(*item).Cost, &(*item).Extension, &(*item).Description, &(*item).IsBoon)
}

type bareLItem struct {
	id   int
	name string
}

func (b bareLItem) SetName(name string) {
	b.name = name
}

func (b bareLItem) SetLevel(level int) {
	b.id = level
}

func (b bareLItem) Pk() string {
	return b.name
}

func (b bareLItem) GetLevel() int {
	return b.id
}

func scanAbility(rows *sql.Rows, item *models.Ability) (err error) {
	var ablStr string
	var affStr string
	var sklStr string
	if err = rows.Scan(&(*item).Name, &(*item).Cost, &(*item).Extension, &(*item).Effect, &ablStr, &affStr,
		&sklStr); err == nil {
		name := (*item).Name
		reqs := make([]models.AbilityReq, 0)
		for _, abl := range resolveLevelList(ablStr, bareLItem{}) {
			reqs = append(reqs, models.AbilityReq{Ability: name, Id: abl.id, RAbility: abl.name, RequiredType: models.ChoiceAbility})
		}
		for _, aff := range resolveLevelList(affStr, bareLItem{}) {
			reqs = append(reqs, models.AbilityReq{Ability: name, Id: aff.id, RAffinity: aff.name, RequiredType: models.ChoiceAffinity})
		}
		for _, skl := range resolveLevelList(sklStr, bareLItem{}) {
			reqs = append(reqs, models.AbilityReq{Ability: name, Id: skl.id, RSkill: skl.name, RequiredType: models.ChoiceSkill})
		}
		(*item).Requires = reqs
	}
	return
}

func scanSkill(rows *sql.Rows, item *models.Skill) error {
	return rows.Scan(&(*item).Name, &(*item).Cost, &(*item).Extension, &(*item).Stat, &(*item).Description)
}

func scanCharacter(rows *sql.Rows, item *models.Character) (err error) {
	var extStr string
	var boonStr string
	var baneStr string
	var abilityStr string
	var sklStr string
	if err = rows.Scan(&(*item).Id, &(*item).Name, &(*item).GP, &(*item).XP, &(*item).Stats.Id, &(*item).Stats.Cr,
		&(*item).Stats.Int, &(*item).Stats.Ins, &(*item).Stats.Ch, &(*item).Stats.Ag, &(*item).Stats.Dex,
		&(*item).Stats.Con, &(*item).Stats.Str, &(*item).Race.Name, &(*item).Background.Name, &(*item).Description,
		&extStr, &boonStr, &baneStr, &abilityStr, &sklStr); err == nil {

		item.Extensions = strings.Split(extStr, ",")

		item.Boons = resolveLevelList(boonStr, models.Affinity{})
		item.Banes = resolveLevelList(baneStr, models.Affinity{})

		item.Abilities = resolveList(abilityStr, models.Ability{})

		item.Skills = resolveLevelList(sklStr, models.Skill{})
	}
	return
}

func toPkList[T models.DbItem](items []*T) []string {
	pks := make([]string, len(items))
	for i, item := range items {
		pks[i] = (*item).Pk()
	}
	return pks
}

func populateLevelItem[T models.LevelItem](stmt *sql.Stmt, items *[]*T, pks []string, scan func(*sql.Rows, *T) error) (err error) {
	var dbItems []T
	if dbItems, err = fetchItems(stmt, pks, scan); err != nil {
		return
	}
	// merge affinities since some might already have data stored
	for i, item := range *items {
		j := 0
		for ; dbItems[j].Pk() != (*item).Pk() && j < len(dbItems); j++ {
		}
		if j == len(dbItems) {
			return fmt.Errorf("could not find %s to pupulate affinities", (*item).Pk())
		}
		// affinities = append(affinities[:j], affinities[j+1:]...) // FIXME performance
		dbItem := dbItems[j]
		dbItem.SetLevel((*item).GetLevel())
		(*items)[i] = &dbItem // can change item in iterator?
	}
	return
}

func populateRaceBg(race *models.Race, stmt *sql.Stmt) (err error) {
	pks := toPkList(race.Backgrounds)
	var bgs []models.Background
	if bgs, err = fetchItems(stmt, pks, scanBg); err != nil {
		err = fmt.Errorf("failed to populate backgrounds for %s: %w", race.Name, err)
	} else {
		// note: race.Background already exists => for simplicity overwrite all items instead of
		// matching by name since no additional data is stored
		for j, bg := range bgs {
			(*race).Backgrounds[j] = &bg
		}
	}
	pks = toPkList(race.NotBackgrounds)
	if bgs, err = fetchItems(stmt, pks, scanBg); err != nil {
		if err != nil {
			err = fmt.Errorf("%w\n  failed to populate notBackgrounds for %s: %w", err, race.Name, err)
		} else {
			err = fmt.Errorf("failed to populate notBackgrounds for %s: %w", race.Name, err)
		}
	} else {
		for j, bg := range bgs {
			(*race).NotBackgrounds[j] = &bg
		}
	}
	return
}

// FetchRaces fetches all races
// if populate is true the values of the subitems like backgrounds/boons/etc
// will be populated otherwise they will just list the names
func FetchRaces(populate bool) (races []models.Race, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(raceQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchRaces query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	if races, err = fetchItems(stmt, nil, scanRace); err != nil {
		return nil, fmt.Errorf("failed to fetchItems: %w", err)
	}
	if populate && err == nil {
		errs := make([]error, 0)
		if stmt, err = db.Prepare(bgQuery + " WHERE b.name IN ?"); err != nil {
			return nil, fmt.Errorf("failed to prepare fetchRaces background population query: %w", err)
		} else {
			for i := range races {
				if err = populateRaceBg(&races[i], stmt); err != nil {
					errs = append(errs, err)
				}
			}
			_ = stmt.Close()
		}
		if stmt, err = db.Prepare(affinityQuery + " WHERE a.name IN ?"); err == nil {
			for i, race := range races {
				if err = populateLevelItem(stmt, &races[i].Boons, toPkList(races[i].Boons), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate boons for %s: %w", race.Name, err))
				}
				if err = populateLevelItem(stmt, &races[i].Banes, toPkList(races[i].Banes), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate banes for %s: %w", race.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, fmt.Errorf("failed to prepare fetchRaces affinity population query: %w", err))
		}

		if stmt, err = db.Prepare(abilityQuery + " WHERE a.name IN ?"); err == nil {
			for i, race := range races {
				if err = populateLevelItem(stmt, &races[i].Abilities, toPkList(races[i].Abilities), scanAbility); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate abilities for %s: %w", race.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if stmt, err = db.Prepare(skillQuery + " WHERE b.name IN ?"); err == nil {
			for i, race := range races {
				if err = populateLevelItem(stmt, &races[i].Skills, toPkList(races[i].Skills), scanSkill); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate skills for %s: %w", race.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			if len(errs) == 1 {
				err = fmt.Errorf("an error occured while populating races: %w", errs[0])
			} else {
				err = fmt.Errorf("multiple errors occured while populating races:\n")
				for _, subError := range errs {
					err = fmt.Errorf("%w\n  %w", err, subError)
				}
			}
		}
	}
	return
}

func FetchRace(name string) (race models.Race, err error) {
	race.Name = name
	if err = fetchItem(raceQuery+" WHERE r.name = ?", &race, scanRace); err != nil {
		return race, fmt.Errorf("failed to fetch race %s: %w", name, err)
	}
	errs := make([]error, 0)
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(bgQuery + " WHERE b.name IN ?"); err != nil {
		return race, fmt.Errorf("failed to prepare fetchRaces background population query: %w", err)
	} else {
		if err = populateRaceBg(&race, stmt); err != nil {
			errs = append(errs, err)
		}
		_ = stmt.Close()
	}
	if stmt, err = db.Prepare(affinityQuery + " WHERE a.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &race.Boons, toPkList(race.Boons), scanAffinity); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate boons for %s: %w", race.Name, err))
		}
		if err = populateLevelItem(stmt, &race.Banes, toPkList(race.Banes), scanAffinity); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate banes for %s: %w", race.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, fmt.Errorf("failed to prepare fetchRaces affinity population query: %w", err))
	}

	if stmt, err = db.Prepare(abilityQuery + " WHERE a.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &race.Abilities, toPkList(race.Abilities), scanAbility); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate abilities for %s: %w", race.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, err)
	}
	if stmt, err = db.Prepare(skillQuery + " WHERE b.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &race.Skills, toPkList(race.Skills), scanSkill); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate skills for %s: %w", race.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		if len(errs) == 1 {
			err = fmt.Errorf("an error occured while populating races: %w", errs[0])
		} else {
			err = fmt.Errorf("multiple errors occured while populating races:\n")
			for _, subError := range errs {
				err = fmt.Errorf("%w\n  %w", err, subError)
			}
		}
	}
	return
}

func FetchBackgrounds(populate bool) (bgs []models.Background, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(bgQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchBgs query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	if bgs, err = fetchItems(stmt, nil, scanBg); err != nil {
		return bgs, fmt.Errorf("failed to fetchItems: %w", err)
	}
	if populate && err == nil {
		errs := make([]error, 0)
		if stmt, err = db.Prepare(affinityQuery + " WHERE a.name IN ?"); err == nil {
			for i, bg := range bgs {
				if err = populateLevelItem(stmt, &bgs[i].Boons, toPkList(bgs[i].Boons), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate boons for %s: %w", bg.Name, err))
				}
				if err = populateLevelItem(stmt, &bgs[i].Banes, toPkList(bgs[i].Banes), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate banes for %s: %w", bg.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, fmt.Errorf("failed to prepare fetchRaces affinity population query: %w", err))
		}

		if stmt, err = db.Prepare(abilityQuery + " WHERE a.name IN ?"); err == nil {
			for i, bg := range bgs {
				if err = populateLevelItem(stmt, &bgs[i].Abilities, toPkList(bgs[i].Abilities), scanAbility); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate abilities for %s: %w", bg.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if stmt, err = db.Prepare(skillQuery + " WHERE b.name IN ?"); err == nil {
			for i, bg := range bgs {
				if err = populateLevelItem(stmt, &bgs[i].Skills, toPkList(bgs[i].Skills), scanSkill); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate skills for %s: %w", bg.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			if len(errs) == 1 {
				err = fmt.Errorf("an error occured while populating backgrounds: %w", errs[0])
			} else {
				err = fmt.Errorf("multiple errors occured while populating backgrounds:\n")
				for _, subError := range errs {
					err = fmt.Errorf("%w\n  %w", err, subError)
				}
			}
		}
	}
	return
}

func FetchBackground(name string) (bg models.Background, err error) {
	bg.Name = name
	if err = fetchItem(bgQuery+" WHERE b.name LIKE ?", &bg, scanBg); err != nil {
		err = fmt.Errorf("failed to fetch background %s: %w", name, err)
	}
	errs := make([]error, 0)
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(affinityQuery + " WHERE a.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &bg.Boons, toPkList(bg.Boons), scanAffinity); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate boons for %s: %w", bg.Name, err))
		}
		if err = populateLevelItem(stmt, &bg.Banes, toPkList(bg.Banes), scanAffinity); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate banes for %s: %w", bg.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, fmt.Errorf("failed to prepare fetchRaces affinity population query: %w", err))
	}

	if stmt, err = db.Prepare(abilityQuery + " WHERE a.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &bg.Abilities, toPkList(bg.Abilities), scanAbility); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate abilities for %s: %w", bg.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, err)
	}
	if stmt, err = db.Prepare(skillQuery + " WHERE b.name IN ?"); err == nil {
		if err = populateLevelItem(stmt, &bg.Skills, toPkList(bg.Skills), scanSkill); err != nil {
			errs = append(errs, fmt.Errorf("failed to populate skills for %s: %w", bg.Name, err))
		}
		_ = stmt.Close()
	} else {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		if len(errs) == 1 {
			err = fmt.Errorf("an error occured while populating backgrounds: %w", errs[0])
		} else {
			err = fmt.Errorf("multiple errors occured while populating backgrounds:\n")
			for _, subError := range errs {
				err = fmt.Errorf("%w\n  %w", err, subError)
			}
		}
	}
	return
}

func FetchAffinities(isBoon *bool) (_ []models.Affinity, err error) {
	var stmt *sql.Stmt
	if isBoon != nil {
		if stmt, err = db.Prepare(affinityQuery + " WHERE a.isBoon = ?"); err != nil {
			return nil, fmt.Errorf("failed to prepare fetchAffinities query: %w", err)
		}
	} else {
		if stmt, err = db.Prepare(affinityQuery); err != nil {
			return nil, fmt.Errorf("failed to prepare fetchAffinities query: %w", err)
		}
	}
	defer func() {
		_ = stmt.Close()
	}()
	return fetchItems(stmt, isBoon, scanAffinity)
}

func FetchAffinity(name string) (aff models.Affinity, err error) {
	aff.Name = name
	if err = fetchItem(affinityQuery+" WHERE a.name LIKE ?", &aff, scanAffinity); err != nil {
		err = fmt.Errorf("failed to fetch affinity %s: %w", name, err)
	}
	return
}

func FetchAbilities() (abilities []models.Ability, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(abilityQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchAbilities query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	return fetchItems(stmt, nil, scanAbility)
}

func FetchAbility(name string) (ability models.Ability, err error) {
	ability.Name = name
	if err = fetchItem(abilityQuery+" WHERE a.name LIKE ?", &ability, scanAbility); err != nil {
		err = fmt.Errorf("failed to fetch ability %s: %w", name, err)
	}
	return
}

func FetchSkills() (skills []models.Skill, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(skillQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchSkills query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	return fetchItems(stmt, nil, scanSkill)
}

func FetchSkill(name string) (skill models.Skill, err error) {
	skill.Name = name
	if err = fetchItem(skillQuery+" WHERE a.name LIKE ?", &skill, scanSkill); err != nil {
		err = fmt.Errorf("failed to fetch skill %s: %w", name, err)
	}
	return
}

func FetchCharacters(populate bool) (chars []models.Character, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(characterQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchCharacters query: %w", err)
	}
	chars, err = fetchItems(stmt, nil, scanCharacter)
	_ = stmt.Close()
	if populate && err == nil {
		errs := make([]error, 0)
		for _, chr := range chars {
			if err = fetchItem(raceQuery+" WHERE r.name = ?", &chr.Race, scanRace); err != nil {
				errs = append(errs, err)
			}
			if err = fetchItem(bgQuery+" WHERE b.name = ?", &chr.Background, scanBg); err != nil {
				errs = append(errs, err)
			}
		}

		if stmt, err = db.Prepare(affinityQuery + " WHERE a.name IN ?"); err == nil {
			for i, chr := range chars {
				if err = populateLevelItem(stmt, &chars[i].Boons, toPkList(chars[i].Boons), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate boons for %s: %w", chr.Name, err))
				}
				if err = populateLevelItem(stmt, &chars[i].Banes, toPkList(chars[i].Banes), scanAffinity); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate banes for %s: %w", chr.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, fmt.Errorf("failed to prepare fetchRaces affinity population query: %w", err))
		}

		if stmt, err = db.Prepare(abilityQuery + " WHERE a.name IN ?"); err == nil {
			for i, chr := range chars {
				if err = populateLevelItem(stmt, &chars[i].Abilities, toPkList(chars[i].Abilities), scanAbility); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate abilities for %s: %w", chr.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if stmt, err = db.Prepare(skillQuery + " WHERE b.name IN ?"); err == nil {
			for i, chr := range chars {
				if err = populateLevelItem(stmt, &chars[i].Skills, toPkList(chars[i].Skills), scanSkill); err != nil {
					errs = append(errs, fmt.Errorf("failed to populate skills for %s: %w", chr.Name, err))
				}
			}
			_ = stmt.Close()
		} else {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			if len(errs) == 1 {
				err = fmt.Errorf("an error occured while populating character: %w", errs[0])
			} else {
				err = fmt.Errorf("multiple errors occured while populating character:\n")
				for _, subError := range errs {
					err = fmt.Errorf("%w\n  %w", err, subError)
				}
			}
		}
	}
	return
}

type paragraph struct {
	text    string
	css     string
	ordinal int
}

type table struct {
	title   string
	ordinal int
	header  []string
	rows    [][]string
}

type article struct {
	id         int64
	title      string
	sub        []*article
	paragraphs []paragraph
	table      *table
}

func fetchTable(article int64) (tbl *table, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(`SELECT ptable.id, ordinal, MAX(tc.col), MAX(tc.row), ptable.title FROM ptable 
    JOIN table_col tc on ptable.id = tc.id WHERE ptable.article = ? HAVING COUNT(*) > 0`, article); err != nil {
		return
	}
	var id int
	var colCnt int
	var rowCnt int
	if rows.Next() {
		tbl = &table{}
		if err = rows.Scan(&id, &tbl.ordinal, &colCnt, &rowCnt, &tbl.title); err != nil {
			return
		}
	} else {
		return nil, nil
	}
	_ = rows.Close()
	if rows, err = db.Query("SELECT row, col, value FROM table_col WHERE id = ? ORDER BY row", id); err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	tbl.header = make([]string, colCnt)
	tbl.rows = make([][]string, rowCnt)
	for rows.Next() {
		row := struct {
			row   int
			col   int
			value string
		}{}
		if err = rows.Scan(&row.row, &row.col, &row.value); err != nil {
			return
		}
		if row.row == 0 {
			tbl.header[row.col-1] = row.value
		} else {
			if tbl.rows[row.row-1] == nil {
				tbl.rows[row.row-1] = make([]string, colCnt)
			}
			tbl.rows[row.row-1][row.col-1] = row.value
		}
	}
	return
}

func fetchParagraphs(article int64) (p []paragraph, err error) {
	var rows *sql.Rows
	if rows, err = db.Query("SELECT ordinal, text, COALESCE(css, '') FROM paragraph WHERE article = ?", article); err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	p = make([]paragraph, 0)
	for rows.Next() {
		para := paragraph{}
		if err = rows.Scan(&para.ordinal, &para.text, &para.css); err != nil {
			return
		}
		p = append(p, para)
	}
	return
}

func fetchArticle(title string) (result *article, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(articleTitleQuery); err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	var rows *sql.Rows
	if rows, err = stmt.Query(title); err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	articleMap := make(map[int64]*article)
	for rows.Next() {
		art := article{}
		var p sql.NullInt64
		if err = rows.Scan(&art.id, &p, &art.title); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		var parent int64
		if p.Valid {
			parent = p.Int64
		} else {
			parent = 0
		}
		if art.title == title {
			result = &art
			articleMap[art.id] = result
			continue
		}
		articleMap[art.id] = &art
		if (*articleMap[parent]).sub == nil {
			(*articleMap[parent]).sub = make([]*article, 0)
		}
		(*articleMap[parent]).sub = append(articleMap[parent].sub, &art)
	}
	_ = rows.Close()
	_ = stmt.Close()
	for k := range articleMap {
		if articleMap[k].paragraphs, err = fetchParagraphs(articleMap[k].id); err != nil {
			return nil, fmt.Errorf("failed to fetch paragraphs: %w", err)
		}
		if (articleMap[k]).table, err = fetchTable(articleMap[k].id); err != nil {
			return nil, fmt.Errorf("failed to fetch table: %w", err)
		}
	}
	return
}

func fetchArticles() (articles []*article, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(articleParentQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchArticles query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	var rows *sql.Rows
	if rows, err = stmt.Query(); err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	articleMap := make(map[int64]*article)
	articleMap[0] = &article{}
	for rows.Next() {
		art := article{}
		var p sql.NullInt64
		if err = rows.Scan(&art.id, &p, &art.title); err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)
		}
		var parent int64
		if p.Valid {
			parent = p.Int64
		} else {
			parent = 0
		}
		articleMap[art.id] = &art
		if (*articleMap[parent]).sub == nil {
			(*articleMap[parent]).sub = make([]*article, 0)
		}
		articleMap[parent].sub = append(articleMap[parent].sub, &art)
	}
	articles = (*articleMap[0]).sub
	_ = rows.Close()
	_ = stmt.Close()
	for k := range articleMap {
		if articleMap[k].paragraphs, err = fetchParagraphs(articleMap[k].id); err != nil {
			return nil, fmt.Errorf("failed to fetch paragraphs: %w", err)
		}
		if articleMap[k].table, err = fetchTable(articleMap[k].id); err != nil {
			return nil, fmt.Errorf("failed to fetch table: %w", err)
		}
	}
	return
}

func articleToString(art *article) (text string, err error) {
	num := len(art.paragraphs)
	if art.table != nil {
		num++
	}
	ps := make([]string, num)
	for _, para := range art.paragraphs {
		html := "<p"
		if len(para.css) > 0 {
			html += " class=\"" + para.css + "\""
		}
		ps[para.ordinal] = html + ">" + para.text + "</p>"
	}
	if art.table != nil {
		var sb strings.Builder
		if art.table.title[:5] == "Table" {
			sb.WriteString("<a class=\"d-none\" id=\"")
			sb.WriteString(strings.Replace(strings.ToLower(art.table.title), " ", "-", -1))
			sb.WriteString("\">")
			sb.WriteString(art.table.title)
			sb.WriteString("</a><p class=\"d-none\">")
			for i := range art.table.rows {
				if i != 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(art.table.rows[i][0])
			}
		} else {
			sb.WriteString("<p class=\"d-none\">")
			sb.WriteString(art.table.title)
		}
		sb.WriteString("</p><table class=\"table table-striped mb-4\"><thead><tr>")
		for _, col := range art.table.header {
			sb.WriteString("<th>")
			sb.WriteString(col)
			sb.WriteString("</th>")
		}
		sb.WriteString("</tr></thead><tbody>")
		for _, row := range art.table.rows {
			sb.WriteString("<tr>")
			for _, col := range row {
				sb.WriteString("<td>")
				sb.WriteString(col)
				sb.WriteString("</td>")
			}
			sb.WriteString("</tr>")
		}
		sb.WriteString("</tbody></table>")
		ps[art.table.ordinal] = sb.String()
	}
	for _, p := range ps {
		text += p
	}
	return
}

type articleDepth struct {
	art   *article
	depth int
}

func resolveDepths(art *article, depth int) (d []articleDepth) {
	d = []articleDepth{{art, depth}}
	if art.sub != nil {
		for _, sub := range art.sub {
			d = append(d, resolveDepths(sub, depth+1)...)
		}
	}
	return d
}

func div(class string) string {
	return "<div class=\"" + class + "\">"
}

func articleToArticle(art *article) (article models.Article, err error) {
	depths := resolveDepths(art, 0)
	var sb strings.Builder
	sb.WriteString(div("row"))
	pDepth := 0
	for _, depth := range depths {
		if pDepth < depth.depth {
			sb.WriteString(div("row px-4 mt-2"))
		} else if pDepth > depth.depth {
			sb.WriteString("</div>")
		}
		pDepth = depth.depth
		sb.WriteString(div("col col-md-12 col-sm-12"))
		sb.WriteString("<h")
		sb.WriteRune(rune(depth.depth + 0x31))
		sb.WriteString(" id=\"")
		sb.WriteString(strings.ToLower(strings.Replace(depth.art.title, " ", "-", -1)))
		sb.WriteString("\">")
		sb.WriteString(depth.art.title)
		sb.WriteString("</h")
		sb.WriteRune(rune(depth.depth + 0x31))
		sb.WriteRune('>')
		var s string
		if s, err = articleToString(depth.art); err != nil {
			return article, fmt.Errorf("failed to parse article: %w", err)
		}
		sb.WriteString(s)
		sb.WriteString("</div>")
	}
	for ; pDepth > 0; pDepth-- {
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")
	article = models.Article{
		Title: art.title,
		Text:  template.HTML(sb.String()),
		Table: "article",
	}
	return
}

func FetchArticle(title string) (result models.Article, err error) {
	var art *article
	if art, err = fetchArticle(title); err != nil {
		return
	}
	return articleToArticle(art)
}

func FetchArticles() (articles []models.Article, err error) {
	var arts []*article
	if arts, err = fetchArticles(); err != nil {
		return
	}
	articles = make([]models.Article, len(arts))
	for i, art := range arts {
		if articles[i], err = articleToArticle(art); err != nil {
			return
		}
	}
	return
}

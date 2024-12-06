package data

import (
	"NeonManager/models"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
var articleQuery = `
SELECT id, title, txt, tags, tbl FROM (
SELECT a.id as id, a.title as title, a.document as txt, COALESCE(GROUP_CONCAT(t.tag, ','), '') as tags, 'article' as tbl
FROM article a 
    LEFT JOIN _articles_tags t ON a.id = t.article
GROUP BY a.id
UNION SELECT -1 as id, r.name || '(' || r.cost || ' GP)' as title, r.description as txt, e.name || 'race' as tags, 'race' as tbl
      FROM race r JOIN extension e ON e.id = r.extension
UNION SELECT -1 as id, b.name || '(' || b.cost || ' GP)' as title, b.description as txt, e.name || 'background,bg' as tags, 'background' as tbl
      FROM background b JOIN extension e ON e.id = b.extension
UNION SELECT -1 as id, aff.name as title, aff.description as txt, e.name || 'boon' as tags, 'affinity' as tbl
      FROM affinity aff JOIN extension e ON e.id = aff.extension WHERE aff.isBoon = 1
UNION SELECT -1 as id, aff.name as title, aff.description as txt, e.name || 'bane' as tags, 'affinity' as tbl
      FROM affinity aff JOIN extension e ON e.id = aff.extension WHERE aff.isBoon = 0
UNION SELECT -1 as id, a.name || '(' || a.cost || ' XP)' as title, a.effect as txt, e.name || 'ability' as tags, 'ability' as tbl
      FROM ability a JOIN extension e ON e.id = a.extension
UNION SELECT -1 as id, s.name || '(' || CASE WHEN s.cost == 0 THEN 'advanced' ELSE 'basic' END || ')' as title, s.description as txt, e.name || 'skill' as tags, 'skill' as tbl
      FROM skill s JOIN extension e ON e.id = s.extension
) WHERE title IS NOT NULL
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
			return fmt.Errorf("item with pk \"%s\" does not exists (query: \"%s\"): NULL", (*item).Pk(), query)
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

func scanArticle(rows *sql.Rows, item *models.Article) (err error) {
	var tagsStr string
	if err = rows.Scan(&(*item).Id, &(*item).Title, &(*item).Text, &tagsStr, &(*item).Table); err == nil {
		(*item).Tags = strings.Split(tagsStr, ",")
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

func FetchArticles() (articles []models.Article, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(articleQuery); err != nil {
		return nil, fmt.Errorf("failed to prepare fetchArticles query: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	return fetchItems(stmt, nil, scanArticle)
}

func FetchArticle(title string) (article models.Article, err error) {
	var stmt *sql.Stmt
	if stmt, err = db.Prepare(articleQuery + " AND title LIKE ?"); err != nil {
		return article, fmt.Errorf("failed to prepare fetchArticle query: %w", err)
	}
	var rows *sql.Rows
	if rows, err = stmt.Query(title); err != nil {
		return article, fmt.Errorf("failed to fetch article: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		err = scanArticle(rows, &article)
	} else {
		err = fmt.Errorf("article \"%s\" not found", title)
	}
	return
}

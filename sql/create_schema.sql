CREATE TABLE IF NOT EXISTS stats
(
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    cr  INTEGER NOT NULL DEFAULT 0,
    int INTEGER NOT NULL DEFAULT 0,
    ins INTEGER NOT NULL DEFAULT 0,
    ch  INTEGER NOT NULL DEFAULT 0,
    ag  INTEGER NOT NULL DEFAULT 0,
    dex INTEGER NOT NULL DEFAULT 0,
    con INTEGER NOT NULL DEFAULT 0,
    str INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS extension
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(30)
);

CREATE INDEX IF NOT EXISTS extension_name_IDX ON extension (name);

BEGIN TRANSACTION;

/* ensure 'core rules' is id 0 */
INSERT OR REPLACE INTO extension(id, name)
VALUES (0, 'core rules');

/* ensure sequence did not get reset */
UPDATE sqlite_sequence
SET seq = (SELECT MAX(id) FROM extension)
WHERE name = 'extension';

COMMIT;

CREATE TABLE IF NOT EXISTS race
(
    name        VARCHAR(30) PRIMARY KEY,
    cost        INTEGER NOT NULL,
    extension   INTEGER NOT NULL DEFAULT 0,
    description TEXT    NOT NULL,
    stats       INTEGER NOT NULL,
    FOREIGN KEY (stats) REFERENCES stats (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (extension) REFERENCES extension (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS background
(
    name        VARCHAR(30) PRIMARY KEY,
    cost        INTEGER NOT NULL,
    extension   INTEGER NOT NULL DEFAULT 0,
    description TEXT    NOT NULL,
    stats       INTEGER NOT NULL,
    FOREIGN KEY (stats) REFERENCES stats (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (extension) REFERENCES extension (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS affinity
(
    name        VARCHAR(30) PRIMARY KEY,
    cost        INTEGER NOT NULL,
    extension   INTEGER NOT NULL DEFAULT 0,
    description TEXT    NOT NULL,
    isBoon      BOOLEAN NOT NULL,
    FOREIGN KEY (extension) REFERENCES extension (id) ON DELETE RESTRICT ON UPDATE RESTRICT
    /* TODO mods */
);

CREATE INDEX IF NOT EXISTS affinity_boon_IDX ON affinity (isBoon);

CREATE TABLE IF NOT EXISTS ability
(
    name      VARCHAR(30) PRIMARY KEY,
    cost      INTEGER NOT NULL,
    extension INTEGER NOT NULL DEFAULT 0,
    effect    TEXT    NOT NULL,
    FOREIGN KEY (extension) REFERENCES extension (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS skill
(
    name        VARCHAR(30) PRIMARY KEY,
    cost        INTEGER NOT NULL,
    extension   INTEGER NOT NULL DEFAULT 0,
    stat        TEXT    NOT NULL,
    description TEXT    NOT NULL,
    FOREIGN KEY (extension) REFERENCES extension (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

/* race join tables */
CREATE TABLE IF NOT EXISTS _races_backgrounds
(
    race       TEXT    NOT NULL,
    background TEXT    NOT NULL,
    allowed    BOOLEAN NOT NULL,
    PRIMARY KEY (race, background),
    FOREIGN KEY (race) REFERENCES race (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (race) REFERENCES background (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _races_affinities
(
    race     TEXT    NOT NULL,
    affinity TEXT    NOT NULL,
    level    INTEGER NOT NULL DEFAULT 0,
    grp      INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount   INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (race, affinity),
    FOREIGN KEY (race) REFERENCES race (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (affinity) REFERENCES affinity (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _races_abilities
(
    race    TEXT    NOT NULL,
    ability TEXT    NOT NULL,
    grp     INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount  INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (race, ability),
    FOREIGN KEY (race) REFERENCES race (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (ability) REFERENCES ability (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _races_skills
(
    race   TEXT    NOT NULL,
    skill  TEXT    NOT NULL,
    level  INTEGER NOT NULL DEFAULT 0,
    grp    INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (race, skill),
    FOREIGN KEY (race) REFERENCES race (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (skill) REFERENCES skill (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS race_bg_race_IDX ON _races_backgrounds (race, allowed);
CREATE INDEX IF NOT EXISTS race_bg_bg_IDX ON _races_backgrounds (background, allowed);

CREATE INDEX IF NOT EXISTS race_af_race_IDX ON _races_affinities (race);
CREATE INDEX IF NOT EXISTS race_af_af_IDX ON _races_affinities (affinity);

CREATE INDEX IF NOT EXISTS race_ab_race_IDX ON _races_abilities (race);
CREATE INDEX IF NOT EXISTS race_ab_ab_IDX ON _races_abilities (ability);

CREATE INDEX IF NOT EXISTS race_sk_race_IDX ON _races_skills (race);
CREATE INDEX IF NOT EXISTS race_sk_sk_IDX ON _races_skills (skill);

/* background join tables */
CREATE TABLE IF NOT EXISTS _backgrounds_affinities
(
    background TEXT    NOT NULL,
    affinity   TEXT    NOT NULL,
    level      INTEGER NOT NULL DEFAULT 0,
    grp        INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount     INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (background, affinity),
    FOREIGN KEY (background) REFERENCES background (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (affinity) REFERENCES affinity (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _backgrounds_abilities
(
    background TEXT    NOT NULL,
    ability    TEXT    NOT NULL,
    grp        INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount     INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (background, ability),
    FOREIGN KEY (background) REFERENCES background (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (ability) REFERENCES ability (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _backgrounds_skills
(
    background TEXT    NOT NULL,
    skill      TEXT    NOT NULL,
    level      INTEGER NOT NULL DEFAULT 0,
    grp        INTEGER NOT NULL DEFAULT 0, /* groupings make optional if grp > 0 */
    gcount     INTEGER NOT NULL DEFAULT 0, /* number of items that may be chosen from a group */
    PRIMARY KEY (background, skill),
    FOREIGN KEY (background) REFERENCES background (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (skill) REFERENCES skill (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS bg_af_bg_IDX ON _backgrounds_affinities (background);
CREATE INDEX IF NOT EXISTS bg_af_af_IDX ON _backgrounds_affinities (affinity);

CREATE INDEX IF NOT EXISTS bg_ab_bg_IDX ON _backgrounds_abilities (background);
CREATE INDEX IF NOT EXISTS bg_ab_ab_IDX ON _backgrounds_abilities (ability);

CREATE INDEX IF NOT EXISTS bg_sk_bg_IDX ON _backgrounds_skills (background);
CREATE INDEX IF NOT EXISTS bg_sk_sk_IDX ON _backgrounds_skills (skill);

/* ability requirements */
CREATE TABLE IF NOT EXISTS _abilities_requires
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    ability   VARCHAR(30) NOT NULL,
    reqType   INTEGER     NOT NULL DEFAULT 0, /* ability: 0, affinity: 1, skill: 2 */
    rAffinity VARCHAR(30),
    rAbility  VARCHAR(30),
    rSkill    VARCHAR(30),
    FOREIGN KEY (ability) REFERENCES ability (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (rAffinity) REFERENCES affinity (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (rAbility) REFERENCES ability (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (rSkill) REFERENCES skill (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS ab_ab_ab_IDX ON _abilities_requires (ability);

CREATE TABLE IF NOT EXISTS character
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(255),
    gp          INTEGER NOT NULL DEFAULT 100,
    xp          INTEGER NOT NULL DEFAULT 2500,
    stats       INTEGER NOT NULL,
    race        VARCHAR(30), /* directly resolve affinities/abilities/skills do not take from race (=> choice groups) */
    background  VARCHAR(30), /* same as race comment */
    description TEXT,
    FOREIGN KEY (stats) REFERENCES stats (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (race) REFERENCES race (name) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (background) REFERENCES background (name) ON DELETE CASCADE ON UPDATE CASCADE
);

/* extensions used/selected for the character; no entries for core */
CREATE TABLE IF NOT EXISTS _characters_extensions
(
    character INTEGER,
    extension INTEGER,
    PRIMARY KEY (character, extension),
    FOREIGN KEY (character) REFERENCES character(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (extension) REFERENCES extension(id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE IF NOT EXISTS _characters_affinities
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT, /* affinities allows duplicate entries with different mods */
    character INTEGER NOT NULL,
    affinity  VARCHAR(30),
    level     INTEGER NOT NULL DEFAULT 0,
    /* TODO mods */
    FOREIGN KEY (character) REFERENCES character (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (affinity) REFERENCES affinity (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _characters_abilities
(
    character INTEGER,
    ability   VARCHAR(30),
    PRIMARY KEY (character, ability),
    FOREIGN KEY (character) REFERENCES character (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (ability) REFERENCES ability (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS _characters_skills
(
    character INTEGER,
    skill     VARCHAR(30),
    level     INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (character, skill),
    FOREIGN KEY (character) REFERENCES character (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (skill) REFERENCES skill (name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS article (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent INTEGER REFERENCES article(id),
    title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS paragraph (
    article INTEGER NOT NULL REFERENCES article(id),
    ordinal INTEGER NOT NULL DEFAULT 0,
    text TEXT NOT NULL,
    css VARCHAR(255) DEFAULT '',
    PRIMARY KEY (article, ordinal)
);

CREATE TABLE IF NOT EXISTS ptable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article INTEGER NOT NULL UNIQUE REFERENCES article(id),
    ordinal INTEGER NOT NULL,
    title VARCHAR(255) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS table_col (
    id INTEGER REFERENCES ptable(id),
    col INTEGER NOT NULL,
    row INTEGER NOT NULL,
    value VARCHAR(255),
    PRIMARY KEY (id, col, row)
);
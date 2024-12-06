package web

import (
	"NeonManager/data"
	"NeonManager/logger"
	"NeonManager/models"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type HttpError struct {
	Message string
	Status  int
}

func (e HttpError) Error() string {
	return e.Message
}

func Serve() error {

	for addr, handler := range createHandlers() {
		http.HandleFunc(addr, handler)
	}

	// TODO admin panel -> drop data / import + file upload
	//TODO templates
	// error
	// character creation
	// race, background, affinity, ability, skill
	//
	logger.Info("Listening on :8080")
	return http.ListenAndServe(":8080", nil)
}

// handleError accepts a message and a status to create an error response
func handleError(w http.ResponseWriter, err HttpError) {
	// TODO handle gracefully
	if err.Status < 400 || err.Status >= 600 {
		err.Status = http.StatusInternalServerError
	}
	tmpl, content, er := withBase("./templates/error.html")
	if er != nil {
		http.Error(w, err.Message, err.Status)
		return
	}
	content.Title = http.StatusText(err.Status)
	content.Content = models.ErrorSite{Message: err.Message, Status: err.Status}
	if er = tmpl.Execute(w, content); er != nil {
		http.Error(w, err.Message, err.Status)
	}
}

// handleErr accepts an error to create a http 500 response
func handleErr(w http.ResponseWriter, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(strings.ToLower(msg), "not found"):
		handleError(w, HttpError{Message: msg, Status: http.StatusNotFound})
	case strings.Contains(strings.ToLower(msg), "missing") || strings.Contains(strings.ToLower(msg), "invalid"):
		handleError(w, HttpError{Message: msg, Status: http.StatusBadRequest})
	default:
		handleError(w, HttpError{Message: msg})
	}
}

func createHandlers() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateIndex(r); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/favicon.ico"] = func(w http.ResponseWriter, r *http.Request) {
		ico := "./templates/favicon.ico"
		file, err := os.Open(ico)
		if err != nil {
			handleErr(w, err)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		w.Header().Set("Content-Type", "image/png")
		if _, err = io.Copy(w, file); err != nil {
			handleErr(w, err)
		}
	}
	handlers["/foundation"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Foundation"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/dice-checks-and-stats"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Dice Checks and Stats"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/movement"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Movement"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/combat"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Combat"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/healing"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Healing"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/hazards"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Hazards"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/transhumanism"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Transhumanism"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/hacking"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Hacking"); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/character-creation"] = func(w http.ResponseWriter, r *http.Request) {
		article := r.URL.Query().Get("type")
		if tmpl, d, err := createTemplateArticleList(article); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/create"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateCharacterCreation(r); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/equipment"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateEquipment(r); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/addons"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateAddons(r); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/glossary"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateGlossary(r); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	handlers["/search"] = func(w http.ResponseWriter, r *http.Request) {
		table := strings.ToLower(r.URL.Query().Get("t"))
		query := r.URL.Query().Get("q")
		if table == "" || query == "" {
			handleErr(w, fmt.Errorf("missing query parameter 't'/'q'"))
			return
		}
		var tmpl *template.Template
		var d *models.BaseSite
		var err error
		switch table {
		case "article":
			tmpl, d, err = createTemplateArticle(query)
		case "race", "races":
			tmpl, d, err = createTemplateDbItem(query, data.FetchRace, "./templates/race.html")
		case "background", "backgrounds":
			tmpl, d, err = createTemplateDbItem(query, data.FetchBackground, "./templates/background.html")
		case "affinity", "boons", "banes":
			tmpl, d, err = createTemplateDbItem(query, data.FetchAffinity, "./templates/affinity.html")
		case "ability", "abilities":
			tmpl, d, err = createTemplateDbItem(query, data.FetchAbility, "./templates/ability.html")
		case "skill", "skills":
			tmpl, d, err = createTemplateDbItem(query, data.FetchSkill, "./templates/skill.html")
		default:
			handleErr(w, fmt.Errorf("invalid table %s", table))
			return
		}
		if err != nil {
			switch err.Error() {
			case "missing query parameter":
				handleErr(w, fmt.Errorf("missing query parameter 'q'"))
				return
			case "not found":
				handleErr(w, fmt.Errorf("\"%s\" not found", query))
			default:
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
			return
		}
		if err = tmpl.Execute(w, d); err != nil {
			logger.Error("failed to parse template: %v", err)
			handleErr(w, err)
		}
	}
	return handlers
}

// withBase builds a template using the paths given as templates
func withBase(templates ...string) (*template.Template, *models.BaseSite, error) {
	templates = append([]string{"./templates/base.html"}, templates...)
	t, err := template.ParseFiles(templates...)
	return t, &models.BaseSite{Title: "Neon Manager"}, err
}

func createTemplateIndex(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	if r.URL.Path != "/" {
		return t, d, fmt.Errorf("%s not found", r.URL.Path)
	}
	return withBase("./templates/index.html")
}

func createTemplateArticle(title string) (t *template.Template, d *models.BaseSite, err error) {
	var article models.Article
	if article, err = data.FetchArticle(title); err != nil {
		return
	}
	t, d, err = withBase("./templates/article.html")
	d.Title = article.Title
	d.Content = article
	return
}

func createTemplateArticleList(articleType string) (t *template.Template, d *models.BaseSite, err error) {
	var article models.ArticleList
	switch strings.ToLower(articleType) {
	case "":
		return createTemplateArticle("Character Creation")
	case "races":
		article.Article.Title = "Races"
		if article.List, err = data.FetchRaces(false); err != nil {
			return
		}
	case "bgs", "backgrounds":
		article.Article.Title = "Backgrounds"
		if article.List, err = data.FetchBackgrounds(false); err != nil {
			return
		}
	case "boons":
		article.Article.Title = "Boons"
		boon := true
		if article.List, err = data.FetchAffinities(&boon); err != nil {
			return
		}
	case "banes":
		article.Article.Title = "Banes"
		boon := false
		if article.List, err = data.FetchAffinities(&boon); err != nil {
			return
		}
	case "affinities":
		article.Article.Title = "Affinities"
		if article.List, err = data.FetchAffinities(nil); err != nil {
			return
		}
	case "abilities":
		article.Article.Title = "Abilities"
		if article.List, err = data.FetchAbilities(); err != nil {
			return
		}
	case "skills":
		article.Article.Title = "Skills"
		if article.List, err = data.FetchSkills(); err != nil {
			return
		}
	default:
		return t, d, fmt.Errorf("invalid article: %s", articleType)
	}
	if article.Article, err = data.FetchArticle(article.Article.Title); err != nil {
		return t, d, fmt.Errorf("failed to fetch article: %v", err)
	}
	t, d, err = withBase("./templates/article_list.html")
	d.Content = article
	d.Title = article.Article.Title
	return
}

func createTemplateCharacterCreation(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	return withBase("./templates/character-creation.html")
}

func createTemplateEquipment(r *http.Request) (*template.Template, *models.BaseSite, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func createTemplateAddons(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	return withBase("./templates/addons.html")
}

func createTemplateGlossary(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	var articles []models.Article
	if articles, err = data.FetchArticles(); err != nil {
		return
	}
	t, d, err = withBase("./templates/glossary.html")
	d.Title = "Glossary"
	d.Content = articles
	return
}

func createTemplateDbItem[T models.DbItem](name string, fetch func(string) (T, error), template string) (t *template.Template, d *models.BaseSite, err error) {
	name = strings.Trim(strings.Split(name, "(")[0], " \t\n")
	if name == "" {
		return nil, nil, fmt.Errorf("missing query parameter")
	}
	var item T
	if item, err = fetch(name); err != nil {
		if strings.HasSuffix(err.Error(), "NULL") {
			return nil, nil, fmt.Errorf("not found")
		}
		return nil, nil, fmt.Errorf("failed to fetch item from db: %v", err)
	}
	t, d, err = withBase(template)
	d.Title = item.Pk()
	d.Content = item
	return
}

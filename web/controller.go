package web

import (
	"NeonManager/data"
	"NeonManager/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func Serve() error {

	for addr, handler := range createHandlers() {
		http.HandleFunc(addr, handler)
	}

	// TODO admin panel -> drop data / import + file upload
	//TODO templates
	// error
	// race, background, affinity, ability, skill

	return http.ListenAndServe(":8080", nil)
}

// handleError accepts a message and a status to create an error response
func handleError(w http.ResponseWriter, message string, status int) {
	// TODO handle gracefully
	if status < 400 || status >= 600 {
		status = http.StatusInternalServerError
	}
	http.Error(w, message, status)
}

// handleErr accepts an error to create a http 500 response
func handleErr(w http.ResponseWriter, err error) {
	handleError(w, err.Error(), http.StatusInternalServerError)
}

func createHandlers() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, d, _ := createTemplateIndex()
		if err := tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/foundation"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Foundation"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/dice-checks-and-stats"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Dice Checks and Stats"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/movement"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Movement"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/combat"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Combat"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/healing"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Healing"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/hazards"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Hazards"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/transhumanism"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Transhumanism"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/hacking"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateArticle("Hacking"); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/character-creation"] = func(w http.ResponseWriter, r *http.Request) {
		var article string
		articleType := r.URL.Query().Get("type")
		switch strings.ToLower(articleType) {
		case "":
			article = "Character Creation"
		case "race":
			article = "Race"
		case "bg", "background":
			article = "Background"
		case "boons":
			article = "Boons"
		case "banes":
			article = "Banes"
		case "abilities":
			article = "Abilities"
		case "skills":
			article = "Skills"
		}
		if tmpl, d, err := createTemplateArticle(article); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/create"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateCharacterCreation(r); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/equipment"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateEquipment(r); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/addons"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateAddons(r); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/glossary"] = func(w http.ResponseWriter, r *http.Request) {
		if tmpl, d, err := createTemplateGlossary(r); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		} else if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	handlers["/search"] = func(w http.ResponseWriter, r *http.Request) {
		table := r.URL.Query().Get("t")
		query := r.URL.Query().Get("q")
		if table == "" || query == "" {
			handleError(w, "Bad request: missing query parameter 't'/'q'", http.StatusBadRequest)
			return
		}
		var tmpl *template.Template
		var d *models.BaseSite
		var err error
		switch table {
		case "article":
			tmpl, d, err = createTemplateArticle(query)
		case "race":
			tmpl, d, err = createTemplateDbItem(query, data.FetchRace, "./templates/race.html")
		case "background":
			tmpl, d, err = createTemplateDbItem(query, data.FetchBackground, "./templates/background.html")
		case "affinity":
			tmpl, d, err = createTemplateDbItem(query, data.FetchAffinity, "./templates/affinity.html")
		case "ability":
			tmpl, d, err = createTemplateDbItem(query, data.FetchAbility, "./templates/ability.html")
		case "skill":
			tmpl, d, err = createTemplateDbItem(query, data.FetchSkill, "./templates/skill.html")
		}
		if err != nil {
			switch err.Error() {
			case "missing query parameter":
				handleError(w, "Ba request: missing query parameter 'q'", http.StatusBadRequest)
				return
			case "not found":
				handleError(w, "Not found: "+query, http.StatusNotFound)
			}
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
			return
		}
		if err = tmpl.Execute(w, d); err != nil {
			log.Printf("failed to parse template %v", err)
			handleErr(w, err)
		}
	}
	return handlers
}

// ff reads the content of a file or returns "" on err
func ff(path string) string {
	if bytes, err := os.ReadFile(path); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

// base returns the base site's template which requires for which "content" fills the template
func base() (*template.Template, *models.BaseSite) {
	return template.Must(template.New("base").Parse(ff("./templates/base.html"))),
		&models.BaseSite{Title: "Neon Manager"}
}

func createTemplateIndex() (t *template.Template, d *models.BaseSite, err error) {
	t, d = base()
	t, err = t.New("content").Parse(ff("./templates/index.html"))
	return
}

func createTemplateArticle(title string) (t *template.Template, d *models.BaseSite, err error) {
	var article models.Article
	if article, err = data.FetchArticle(title); err != nil {
		return
	}
	t, d = base()
	d.Title = article.Title
	d.Content = article
	t, err = t.New("content").Parse(ff("./templates/article.html"))
	return

}

func createTemplateCharacterCreation(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	t, d = base()
	t, err = t.New("content").Parse(ff("./templates/character-creation.html"))
	return
}

func createTemplateEquipment(r *http.Request) (*template.Template, *models.BaseSite, error) {
	return nil, nil, nil
}

func createTemplateAddons(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	t, d = base()
	t, err = t.New("content").Parse(ff("./templates/addons.html"))
	return
}

func createTemplateGlossary(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	var articles []models.Article
	if articles, err = data.FetchArticles(); err != nil {
		return
	}
	t, d = base()
	d.Title = "Glossary"
	d.Content = articles
	t, err = t.New("content").Parse(ff("./templates/glossary.html"))
	return
}

func pruneName(name string) string {
	return strings.Trim(strings.Split(name, "(")[0], " \t\n")
}

func createTemplateDbItem[T models.DbItem](name string, fetch func(string) (T, error), template string) (t *template.Template, d *models.BaseSite, err error) {
	name = pruneName(name)
	if name == "" {
		return nil, nil, fmt.Errorf("missing query parameter")
	}
	var item T
	if item, err = fetch(name); err != nil {
		if strings.HasPrefix(err.Error(), "NULL") {
			return nil, nil, fmt.Errorf("not found")
		}
		return nil, nil, fmt.Errorf("failed to fetch item from db: %v", err)
	}
	t, d = base()
	d.Title = item.Pk()
	d.Content = item
	t, err = t.New("content").Parse(ff(template))
	return
}

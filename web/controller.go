package web

import (
	"NeonManager/data"
	"NeonManager/logger"
	"NeonManager/models"
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type HttpError struct {
	Message string
	Status  int
}

func (e HttpError) Error() string {
	return e.Message
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

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexHandler struct {
	routes []*route
}

func (h *RegexHandler) Handler(patter *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{patter, handler})
}

func (h *RegexHandler) HandlerFunc(patter *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{patter, http.HandlerFunc(handler)})
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	handleError(w, HttpError{Message: r.URL.String(), Status: http.StatusNotFound})
}

func Serve() error {
	handler := &RegexHandler{routes: make([]*route, 0)}

	if regex, err := regexp.Compile("^/$"); err != nil {
		logger.Error("Failed to compile regex for index: ", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := withBase(); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	sendImg := func(file *os.File, out http.ResponseWriter) {
		defer func() {
			_ = file.Close()
		}()

		out.Header().Set("Content-Type", "image/png")
		if _, err := io.Copy(out, file); err != nil {
			handleErr(out, err)
		}
	}
	if regex, err := regexp.Compile("^/favicon\\.ico$"); err != nil {
		logger.Error("Failed to compile regex for favicon: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			ico := "./media/favicon.ico"
			file, err := os.Open(ico)
			if err != nil {
				handleErr(w, err)
				return
			}
			sendImg(file, w)
		})
	}
	if regex, err := regexp.Compile("^/media/.*$"); err != nil {
		logger.Error("Failed to compile regex for media: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			segments := strings.Split(r.URL.String(), "media/")
			if len(segments) != 2 || strings.Contains(segments[1], "..") {
				handleErr(w, fmt.Errorf("invalid media url: %s", r.URL.String()))
			}
			file, err := os.Open(path.Join("./media", segments[1]))
			if err != nil {
				handleErr(w, err)
			}
			sendImg(file, w)
		})
	}
	if regex, err := regexp.Compile("^/the-foundation$"); err != nil {
		logger.Error("Failed to compile regex for foundation: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("The Foundation"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/dice-checks-and-stats$"); err != nil {
		logger.Error("Failed to compile regex for dice-checks-and-stats: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Dice Checks and Stats"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/movement$"); err != nil {
		logger.Error("Failed to compile regex for movement: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Movement"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/combat$"); err != nil {
		logger.Error("Failed to compile regex for combat: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Combat"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/healing$"); err != nil {
		logger.Error("Failed to compile regex for healing: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Healing"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/hazards$"); err != nil {
		logger.Error("Failed to compile regex for hazards: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Hazards"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/transhumanism$"); err != nil {
		logger.Error("Failed to compile regex for transhumanism: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Transhumanism"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/hacking$"); err != nil {
		logger.Error("Failed to compile regex for hacking: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateArticle("Hacking"); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/character-creation(/.*)?$"); err != nil {
		logger.Error("Failed to compile regex for character-creation: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			var article string
			if segments := strings.Split(r.URL.String(), "character-creation/"); len(segments) == 2 {
				article = segments[1]
			}
			var templateFunc func(string) (*template.Template, *models.BaseSite, error)
			if typeAndName := strings.Split(article, "/"); len(typeAndName) > 1 {
				templateFunc = createTemplateItem
			} else {
				templateFunc = createTemplateArticleList
			}
			if tmpl, d, err := templateFunc(article); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/create(/.*)?$"); err != nil {
		logger.Error("Failed to compile regex for create: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateCharacterCreation(r); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/equipment(/.*)?$"); err != nil {
		logger.Error("Failed to compile regex for equipment: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateEquipment(r); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/addons(/.*)?$"); err != nil {
		logger.Error("Failed to compile regex for healing: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateAddons(r); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}
	if regex, err := regexp.Compile("^/glossary$"); err != nil {
		logger.Error("Failed to compile regex for glossary: %v", err)
	} else {
		handler.HandlerFunc(regex, func(w http.ResponseWriter, r *http.Request) {
			if tmpl, d, err := createTemplateGlossary(r); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			} else if err = tmpl.Execute(w, d); err != nil {
				logger.Error("failed to parse template: %v", err)
				handleErr(w, err)
			}
		})
	}

	logger.Info("Listening on :8080")
	return http.ListenAndServe(":8080", handler)
}

// withBase builds a template using the paths given as templates
func withBase(templates ...string) (*template.Template, *models.BaseSite, error) {
	templates = append([]string{"./templates/base.html"}, templates...)
	t, err := template.ParseFiles(templates...)
	return t, &models.BaseSite{Title: "Neon Manager"}, err
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
		article.Article.Title = "Race & Background"
		if article.List, err = data.FetchRaces(false); err != nil {
			return
		}
	case "bgs", "backgrounds":
		article.Article.Title = "Race & Background"
		if article.List, err = data.FetchBackgrounds(false); err != nil {
			return
		}
	case "boons":
		article.Article.Title = "Boons & Banes"
		boon := true
		if article.List, err = data.FetchAffinities(&boon); err != nil {
			return
		}
	case "banes":
		article.Article.Title = "Boons & Banes"
		boon := false
		if article.List, err = data.FetchAffinities(&boon); err != nil {
			return
		}
	case "affinities":
		article.Article.Title = "Boons & Banes"
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
	switch strings.ToLower(articleType) {
	case "races":
		article.Article.Title = "Races"
	case "bgs", "backgrounds":
		article.Article.Title = "Backgrounds"
	case "boons":
		article.Article.Title = "Boons"
	case "banes":
		article.Article.Title = "Banes"
	}
	t, d, err = withBase("./templates/article_list.html")
	d.Content = article
	d.Title = article.Article.Title
	return
}

func createTemplateRace(name string) (t *template.Template, d *models.BaseSite, err error) {
	return withBase() // TODO
}

func createTemplateBg(name string) (t *template.Template, d *models.BaseSite, err error) {
	return withBase() // TODO
}

func createTemplateAffinity(topic string, name string) (t *template.Template, d *models.BaseSite, err error) {
	return withBase() // TODO
}

func createTemplateSkill(name string) (t *template.Template, d *models.BaseSite, err error) {
	return withBase() // TODO
}

func createTemplateAbility(name string) (t *template.Template, d *models.BaseSite, err error) {
	return withBase() // TODO
}

func createTemplateItem(article string) (t *template.Template, d *models.BaseSite, err error) {
	segments := strings.Split(article, "/")
	topic, item := segments[0], segments[1]
	switch topic {
	case "races":
		return createTemplateRace(item)
	case "bgs", "backgrounds":
		return createTemplateBg(item)
	case "boons":
		return createTemplateAffinity("Boon", item)
	case "banes":
		return createTemplateAffinity("Bane", item)
	case "affinities":
		return createTemplateAffinity("Affinity", item)
	case "abilities":
		return createTemplateAbility(item)
	case "skills":
		return createTemplateSkill(item)
	default:
		return t, d, fmt.Errorf("topic %s not found", topic)
	}
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

type glEntry struct {
	Title string
	Text  string
	Tags  []string
	Link  string
}

func evalAnd(entry glEntry, cond string) bool {
	conditions := strings.Split(cond, "&")
	for _, condition := range conditions {
		if len(condition) < 2 {
			return false
		}
		switch condition[0] {
		case '#':
			contained := false
			for _, tag := range entry.Tags {
				if condition[1:] == strings.ToLower(tag) {
					contained = true
				}
			}
			if !contained {
				return false
			}
		case '"':
			if !(condition[len(condition)-1] == '"' &&
				strings.Contains(strings.ToLower(entry.Text), condition[1:len(condition)-1]) ||
				strings.Contains(strings.ToLower(entry.Text), condition[1:])) {
				return false
			}
		default:
			if !strings.Contains(strings.ToLower(entry.Title), condition) {
				return false
			}
		}

	}
	return true
}

func matchFilter(entry glEntry, filter string) bool {
	if filter == "" {
		return true
	}
	filter = strings.ToLower(filter)
	conditions := strings.Split(filter, "|")
	for _, condition := range conditions {
		if evalAnd(entry, condition) {
			return true
		}
	}
	return false
}

func getId(node *html.Node) string {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "id" {
				return attr.Val
			}
		}
	}
	return ""
}

func getText(node *html.Node) (text string) {
	if node == nil {
		return ""
	}
	if node.Type == html.TextNode {
		text = node.Data
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return
}

func getLink(title string) string {
	switch title {
	case "Races", "Backgrounds", "Boons", "Banes", "Skills", "Abilities":
		return "/character-creation?type=" + strings.ToLower(title)
	default:
		return "/" + strings.Replace(strings.ToLower(title), " ", "-", -1)
	}
}

func trimTx(text string, max int) string {
	if len(text) > max {
		return text[:max]
	}
	return text
}

func convArticleToEntries(article models.Article, maxLength int, filter string) ([]glEntry, error) {
	doc, err := html.Parse(strings.NewReader(string(article.Text)))
	if err != nil {
		return nil, err
	}
	entries := make([]glEntry, 0)
	link := getLink(article.Title)
	var text strings.Builder
	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.TextNode {
			if text.Len() < maxLength {
				text.WriteString(n.Data)
			}
		} else if id := getId(n); id != "" {
			var pText string
			for c := n.NextSibling; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "p" {
					pText = getText(c)
					break
				}
			}
			entry := glEntry{
				Title: getText(n),
				Text:  pText,
				Tags:  append([]string{article.Title}, article.Tags...),
				Link:  link + "#" + id,
			}
			if matchFilter(entry, filter) {
				if len(entry.Text) > maxLength {
					entry.Text = entry.Text[:maxLength]
				}
				entries = append(entries, entry)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}
	parseNode(doc)
	if len(filter) == 0 {
		if len(entries) > 0 && len(entries[0].Text) == 0 {
			entries[0].Tags = append([]string{"Topic"}, entries[0].Tags...)
			if text.Len() <= maxLength {
				entries[0].Text = text.String()
			} else {
				entries[0].Text = text.String()[:maxLength]
			}
		}
	} else {
		topic := glEntry{
			Title: article.Title,
			Text:  text.String(),
			Tags:  append([]string{"Topic", article.Title}, article.Tags...),
			Link:  link,
		}
		if matchFilter(topic, filter) {
			if len(topic.Text) > maxLength {
				topic.Text = topic.Text[:maxLength]
			}
			entries = append([]glEntry{topic}, entries...)
		}
	}
	return entries, nil
}

func createTemplateGlossary(r *http.Request) (t *template.Template, d *models.BaseSite, err error) {
	filter := r.URL.Query().Get("q")
	var articles []models.Article
	glEntries := make([]glEntry, 0)
	if articles, err = data.FetchArticles(); err != nil {
		return
	}
	for _, article := range articles {
		if items, err1 := convArticleToEntries(article, 300, filter); err1 != nil {
			logger.Warn("Failed to convert article to entries: %v", err1)
		} else {
			if items != nil {
				glEntries = append(glEntries, items...)
			}
		}
	}
	t, d, err = withBase("./templates/glossary.html")
	d.Title = "Glossary"
	d.Content = glEntries
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

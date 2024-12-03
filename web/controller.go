package web

import (
	"html/template"
	"log"
	"net/http"
)

func Serve() error {

	for addr, handler := range createHandlers() {
		http.HandleFunc(addr, handler)
	}

	// TODO admin panel -> drop data / import + file upload

	return http.ListenAndServe(":8080", nil)
}

func createHandlers() map[string]func(http.ResponseWriter, *http.Request) {
	handlers := make(map[string]func(http.ResponseWriter, *http.Request))
	handlers["/"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl := createTemplateIndex()
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/foundation"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateFoundation()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/dice-checks-and-stats"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateDiceChecksAndStats()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/movement"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateMovement()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/combat"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateCombat()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/healing"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateHealing()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/hazards"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateHazards()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/transhumanism"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateTranshumanism()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/hacking"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateHacking()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/character-creation"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateCharacterCreation(r)
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/equipment"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateEquipment()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/addons"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateAddons()
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	handlers["/glossary"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl, data := createTemplateGlossary(r)
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("failed to parse template %v", err)
		}
	}
	return handlers
}

func createTemplateIndex() *template.Template {
	return template.Must(template.ParseFiles("index.html"))
}

func createTemplateFoundation() (*template.Template, any) {
	return nil, nil
}

func createTemplateDiceChecksAndStats() (*template.Template, any) {
	return nil, nil
}

func createTemplateMovement() (*template.Template, any) {
	return nil, nil
}

func createTemplateCombat() (*template.Template, any) {
	return nil, nil
}

func createTemplateHealing() (*template.Template, any) {
	return nil, nil
}

func createTemplateHazards() (*template.Template, any) {
	return nil, nil
}

func createTemplateTranshumanism() (*template.Template, any) {
	return nil, nil
}

func createTemplateHacking() (*template.Template, any) {
	return nil, nil
}

func createTemplateCharacterCreation(r *http.Request) (*template.Template, any) {
	return template.Must(template.ParseFiles("character-creation.html")), nil
}

func createTemplateEquipment() (*template.Template, any) {
	return nil, nil
}

func createTemplateAddons() (*template.Template, any) {
	return template.Must(template.ParseFiles("addons.html")), nil
}

func createTemplateGlossary(r *http.Request) (*template.Template, any) {
	return nil, nil
}

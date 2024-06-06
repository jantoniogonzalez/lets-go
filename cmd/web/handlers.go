package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jantoniogonzalez/lets-go/internal/models"
	"github.com/jantoniogonzalez/lets-go/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil || id < 1 {
		app.notFound(w)
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	var createForm snippetCreateForm

	err := app.decodePostForm(r, &createForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	createForm.CheckField(validator.NotEmpty(createForm.Title), "title", "This field cannot be blank")
	createForm.CheckField(validator.MaxChars(createForm.Title, 100), "title", "This field cannot be more than 100 characters long")
	createForm.CheckField(validator.NotEmpty(createForm.Content), "content", "This field cannot be blank")
	createForm.CheckField(validator.PermittedInt(createForm.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !createForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = createForm
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(createForm.Title, createForm.Content, createForm.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

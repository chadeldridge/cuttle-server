package htmx

//
// HTMX enabled components.
//

const html_tab = "\t"

type Component interface {
	// render returns the html string for the component.
	// string: Element id override.
	// Attributes: Shared Attributes.
	render(string, Attributes) string
}

// type Renderer func(next Component) Component
type Renderer func(string, Attributes) string

func (c Renderer) render(id string, shared Attributes) string {
	return c(id, shared)
}

//
// Rendering
//

// Render renders a component chain witht he given string and Attributes.
func Render(c Component, id string, shared Attributes) string {
	if shared == nil {
		shared = Attributes{}
	}

	if _, ok := shared["indent"]; !ok {
		shared["indent"] = "0"
	}

	return c.render(id, shared)
}

//
// Components return a ComponentFunc that can render html.
//

func NewComponent(c Component, compID string, compAttrs Attributes) Component {
	return Renderer(func(id string, attrs Attributes) string {
		return c.render(compID, compAttrs)
	})
}

func HeaderComponent(title string, attrs Attributes, comps ...Component) Component {
	return Renderer(func(id string, shared Attributes) string {
		open := `<html>
  <head>
    <script src="https://unpkg.com/htmx.org@2.0.1"></script>
    <title>` + title + `</title>
`
		var elems string
		for _, c := range comps {
			elems += shared.ComponentWrapper(c, id)
		}

		return shared.LnWrapper(open+elems) + `  </head>`
	})
}

func FormComponent(method, target, hxTarget, hxSwap string, attrs Attributes, comps ...Component) Component {
	return Renderer(func(id string, shared Attributes) string {
		if id == "" {
			id = string(attrs["id"])
		}

		elem := `<form` +
			attrHXMethod(method, target) +
			attrHXSwap(attrs["hx-swap"]) +
			attrClass(attrs["class"]) +
			attrID(id) +
			`>`
		if len(comps) == 0 {
			return elem + `</form>`
		}

		var elems string
		for _, c := range comps {
			elems += shared.ComponentWrapper(c, id)
		}

		return elem + shared.LnWrapper(elems) + `</form>`
	})
}

// DivComponent combines all Component arguments sequentially inside of a single div element.
func DivComponent(attrs Attributes, c ...Component) Component {
	return Renderer(func(id string, shared Attributes) string {
		if id == "" {
			id = attrs["id"]
		}

		elem := `<div` +
			attrHXTarget(attrs["hx-target"]) +
			attrHXSwap(attrs["hx-swap"]) +
			attrClass(shared["class"]) +
			attrID(id) +
			`>`
		if len(c) == 0 {
			return elem + `</div>`
		}

		var comps string
		for _, next := range c {
			if attrs["inline"] == "true" {
				comps += next.render(id, shared)
				continue
			}

			comps += shared.ComponentWrapper(next, id)
		}

		return elem + shared.LnWrapper(comps) + `</div>`
	})
}

func InlineDivComponent(attrs Attributes, comps ...Component) Component {
	attrs["inline"] = "true"
	return DivComponent(attrs, comps...)
}

// InputComponent renders an input element. Optionally wraps the input with a label.
func InputComponent(inputID string, inputType string, attrs Attributes, label Component) Component {
	return Renderer(func(id string, shared Attributes) string {
		var input, l string
		if label != nil {
			l = shared.LnWrapper(label.render(id, shared))
		}

		input = `<input type="` + inputType + `"` +
			attrID(inputID) +
			attrClass(attrs["class"]) +
			attrPlaceholder(attrs["placeholder"]) +
			attrValue(attrs["value"]) +
			attrOnChange(attrs["onChange"]) +
			`>`

		return l + input
	})
}

//
// Elements return a ComponentFunc that can render an individual html element.
//

// LabelElement wraps the next component with a label.
// func LabelElement(next Component) Component {
func LabelElement(id, labelFor string, label string, attrs Attributes) Renderer {
	return func(id string, shared Attributes) string {
		return `<label` +
			attrID(id) +
			attrClass(attrs["class"]) +
			attrFor(labelFor) +
			`>` + label + `</label>`
	}
}

func ButtonElement(text string, attrs Attributes, optDiv Component) Renderer {
	return func(id string, shared Attributes) string {
		if id == "" {
			id = string(attrs["id"])
		}

		return `<button` +
			attrID(id) +
			attrClass(attrs["class"]) +
			`>` + text + `</button>`
	}
}

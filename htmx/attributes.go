package htmx

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Attributes map[string]string

func (a Attributes) AddIndent() {
	indent, _ := strconv.Atoi(a["indent"])
	a["indent"] = strconv.Itoa(indent + 1)
}

func (a Attributes) RemoveIndent() {
	indent, _ := strconv.Atoi(a["indent"])
	a["indent"] = strconv.Itoa(indent - 1)
}

//
// Wrappers always return "" if the string is empty.
//

func (a Attributes) ComponentWrapper(c Component, id string) string {
	if c == nil {
		return ""
	}

	// i, _ := strconv.Atoi(a["indent"])
	// d := fmt.Sprintf("\n%s%s\n", strings.Repeat(html_tab, i), c.render(id, a))
	a.AddIndent()
	d := "\n" + a.IndentWrapper(c.render(id, a))
	a.RemoveIndent()

	return d
	// return "\n" + s + "\n"
}

func (a Attributes) IndentWrapper(s string) string {
	if s == "" {
		return ""
	}

	var data string
	indent, err := strconv.Atoi(a["indent"])
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(s, "\n") {
		data += fmt.Sprintf("%s%s\n", strings.Repeat(html_tab, indent), line)
	}
	/*
		indent, _ := strconv.Atoi(a["indent"])
		scanner := bufio.NewScanner(strings.NewReader(s))
		for scanner.Scan() {
			data += fmt.Sprintf("%s%s", strings.Repeat(html_tab, indent), s)
		}
	*/
	return strings.TrimRight(data, "\n")
	// return fmt.Sprintf("%s%s", strings.Repeat(html_tab, indent), s)
}

func (a Attributes) PrefixLnWrapper(s string) string {
	if s == "" {
		return ""
	}

	return "\n" + s
	// return fmt.Sprintf("\n%s", s)
}

func (a Attributes) LnWrapper(s string) string {
	if s == "" {
		return ""
	}

	return s + "\n"
	// return fmt.Sprintln(s)
}

//
// Additive formaters always return a non-empty string.
//

func (a Attributes) Indent() string {
	i, _ := strconv.Atoi(a["indent"])
	return strings.Repeat(html_tab, i)
}

func (a Attributes) NewLine() string {
	return "\n"
}

//
// attr functions return an attribute as a string. Always returns "" if param is empty.
//

func attrHXTarget(target string) string {
	if target == "" {
		return ""
	}

	return ` hx-target="` + target + `"`
}

func attrHXSwap(swap string) string {
	if swap == "" {
		return ""
	}

	return ` hx-swap="` + swap + `"`
}

func attrHXMethod(method, resource string) string {
	switch method {
	case "GET":
		return attrHXGet(resource)
	case "POST":
		return attrHXPost(resource)
	case "PUT":
		return attrHXPut(resource)
	default:
		return ""
	}
}

func attrHXGet(resource string) string {
	if resource == "" {
		return ""
	}

	return ` hx-get="` + resource + `"`
}

func attrHXPost(resource string) string {
	if resource == "" {
		return ""
	}

	return ` hx-post="` + resource + `"`
}

func attrHXPut(resource string) string {
	if resource == "" {
		return ""
	}

	return ` hx-put="` + resource + `"`
}

func attrID(id string) string {
	if id == "" {
		return ""
	}

	return ` id="` + string(id) + `"`
}

func attrClass(class string) string {
	if class == "" {
		return ""
	}

	return ` class="` + class + `"`
}

func attrFor(id string) string {
	if id == "" {
		return ""
	}

	return ` for="` + string(id) + `"`
}

func attrValue(value string) string {
	if value == "" {
		return ""
	}

	return ` value="` + value + `"`
}

func attrPlaceholder(placeholder string) string {
	if placeholder == "" {
		return ""
	}

	return ` placeholder="` + placeholder + `"`
}

func attrOnChange(onChange string) string {
	if onChange == "" {
		return ""
	}

	return ` onchange="` + onChange + `"`
}

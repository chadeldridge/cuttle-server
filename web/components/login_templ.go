// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Login(redirect string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"app\" class=\"\"><form hx-post=\"/login.html\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Card(
			CardTitle("Cuttle Login"),
			LoginFields(redirect),
			CardActions(
				ButtonSubmit("3", "Login", "3", nil),
				ButtonOutlined("4", "Sign Up", "4", templ.Attributes{"onclick": "location.href='/signup.html'"}),
			),
		).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</form><output id=\"errors\" class=\"CardErrors\"></output></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func LoginFields(redirect string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if redirect == "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("redirect = \"/\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"block py-5\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InputOutlined("username", "text", "Username",
			templ.Attributes{
				"placeholder": "Username",
				"tabindex":    "1",
				"onkeypress":  OnKeyFocusID(13, "2"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InputOutlined("password", "password", "Password",
			templ.Attributes{
				"placeholder": "Password",
				"tabindex":    "2",
				"onkeypress":  OnKeyFocusID(13, "3"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input id=\"redirect\" type=\"hidden\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(redirect)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `login.templ`, Line: 36, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func Signup() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"app\" class=\"\"><form hx-post=\"/signup.html\" hx-target=\"errors\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Card(
			CardTitle("Cuttle Login"),
			CardContent(SignupFields()),
			CardActions(
				ButtonSubmit("signupButton", "Sign Up", "5", nil),
				ButtonOutlined("6", "Login", "6", templ.Attributes{"onclick": "location.href='/login.html'"}),
			),
		).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</form></div><script>\n\t\tlet pwMatch = false;\n\t\tlet name = document.getElementById(\"name\");\n\t\tlet username = document.getElementById(\"username\");\n\t\tlet signupButton = document.getElementById(\"signupButton\");\n\t\tlet password = document.getElementById(\"password\");\n\t\tlet confirmPassword = document.getElementById(\"confirmPassword\");\n\t\tlet confirmPasswordDiv = document.getElementById(\"confirmPassword-div\");\n\t\tlet confirmPasswordLabel = document.getElementById(\"confirmPassword-label\");\n\t\tlet errorOutput = document.getElementById(\"errors\");\n\n\t\tif (password.value == \"\" || password.value !== confirmPassword.value) {\n\t\t\tsignupButton.disabled = true;\n\t\t}\n\t\n\t\tlet checkPasswords = () => {\n\t\t\tif (password.value === \"\" || password.value !== confirmPassword.value) {\n\t\t\t\tpwMatch = false;\n\t\t\t\tsignupButton.disabled = true;\n\t\t\t\terrorOutput.style.display = \"block\";\n\t\t\t\terrorOutput.innerHTML = \"Passwords do not match.\";\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"border-primary-light\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"hover:border-primary-highlight\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"has-[:focus]:border-secondary-base\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"peer-has-[:hover]:text-primary-highlight\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"peer-has-[:focus]:text-secondary-base\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"border-success\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"text-success\");\n\t\t\t\tconfirmPasswordDiv.classList.add(\"border-error\");\n\t\t\t\tconfirmPasswordLabel.classList.add(\"text-error\");\n\t\t\t} else {\n\t\t\t\tpwMatch = true;\n\t\t\t\tsignupButton.disabled = false;\n\t\t\t\terrorOutput.style.display = \"none\";\n\t\t\t\terrorOutput.innerHTML = \"\";\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"border-primary-light\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"hover:border-primary-highlight\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"has-[:focus]:border-secondary-base\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"peer-has-[:hover]:text-primary-highlight\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"peer-has-[:focus]:text-secondary-base\");\n\t\t\t\tconfirmPasswordLabel.classList.remove(\"text-error\");\n\t\t\t\tconfirmPasswordDiv.classList.remove(\"border-error\");\n\t\t\t\tconfirmPasswordDiv.classList.add(\"border-success\");\n\t\t\t\tconfirmPasswordLabel.classList.add(\"text-success\");\n\t\t\t}\n\n\t\t\tenableButton();\n\t\t};\n\n\t\tname.addEventListener('keyup', enableButton);\n\t\tusername.addEventListener('keyup', enableButton);\n\t\tpassword.addEventListener('keyup', checkPasswords);\n\t\tconfirmPassword.addEventListener('keyup', checkPasswords);\n\n\t\tfunction enableButton() {\n\t\t\tif (allFilled()) {\n\t\t\t\tsignupButton.disabled = false;\n\t\t\t} else {\n\t\t\t\tsignupButton.disabled = true;\n\t\t\t}\n\t\t}\n\n\t\tfunction allFilled() {\n\t\t\treturn document.getElementById(\"name\").value !== \"\" &&\n\t\t\t\tdocument.getElementById(\"username\").value !== \"\" &&\n\t\t\t\tdocument.getElementById(\"password\").value !== \"\" &&\n\t\t\t\tdocument.getElementById(\"confirmPassword\").value !== \"\" &&\n\t\t\t\tpwMatch;\n\t\t}\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func SignupFields() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = InputOutlined("name", "text", "Name",
			templ.Attributes{
				"placeholder": "Bob",
				"tabindex":    "1",
				"onkeypress":  OnKeyFocusID(13, "2"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InputOutlined("username", "text", "Username",
			templ.Attributes{
				"placeholder": "myusername",
				"tabindex":    "2",
				"onkeypress":  OnKeyFocusID(13, "3"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InputOutlined("password", "password", "Password",
			templ.Attributes{
				"placeholder": "Make It A Strong One",
				"tabindex":    "3",
				"onkeypress":  OnKeyFocusID(13, "4"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = InputOutlined("confirmPassword", "password", "Confirm Password",
			templ.Attributes{
				"placeholder": "Retype Password",
				"tabindex":    "4",
				"onkeypress":  OnKeyFocusID(13, "5"),
			}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"errors\" class=\"display-hidden content-center relative mt-6 min-h-0 flex-auto flex flex-col text-text-disabled border border-primary-base border-opacity-75\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

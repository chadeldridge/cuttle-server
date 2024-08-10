package components

var theme ColorTheme

func init() {
	theme = DefaultTheme()
}

type ColorTheme struct {
	Primary   Colors
	Secondary Colors
}

type Colors struct {
	Color Color
	Light Color
	Dark  Color
}

type Color struct {
	Color string
	Text  string
}

func DefaultTheme() ColorTheme {
	return ColorTheme{
		Primary: Colors{
			Color: Color{
				Color: "gray-800",
				Text:  "gray-300",
			},
			Light: Color{
				Color: "gray-700",
				Text:  "gray-300",
			},
			Dark: Color{
				Color: "gray-900",
				Text:  "gray-300",
			},
		},
		Secondary: Colors{
			Color: Color{
				Color: "blue-500",
				Text:  "gray-300",
			},
			Light: Color{
				Color: "blue-300",
				Text:  "gray-900",
			},
			Dark: Color{
				Color: "blue-700",
				Text:  "gray-300",
			},
		},
	}
}

func SetTheme(theme ColorTheme) {
	theme = theme
}

func GetTheme() ColorTheme {
	return theme
}

func bgPrimary() string {
	return "bg-" + theme.Primary.Color.Color
}

func bgPrimaryLight() string {
	return "bg-" + theme.Primary.Light.Color
}

func bgPrimaryDark() string {
	return "bg-" + theme.Primary.Dark.Color
}

func bgSecondary() string {
	return "bg-" + theme.Secondary.Color.Color
}

func bgSecondaryLight() string {
	return "bg-" + theme.Secondary.Light.Color
}

func bgSecondaryDark() string {
	return "bg-" + theme.Secondary.Dark.Color
}

func textPrimary() string {
	return "text-" + theme.Primary.Color.Text
}

func textPrimaryLight() string {
	return "text-" + theme.Primary.Light.Text
}

func textPrimaryDark() string {
	return "text-" + theme.Primary.Dark.Text
}

func textSecondary() string {
	return "text-" + theme.Secondary.Color.Text
}

func textSecondaryLight() string {
	return "text-" + theme.Secondary.Light.Text
}

func textSecondaryDark() string {
	return "text-" + theme.Secondary.Dark.Text
}

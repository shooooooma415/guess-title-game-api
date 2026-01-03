package theme

// Theme represents a theme entity
type Theme struct {
	id    ThemeID
	title ThemeTitle
	hint  string
}

// NewTheme creates a new Theme
func NewTheme(id ThemeID, title ThemeTitle, hint string) *Theme {
	return &Theme{
		id:    id,
		title: title,
		hint:  hint,
	}
}

// Getters
func (t *Theme) ID() ThemeID {
	return t.id
}

func (t *Theme) Title() ThemeTitle {
	return t.title
}

func (t *Theme) Hint() string {
	return t.hint
}

// UpdateHint updates the theme hint
func (t *Theme) UpdateHint(hint string) {
	t.hint = hint
}

package ui

import (
	"errors"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	nameField = iota
	portField
	passwordField
	databaseTypeField
	detachedField
	actionField
)

type createForm struct {
	inputs       []textinput.Model
	active       int
	dbTypeOption int
	detach       int
	action       int
	cancelled    bool
	initialCmd   tea.Cmd
}

func newCreateForm() createForm {
	inputs := []textinput.Model{
		newFormTextInput("my-database", false),
		newFormTextInput("5432", false),
		newFormTextInput("", true),
	}
	initialCmd := inputs[nameField].Focus()
	return createForm{inputs: inputs, initialCmd: initialCmd}
}

func newFormTextInput(placeholder string, password bool) textinput.Model {
	input := textinput.New()
	input.Prompt = "› "
	input.Placeholder = placeholder
	input.SetWidth(42)
	if password {
		input.EchoMode = textinput.EchoPassword
	}

	styles := input.Styles()
	accent := lipgloss.Color("69")
	styles.Focused.Prompt = lipgloss.NewStyle().Foreground(accent).Bold(true)
	styles.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	styles.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	input.SetStyles(styles)
	return input
}

func (m createForm) run() (createForm, error) {
	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return createForm{}, err
	}
	form, ok := result.(createForm)
	if !ok {
		return createForm{}, errors.New("unexpected form result")
	}
	return form, nil
}

func (m createForm) Init() tea.Cmd { return m.initialCmd }

func (m createForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if m.active == actionField {
				if m.action == 1 {
					m.cancelled = true
				}
				return m, tea.Quit
			}
			m.active++
			return m, m.focusActiveField()
		case "left", "h":
			if m.active >= databaseTypeField {
				m.choosePrevious()
				return m, nil
			}
		case "right", "l":
			if m.active >= databaseTypeField {
				m.chooseNext()
				return m, nil
			}
		}
	}

	if m.active >= databaseTypeField {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
	return m, cmd
}

func (m *createForm) focusActiveField() tea.Cmd {
	for index := range m.inputs {
		m.inputs[index].Blur()
	}
	if m.active < databaseTypeField {
		return m.inputs[m.active].Focus()
	}
	return nil
}

func (m *createForm) choosePrevious() {
	switch m.active {
	case databaseTypeField:
		m.dbTypeOption = (m.dbTypeOption + 1) % 2
	case detachedField:
		m.detach = (m.detach + 1) % 2
	case actionField:
		m.action = (m.action + 1) % 2
	}
}

func (m *createForm) chooseNext() { m.choosePrevious() }

func (m createForm) View() tea.View {
	rows := []string{
		m.fieldView(nameField, "Database name", m.name()),
		m.fieldView(portField, "Port", m.port()),
		m.fieldView(passwordField, "Password", m.password()),
		m.choiceView(databaseTypeField, "Database type", m.dbType()),
		m.choiceView(detachedField, "Run container in detached mode?", map[bool]string{true: "Yes", false: "No"}[m.detached()]),
	}
	if m.active == actionField {
		rows = append(rows, m.actionsView())
	}
	return tea.NewView(strings.Join(rows, "\n\n"))
}

func (m createForm) fieldView(index int, label, value string) string {
	if m.active == index {
		return labelStyle.Render(label) + "\n" + m.inputs[index].View()
	}
	if m.active > index {
		return labelStyle.Render(label+":") + " " + valueStyle.Render(value)
	}
	return mutedStyle.Render(label)
}

func (m createForm) choiceView(index int, label, value string) string {
	if m.active == index {
		return labelStyle.Render(label) + "\n" + promptStyle.Render("› ") + valueStyle.Render(value) + mutedStyle.Render("  ←/→")
	}
	if m.active > index {
		return labelStyle.Render(label+":") + " " + valueStyle.Render(value)
	}
	return mutedStyle.Render(label)
}

func (m createForm) actionsView() string {
	submit := buttonStyle.Render("Submit")
	cancel := buttonStyle.Render("Cancel")
	if m.action == 0 {
		submit = selectedButtonStyle.Render("Submit")
	} else {
		cancel = cancelButtonStyle.Render("Cancel")
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, submit, "  ", cancel)
	return buttons + "\n" + mutedStyle.Render("←/→ choose • Enter confirm")
}

func (m createForm) name() string { return valueOrPlaceholder(m.inputs[nameField]) }
func (m createForm) port() string { return valueOrPlaceholder(m.inputs[portField]) }
func (m createForm) password() string {
	if value := m.inputs[passwordField].Value(); value != "" {
		return strings.Repeat("•", len([]rune(value)))
	}
	return ""
}
func (m createForm) dbType() string {
	return []string{"PostgreSQL", "MySQL"}[m.dbTypeOption]
}
func (m createForm) detached() bool { return m.detach == 0 }

func valueOrPlaceholder(input textinput.Model) string {
	if value := input.Value(); value != "" {
		return value
	}
	return input.Placeholder
}

var (
	labelStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)
	valueStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	promptStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
	mutedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	buttonStyle         = lipgloss.NewStyle().Width(10).Height(1).Align(lipgloss.Center).Foreground(lipgloss.Color("252")).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("241"))
	selectedButtonStyle = buttonStyle.Copy().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("69")).BorderForeground(lipgloss.Color("69")).Bold(true)
	cancelButtonStyle   = buttonStyle.Copy().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("204")).BorderForeground(lipgloss.Color("204")).Bold(true)
)

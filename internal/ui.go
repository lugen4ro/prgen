package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI styles using semantic color naming
var (
	// Semantic color palette - keeping original colors but remapping usage for better contrast
	primaryColor    = lipgloss.Color("#F7C59F") // Peach - Primary highlights
	secondaryColor  = lipgloss.Color("#087E8B") // Teal - Secondary actions, info
	textColor       = lipgloss.Color("#E1E5EE") // Alice Blue - Standard text (slightly gray-ish white)
	neutralColor    = lipgloss.Color("#C7CCDB") // French Gray - Borders, muted text (keeping for compatibility)
	backgroundColor = lipgloss.Color("#0B3954") // Prussian Blue - Background (for light text)
	errorColor      = lipgloss.Color("#FF6B6B") // Warm red - Errors

	titleStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Margin(0, 0)

	configHeaderStyle = lipgloss.NewStyle().
				Foreground(neutralColor).
				Bold(false).
				Margin(0, 0)

	successStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	warningStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1)

	configTableStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(neutralColor).
				Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(neutralColor).
			Padding(0, 1).
			Margin(0, 0)
)

// InitializeUI sets up Bubble Tea for beautiful output
func InitializeUI() {
	// Bubble Tea setup is handled per operation
}

// ShowStartupBanner displays a welcome banner
func ShowStartupBanner() {
	banner := lipgloss.NewStyle().
		Foreground(textColor).
		Background(backgroundColor).
		Bold(true).
		Padding(1, 3).
		Render("🚀 PR Generator")

	centered := lipgloss.NewStyle().
		Width(80).
		Align(lipgloss.Center).
		Render(banner)
	fmt.Println(centered)
	fmt.Println()
}

// ShowConfigSummary displays essential configuration information in a compact format
func ShowConfigSummary(config *Config) {
	fmt.Println(configHeaderStyle.Render("Configuration Summary"))

	configData := []string{
		fmt.Sprintf("Config Directory: %s", config.ConfigDir),
		fmt.Sprintf("LLM Provider: %v", config.MainConfig["llm_provider"]),
		fmt.Sprintf("Model: %v", config.MainConfig["model"]),
	}

	content := strings.Join(configData, "\n")
	fmt.Println(configTableStyle.Render(content))
}

// ShowDiffInfo displays git diff information
func ShowDiffInfo(diffLength int) {
	if diffLength == 0 {
		fmt.Println(warningStyle.Render("⚠️  No changes detected. Nothing to generate PR for."))
		return
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("ℹ️  Found changes (%d characters)", diffLength)))
}

// ShowGeneratedContent displays the generated title and body in styled panels
func ShowGeneratedContent(title, body string) {
	// Title label and panel
	titleLabel := lipgloss.NewStyle().
		Foreground(backgroundColor).
		Background(primaryColor).
		Bold(true).
		Padding(0, 1).
		Render("PR Title")
	titlePanel := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		MarginTop(-1).
		Render(title)
	titleSection := lipgloss.JoinVertical(lipgloss.Left, titleLabel, titlePanel)

	// Body label and panel
	bodyLabel := lipgloss.NewStyle().
		Foreground(textColor).
		Background(secondaryColor).
		Bold(true).
		Padding(0, 1).
		Render("PR Body")
	bodyPanel := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(0, 1).
		MarginTop(-1).
		Render(body)
	bodySection := lipgloss.JoinVertical(lipgloss.Left, bodyLabel, bodyPanel)

	// Stack sections vertically
	vertical := lipgloss.JoinVertical(lipgloss.Left, titleSection, bodySection)
	fmt.Println(vertical)
}

// SpinnerModel represents a Bubble Tea spinner
type SpinnerModel struct {
	spinner   spinner.Model
	message   string
	done      bool
	success   bool
	result    string
	taskError error
}

// NewSpinner creates a new spinner model
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(secondaryColor)

	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

// Init implements tea.Model
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update implements tea.Model
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case string:
		// Handle success/failure messages
		if msg == "success" {
			m.done = true
			m.success = true
		} else if msg == "fail" {
			m.done = true
			m.success = false
		}
		return m, tea.Quit
	case error:
		// Handle error messages
		m.done = true
		m.success = false
		m.taskError = msg
		return m, tea.Quit
	}
	return m, nil
}

// View implements tea.Model
func (m SpinnerModel) View() string {
	if m.done {
		if m.success {
			return successStyle.Render("✅ " + m.message + " - Done!")
		} else {
			return errorStyle.Render("❌ " + m.message + " - Failed!")
		}
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// RunSpinnerWithTask runs a spinner while executing a task
func RunSpinnerWithTask(message string, task func() error) error {
	model := NewSpinner(message)

	p := tea.NewProgram(model)

	// Run task in background
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to show spinner
		err := task()
		if err != nil {
			p.Send(err) // Send the actual error
		} else {
			p.Send("success")
		}
	}()

	finalModel, err := p.Run()
	if err != nil {
		return err // Bubble Tea error
	}

	// Return the task error if there was one
	if spinnerModel, ok := finalModel.(SpinnerModel); ok {
		return spinnerModel.taskError
	}

	return nil
}

// ShowProgress displays a simple progress message (fallback)
func ShowProgress(message string) {
	fmt.Println(infoStyle.Render("⏳ " + message))
}

// ShowSuccess displays a success message
func ShowSuccess(message string) {
	fmt.Println(successStyle.Render("✅ " + message))
}

// ShowError displays an error message with styling
func ShowError(message string, err error) {
	fmt.Println(errorStyle.Render(fmt.Sprintf("❌ %s: %v", message, err)))
}

// ShowPRSuccess displays successful PR creation with prominent URL
func ShowPRSuccess(prURL string) {
	fmt.Println()
	fmt.Println(successStyle.Render("🎉 Pull Request Created Successfully as Draft!"))

	urlHeader := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")). // Black text
		Background(neutralColor).
		Bold(true).
		Padding(0, 1).
		Render("PR URL")
	urlContent := fmt.Sprintf("%s\n%s", urlHeader, prURL)
	urlPanel := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(neutralColor).
		Padding(0, 1).
		Margin(0, 0).
		Render(urlContent)
	fmt.Println(urlPanel)
}

// AskConfirmation prompts the user for confirmation
func AskConfirmation(message string) bool {
	fmt.Print(warningStyle.Render("❓ " + message + " (Y/n): "))

	var response string
	fmt.Scanln(&response)

	// Default to yes if empty response
	if response == "" {
		return true
	}

	return strings.ToLower(response) != "n" && strings.ToLower(response) != "no"
}

// AskBackgroundInfo prompts the user to provide background information for PR generation
func AskBackgroundInfo() string {
	fmt.Println()
	fmt.Println(headerStyle.Render("📝 Background Information"))
	fmt.Println(infoStyle.Render("Provide context about your changes to help generate a better PR:"))
	fmt.Println(infoStyle.Render("• What problem does this solve?"))
	fmt.Println(infoStyle.Render("• What approach did you take?"))
	fmt.Println(infoStyle.Render("• Any important details or considerations?"))
	fmt.Println()
	fmt.Print(titleStyle.Render("Enter background info (press Ctrl+D when done): "))
	fmt.Println()

	var lines []string
	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		lines = append(lines, reader.Text())
	}

	if err := reader.Err(); err != nil {
		fmt.Println(errorStyle.Render("❌ Error reading input"))
		return ""
	}

	background := strings.Join(lines, "\n")

	if strings.TrimSpace(background) == "" {
		fmt.Println(infoStyle.Render("ℹ️  No background information provided"))
		return ""
	}

	fmt.Println(successStyle.Render("✅ Background information recorded"))
	return background
}

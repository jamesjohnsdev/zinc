package ui

// ConfirmDialog is a modal yes/no prompt shown before a destructive action.
type ConfirmDialog struct {
	active  bool
	message string
}

// NewConfirmDialog constructs a closed ConfirmDialog.
func NewConfirmDialog() ConfirmDialog {
	return ConfirmDialog{}
}

// Open shows the dialog with the given message.
func (c *ConfirmDialog) Open(message string) {
	c.active = true
	c.message = message
}

// Close hides the dialog.
func (c *ConfirmDialog) Close() {
	c.active = false
	c.message = ""
}

// Active reports whether the dialog is currently shown.
func (c ConfirmDialog) Active() bool {
	return c.active
}

// View renders the dialog.
func (c ConfirmDialog) View(width int) string {
	content := panelTitleStyle.Render("Confirm") + "\n" +
		c.message + "\n" +
		helpDescStyle.Render("y confirm   n/esc cancel")

	return panelFocusedStyle.Width(width).Render(content)
}

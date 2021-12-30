package browser

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

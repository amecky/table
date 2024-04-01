package term

type Model interface {
	View(m Model) string
}

func Run(m Model) {

}

package webapp

// DefaultHTMLTemplate объединяет все части шаблона
const DefaultHTMLTemplate = HTMLHead + HTMLBody + HelperFunctions + MainScript

// GetTemplate возвращает HTML шаблон с данными
func GetTemplate() string {
	return DefaultHTMLTemplate
}

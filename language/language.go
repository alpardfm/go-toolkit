package language

const (
	English    string = "en"
	Indonesian string = "id"
)

func HTTPStatusText(lang string, code int) string {
	var statusText string

	if lang != English && lang != Indonesian {
		lang = English
	}

	switch lang {
	case English:
		value, ok := statusTextEn[code]
		if ok {
			statusText = value
		} else {
			statusText = "Response Unknown"
		}
	case Indonesian:
		value, ok := statusTextId[code]
		if ok {
			statusText = value
		} else {
			statusText = "Tanggapan Tidak Diketahui"
		}
	}

	return statusText
}

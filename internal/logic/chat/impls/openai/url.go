package openai

type urls struct {
	BaseURL    string
	ModelList  string
	Completion string
}

func newURLs(baseURL string) *urls {
	return &urls{
		BaseURL:    baseURL,
		ModelList:  baseURL + "/models",
		Completion: baseURL + "/chat/completions",
	}
}

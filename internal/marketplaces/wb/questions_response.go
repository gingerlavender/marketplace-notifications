package wb

type QuestionsResponse struct {
	Data struct {
		Questions []Question `json:"questions"`
	} `json:"data"`
}

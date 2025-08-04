package wb

type FeedbacksResponse struct {
	Data struct {
		Feedbacks []Feedback `json:"feedbacks"`
	} `json:"data"`
}

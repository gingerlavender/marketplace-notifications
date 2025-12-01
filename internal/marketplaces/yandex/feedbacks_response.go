package yandex

type FeedbacksResponse struct {
	Result struct {
		Feedbacks []Feedback `json:"feedbacks"`
	} `json:"result"`
}

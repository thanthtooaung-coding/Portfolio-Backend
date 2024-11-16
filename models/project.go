package models

type Project struct {
	Title      string `json:"title"`
	GithubLink string `json:"githubLink"`
	ImgUrl     string `json:"imgUrl"`
	Category   string `json:"category"`
}

package api

import (
	"encoding/json"
	"io/ioutil"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
	"net/http"
)

const skyEngUrl = "https://dictionary.skyeng.ru/api/public/v1/"
const searchUrl = "words/search?search="

type TranslateInterface interface {
	GetTranslate(word string) (*models.Word, error)
}

type SkyEngApi struct {
	url string
}

func NewSkyEngApi() *SkyEngApi {
	return &SkyEngApi{
		url: skyEngUrl,
	}
}

func (s *SkyEngApi) GetTranslate(word string) (*models.Word, error) {
	url := s.url + searchUrl + word
	words := models.NewWordIngot()
	body, err := s.getRequest(url)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &words); err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, nil
	}
	result := &words[0]

	return result, nil
}

func (s *SkyEngApi) getRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

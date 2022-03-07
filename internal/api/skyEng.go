package api

import (
	"encoding/json"
	"io/ioutil"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
	"net/http"
	"strconv"
)

const skyEngUrl = "https://dictionary.skyeng.ru/api/public/v1/"
const searchUrl = "words/search?search="
const meaningUrl = "meanings?ids="

type TranslateInterface interface {
	GetTranslate(word string) (*models.Word, error)
	GetMeaning(word string) (*models.Mean, error)
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

func (s *SkyEngApi) GetMeaning(word string) (*models.Mean, error) {
	translate,err := s.GetTranslate(word)
	if err != nil {
		return nil,err
	}
	meanId := translate.Meanings[0].Id
	url := s.url + meaningUrl + strconv.Itoa(meanId)
	mean := models.NewMeanIngot()

	body, err := s.getRequest(url)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &mean); err != nil {
		return nil, err
	}
	res := &mean[0]

	return res,err
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
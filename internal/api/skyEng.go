package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const skyEngUrl = "https://dictionary.skyeng.ru/api/public/v1/"
const searchUrl = "words/search?search="
const meaningUrl = "meanings?ids="

var maxWordId = 99997
var errGeneratingWords = errors.New("generate words error")

type TranslateInterface interface {
	GetTranslate(word string) (*models.Word, error)
	GetMeaning(word string) (*models.Mean, error)
	GenerateRandomWords(countWords int) ([]*models.Mean, error)
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
	translate, err := s.GetTranslate(word)
	if err != nil {
		return nil, err
	}
	meanId := translate.Meanings[0].Id

	return s.getMean(meanId)
}

func (s *SkyEngApi) GenerateRandomWords(countWords int) ([]*models.Mean, error) {
	var model []*models.Mean
	for i := 0; i < countWords; i++ {
		rand.Seed(time.Now().UnixNano())
		id := rand.Intn(maxWordId)
		mean, err := s.getMean(id)
		if err != nil {
			log.Println(err)
			return model, errGeneratingWords
		}
		if mean != nil {
			model = append(model, mean)
		}
	}

	return model, nil
}

func (s *SkyEngApi) getMean(id int) (*models.Mean, error) {
	url := s.url + meaningUrl + strconv.Itoa(id)
	mean := models.NewMeanIngot()

	body, err := s.getRequest(url)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &mean); err != nil {
		return nil, err
	}
	if len(mean) == 0 {
		return nil, nil
	}

	return &mean[0], nil
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

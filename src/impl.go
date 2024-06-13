package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var httpClient = http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}}

type BlumDataImpl interface {
	UpdateUserData() error
	IncreasePoints(value int)
	SetGames(value int)
	GameAndClaim() error
	IsAviableGames() bool
	GetToken() string
}

func (b *BlumData) UpdateUserData() (err error) {
	if err = b.endpoints.ValidateUser(b.token); err != nil {
		return
	}
	err = b.endpoints.GetBalance(b)
	return
}

func (b *BlumData) IncreasePoints(value int) {
	b.totalPoints += value
	log.Println(fmt.Sprintf("баланс = [%d]", b.totalPoints))

}

func (b *BlumData) SetGames(value int) {
	b.totalGames = value
	log.Println(fmt.Sprintf("количество игр = [%d]", b.totalGames))
}

func (b *BlumData) GameAndClaim() (err error) {

	if err = b.endpoints.StartAndClaim(b); err != nil {
		return
	}
	log.Println(fmt.Sprintf("игра [%d] успешно выполнена. текущий баланс: [%d]", b.totalGames, b.totalPoints))

	b.totalGames--
	return
}

func (b *BlumData) IsAviableGames() bool {
	return b.totalGames > 0
}

func (b *BlumData) GetToken() string {
	return b.token
}

type EndpointImpl interface {
	ValidateUser(token string) error
	GetBalance(blum BlumDataImpl) error
	StartAndClaim(blum BlumDataImpl) error
}

func (b *Endpoints) ValidateUser(token string) (err error) {
	var resp *http.Response
	var buf = &bytes.Buffer{}

	req, _ := http.NewRequest("GET", b.Me, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = httpClient.Do(req)
	io.Copy(buf, resp.Body)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("- недействующий токен доступа.\n\tstatus: %d, \n\tmessage: %s, \n\terr: %v", resp.StatusCode, buf.String(), err)
		return
	}
	return
}

func (b *Endpoints) GetBalance(blum BlumDataImpl) (err error) {
	var resp *http.Response
	var buf = &bytes.Buffer{}

	req, _ := http.NewRequest("GET", b.Balance, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", blum.GetToken()))
	resp, err = httpClient.Do(req)
	io.Copy(buf, resp.Body)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("- не удалось получить данные аккаунта.\n\tstatus: %d, \n\tmessage: %s, \n\terr: %v", resp.StatusCode, buf.String(), err)
		return
	}
	var result = make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &result)
	balance, _ := strconv.ParseFloat(result["availableBalance"].(string), 64)
	games := int(result["playPasses"].(float64))
	blum.IncreasePoints(int(balance))
	blum.SetGames(games)
	return
}

func (b *Endpoints) StartAndClaim(blum BlumDataImpl) (err error) {
	var resp *http.Response
	var buf = &bytes.Buffer{}

	points := randRange(200, 250)

	req, _ := http.NewRequest("POST", b.StartGame, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", blum.GetToken()))
	resp, err = httpClient.Do(req)
	io.Copy(buf, resp.Body)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("- не удалось запустить игру.\n\tstatus: %d, \n\tmessage: %s, \n\terr: %v", resp.StatusCode, buf.String(), err)
		return
	}
	time.Sleep(time.Second * 30)

	var result = make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &result)
	result["points"] = points
	c, _ := json.Marshal(result)

	req, _ = http.NewRequest("POST", b.ClaimPoints, bytes.NewBuffer(c))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", blum.GetToken()))
	resp, err = httpClient.Do(req)
	buf.Reset()
	io.Copy(buf, resp.Body)
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		err = fmt.Errorf("- не удалось добавить очки.\n\tstatus: %d, \n\tmessage:%v, \n\terr :%v", resp.StatusCode, buf.String(), err)
		return
	}

	blum.IncreasePoints(points)
	return
}

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Accept-Encoding", ACCEPT_ENCODING)
	req.Header.Set("Accept-Language", ACCEPT_LANG)
	req.Header.Set("Accept", ACCEPT)
	req.Header.Set("Content-Type", CONTENT_TYPE)
	return t.underlyingTransport.RoundTrip(req)
}

func randRange(min, max int) int {
	rnd := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return rnd.Intn(max-min) + min
}

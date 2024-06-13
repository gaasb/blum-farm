package app

var (
	USER_AGENT      = "Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.113 Mobile Safari/537.36"
	ACCEPT_ENCODING = "gzip, deflate, br, zstd"
	ACCEPT_LANG     = "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"
	ACCEPT          = "application/json"
	CONTENT_TYPE    = "application/json"
)

type Endpoints struct {
	Me          string
	Balance     string
	StartGame   string
	ClaimPoints string
}

type BlumData struct {
	token     string
	endpoints EndpointImpl

	totalPoints int
	totalGames  int
}

func (b *BlumData) WithEndpoints(endpoints EndpointImpl) {
	b.endpoints = endpoints
}

func NewBlumData(token string) *BlumData {
	return &BlumData{token: token}
}

func NewEndpoints() *Endpoints {
	return &Endpoints{
		Me:          "https://gateway.blum.codes/v1/user/me",
		Balance:     "https://game-domain.blum.codes/api/v1/user/balance",
		StartGame:   "https://game-domain.blum.codes/api/v1/game/play",
		ClaimPoints: "https://game-domain.blum.codes/api/v1/game/claim",
	}
}

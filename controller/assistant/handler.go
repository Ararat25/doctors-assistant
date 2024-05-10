package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulrzcz/go-gigachat"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"webApp/model"
	"webApp/model/entity"
)

type Handler struct {
	authService *model.Service
}

func NewHandler(authService *model.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type reqBody struct {
	Symptoms string
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	jsonBody := reqBody{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &jsonBody)

	client, err := gigachat.NewInsecureClient(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenFound := entity.Apichats{}
	result := h.authService.Storage.Model(&entity.Apichats{}).Where(entity.Apichats{Name: "gigaChat"}).First(&tokenFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, "InternalServerError", http.StatusInternalServerError)
		return
	}

	if tokenFound.Token == "" {
		err = client.Auth()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result := h.authService.Storage.Model(&entity.Apichats{}).Where(entity.Apichats{Name: "gigaChat"}).Update("token", client.Token)
		if result.Error != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		result = h.authService.Storage.Model(&entity.Apichats{}).Where(entity.Apichats{Name: "gigaChat"}).Update("expiresAt", client.ExpiresAt)
		if result.Error != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if tokenFound.ExpiresAt > int(time.Now().Unix()) {
			err = client.Auth()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	fmt.Println(*client.ExpiresAt)

	msg := []gigachat.Message{
		{
			Role: gigachat.SystemRole,
			Content: "Ответ должен содержать 5 примерных диагнозов от самомого вероятного. А также небольшое описание каждого диагноза.\n" +
				"Перед ответом проанализируй суть запроса: если в запросе информация не связанная с медициной и в запросе нет симптомов, то возвращай такой ответ:" +
				"Неверный запрос, попробуйте ещё раз.\n",
		},
		{
			Role: gigachat.UserRole,
			Content: "Пример правильного запроса:\n" +
				"Симптомы: головная боль, боль в горле\n" +
				"Пример ответа на правильноый запрос:\n" +
				"1. Грипп: это вирусное заболевание, которое часто сопровождается высокой температурой, головной болью, болью в мышцах, сухим кашлем и болью в горле.\n2. Острый синусит: воспаление пазух носа, вызванное бактериальной инфекцией, может вызывать головную боль, заложенность носа, насморк и боль в горле.\n3. Ангина: инфекционное заболевание, вызываемое бактериями или вирусами, характеризующееся острой болью в горле, иногда сопровождающейся лихорадкой и головной болью.\n4. Менингит: воспаление оболочек мозга, обычно вызванное бактериальной или вирусной инфекцией. Симптомы включают сильную головную боль, ригидность затылочных мышц, светобоязнь и высокую температуру.\n5. Гайморит: воспаление верхнечелюстной пазухи, которое может вызывать головную боль, особенно при наклоне головы вниз, заложенность носа и боль в горле." +
				"Пример неправильного запроса:\n" +
				"Привет как дела\n" +
				"Пример ответа на неправильноый запрос:\n" +
				"Неверный запрос, попробуйте ещё раз!" +
				"Запрос:\n" +
				"Симптомы: " +
				jsonBody.Symptoms,
		},
	}

	temper := 0.87
	topP := 0.47
	var n int64
	n = 1
	stream := false
	var maxTokens int64
	maxTokens = 512
	repetitionPenalty := 1.07
	var updateInterval int64
	updateInterval = 0

	resp := &gigachat.ChatResponse{}

	resp, err = client.ChatWithContext(context.Background(), &gigachat.ChatRequest{
		Model:             gigachat.ModelLatest,
		Messages:          msg,
		Temperature:       &temper,
		TopP:              &topP,
		N:                 &n,
		Stream:            &stream,
		MaxTokens:         &maxTokens,
		RepetitionPenalty: &repetitionPenalty,
		UpdateInterval:    &updateInterval,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(resp)
	_, err = res.Write(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

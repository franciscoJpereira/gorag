package localnet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"ragAPI/pkg"
	"ragAPI/pkg/chat/store"

	"github.com/labstack/echo/v4"
)

type LocalControler struct {
	e   *echo.Echo
	rag *pkg.RAG
}

func NewLocalControler(rag *pkg.RAG) *LocalControler {
	return &LocalControler{
		e:   echo.New(),
		rag: rag,
	}
}

func EncodeBase64(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func DecodeBase64(name string) string {
	result, _ := base64.StdEncoding.DecodeString(name)
	return string(result)
}

// Wrapper to create a context and call the given method passing the corresponding data to it
func (lc *LocalControler) callMethod(method echo.HandlerFunc, url string, data any, params map[string]string) (response ResponseWriter, err error) {
	m, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(m))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return
	}

	response = *NewResponseWriter()
	c := lc.e.NewContext(req, &response)
	c.Set(pkg.RAGKey, lc.rag)
	names := make([]string, 0)
	values := make([]string, 0)
	for name, value := range params {
		names = append(names, name)
		values = append(values, value)
	}
	c.SetParamNames(names...)
	c.SetParamValues(values...)
	err = method(c)
	return
}

func (lc *LocalControler) GetAvailableKBs() (s []string, err error) {
	response, err := lc.callMethod(pkg.GetAvailableKBs, "/", []string{}, map[string]string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(response.Buf.Bytes(), &s)
	return
}

func (lc *LocalControler) CreateKB(kbname string) error {
	_, err := lc.callMethod(
		pkg.CreateKB,
		"/",
		[]string{},
		map[string]string{
			"KBName": kbname,
		},
	)
	return err
}

func (lc *LocalControler) AddDataToKB(data pkg.KBAddDataInstruct) error {
	_, err := lc.callMethod(
		pkg.AddDataToKB,
		"/",
		data,
		map[string]string{},
	)
	return err
}

func (lc *LocalControler) SingleShotMessage(data pkg.MessageInstruct) (response pkg.MessageResponse, err error) {
	rData, err := lc.callMethod(
		pkg.SingleShotMessage,
		"/",
		data,
		map[string]string{},
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(rData.Buf.Bytes(), &response)
	return
}

func (lc *LocalControler) SendNewMessageToChat(data pkg.ChatInstruct) (response pkg.MessageResponse, err error) {
	data.ChatName = EncodeBase64(data.ChatName)
	rData, err := lc.callMethod(
		pkg.SendNewMessageToChat,
		"/",
		data,
		map[string]string{},
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(rData.Buf.Bytes(), &response)
	return
}

func (lc *LocalControler) RetrieveAvailableChats() (chats []string, err error) {
	response, err := lc.callMethod(
		pkg.RetrieveAvailableChats,
		"/",
		[]string{},
		map[string]string{},
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(response.Buf.Bytes(), &chats)
	for index, chatname := range chats {
		chats[index] = DecodeBase64(chatname)
	}
	return
}

func (lc *LocalControler) RetrieveChat(chatname string) (c store.ChatHistory, err error) {
	response, err := lc.callMethod(
		pkg.RetrieveAvailableChats,
		fmt.Sprintf("/q?chatID=%s", EncodeBase64(chatname)),
		[]string{},
		map[string]string{},
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(response.Buf.Bytes(), &c)
	if err == nil {
		c.ChatName = DecodeBase64(c.ChatName)
	}
	return
}

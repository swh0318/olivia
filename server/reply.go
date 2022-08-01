package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olivia-ai/olivia/analysis"
	"github.com/olivia-ai/olivia/locales"
	"github.com/olivia-ai/olivia/user"

	"github.com/olivia-ai/olivia/util"
)

// GetReply
func GetReply(writer http.ResponseWriter, req *http.Request) {
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,Olivia-Token"
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	writer.Header().Set("Access-Control-Expose-Headers", "Authorization")

	var request RequestMessage
	if req == nil || req.Body == nil {
		fmt.Print("invalid request")
		return
	}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		fmt.Println(err)
		return
	}

	var responseSentence, responseTag string

	// Send a message from res/datasets/messages.json if it is too long
	if len(request.Content) > 500 {
		responseTag = "too long"
		responseSentence = util.GetMessage(request.Locale, responseTag)
	} else {
		// If the given locale is not supported yet, set english
		locale := request.Locale
		if !locales.Exists(locale) {
			locale = "en"
		}

		responseTag, responseSentence = analysis.NewSentence(
			locale, request.Content,
		).Calculate(*cache, neuralNetworks[locale], request.Token)
	}

	// Marshall the response in json
	response := ResponseMessage{
		Content:     responseSentence,
		Tag:         responseTag,
		Information: user.GetUserInformation(request.Token),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	writer.Write(bytes)
	//json.NewEncoder(writer).Encode(bytes)
}

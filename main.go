package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlekSi/alice"
)

// get dialod ID upon creation at https://dialogs.yandex.ru/developer/
const dialog_ID = "a6ffd39e-f6a2-4367-94d0-a47f3d7ef558"
const track_ID = "track_id"

func internalHandler(ctx context.Context, request *RequestWithState) (*ResponseWithState, error) {
	a := &ResponseWithState{
		Response: alice.ResponsePayload{
			EndSession: false,
		},
		SessionState: map[string]string{},
		Session:      copySession(request.Session),
		Version:      request.Version,
	}

	// Ping request
	if request.Request.OriginalUtterance == "ping" {
		a.Response.Text = "pong"
		return a, nil
	}

	// New session
	if request.Session.New {
		a.Response.Text = "Добро пожаловать! Назовите раздел и урок для включения воспроизведения" +
			"аудио; например `раздел один упражнение три`"
		return a, nil
	}

	// check which track was played
	if strings.Contains(request.Request.OriginalUtterance, "какой трек") {
		fileNum, err := getFileNumFromState(request)
		// check if zero value
		if err != nil || fileNum == 0 {
			a.Response.Text = "не могу определить трек"
			return a, nil
		}
		a.Response.Text = fmt.Sprintf("последним был проигран трек %d", fileNum)
		a.SessionState[track_ID] = strconv.Itoa(fileNum)
		return a, nil
	}

	// repeat
	if strings.Contains(request.Request.OriginalUtterance, "повтори") ||
		strings.Contains(request.Request.OriginalUtterance, "eщё раз") {
		fileNum, err := getFileNumFromState(request)
		// check if zero value
		if err != nil || fileNum == 0 {
			a.Response.Text = "не могу определить урок для повтора"
			return a, nil
		}
		audioID := numToUIDmap[fileNum]
		a.Response.Tts = audioUrlFromFileUID(audioID)
		a.SessionState[track_ID] = strconv.Itoa(fileNum)
		return a, nil
	}

	// next
	if strings.Contains(request.Request.OriginalUtterance, "следующий") {
		fileNum, err := getFileNumFromState(request)
		// check if zero value
		if err != nil || fileNum == 0 {
			a.Response.Text = "не могу определить следующий трек"
			return a, nil
		}
		fileNum = getNextTrackKey(fileNum)
		audioID := numToUIDmap[fileNum]
		a.Response.Tts = audioUrlFromFileUID(audioID)
		a.SessionState[track_ID] = strconv.Itoa(fileNum)
		return a, nil
	}

	// Session Continue
	// Check "Трек" -> explicit track id to play
	if strings.Contains(request.Request.OriginalUtterance, "трек") {
		fmt.Println("check `track` in Utterance")
		var trackID int

		fmt.Printf("NLUs: +%v\n", request.Request.NLU)
		if request.Request.NLU.Entities != nil &&
			len(request.Request.NLU.Entities) > 0 {
			if request.Request.NLU.Entities[0].Type == "YANDEX.NUMBER" {
				if tID, ok := request.Request.NLU.Entities[0].Value.(float64); !ok {
					a.SessionState["type_of_value"] = fmt.Sprintf("%T", request.Request.NLU.Entities[0].Value)
					a.Response.Text = "не могу точно определить трек для проигрывания"
					return a, nil
				} else {
					trackID = int(tID)
				}
			}
		} else {
			a.Response.Text = "не могу определить трек для проигрывания"
			return a, nil
		}
		fileUID := numToUIDmap[trackID]
		a.Response.Tts = audioUrlFromFileUID(fileUID)
		//a.Response.Text = "повторить? включить следующий урок?"
		a.SessionState[track_ID] = strconv.Itoa(trackID)
		return a, nil
	}

	// Check full match
	var matched = false
	// by Section and Exercise numbers
	if fileNum, found := sectionAndExerciseToFileNum[request.Request.Command]; found {
		audioUID := numToUIDmap[fileNum]
		a.Response.Tts = audioUrlFromFileUID(audioUID)
		a.SessionState[track_ID] = strconv.Itoa(fileNum)
		matched = true
	}

	// by Page and Exercise numbers
	if !matched {
		if fileNum, found := pageAndExerciseToFileNum[request.Request.OriginalUtterance]; found {
			audioUID := numToUIDmap[fileNum]
			a.Response.Tts = audioUrlFromFileUID(audioUID)
			a.SessionState[track_ID] = strconv.Itoa(fileNum)
			matched = true
		}
	}

	if !matched {
		a.Response.Text = "урок не обнаружен среди имеющихся; попробуйте снова"
	}

	return a, nil
}

//nolint:deadcode,unused
func Handler(ctx context.Context, request *RequestWithState) (*ResponseWithState, error) {
	fmt.Println("original request utterance:", request.Request.OriginalUtterance)
	fmt.Printf("REQUEST state: +%v\n", request.State)

	responseWithState, err := internalHandler(ctx, request)

	fmt.Printf("RESPONSE state: +%v\n", responseWithState.SessionState)
	return responseWithState, err
}

func getFileNumFromState(req *RequestWithState) (int, error) {
	fmt.Println("getFileNumFromState called")
	if trackID, found := req.State.Session[track_ID]; found {
		i, e := strconv.Atoi(trackID)
		fmt.Println("found", i, "; error", e)
		return i, e
	}
	fmt.Println("not found, zero returned")
	return 0, errors.New("сould not find prev track id")
}

func copySession(reqSession alice.RequestSession) alice.ResponseSession {
	return alice.ResponseSession{
		SessionID: reqSession.SessionID,
		MessageID: reqSession.MessageID,
		UserID:    reqSession.UserID,
	}
}

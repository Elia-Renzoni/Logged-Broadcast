package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log-b/model"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	backoffDeltaFactor int = 500
	backoffMaxRetries  int = 7
)

func RegisterToSeed(seedAddress, personalAddress string) {
	var (
		msg = model.BasicJoinMessage{
			SeedFlagOption:      true,
			NodeCompleteAddress: []string{personalAddress},
		}
		res          *http.Response
		connErr      error
		timeSleeping = 500
		retries      = 0
		exitStatus   bool
	)

EXP_BACKOFF:
	for {
		time.Sleep(time.Duration(timeSleeping) * time.Millisecond)
		client := http.Client{Timeout: 3 * time.Second}
		req, err := prepareRegistrationRequest(seedAddress, msg)
		if err != nil {
			panic(err)
		}
		res, connErr = client.Do(req)
		if errors.Is(connErr, context.DeadlineExceeded) || isNetError(connErr) {
			timeSleeping *= backoffDeltaFactor
			if retries >= backoffMaxRetries-1 {
				exitStatus = false
				break EXP_BACKOFF
			}
			retries += 1
			continue

		} else {
			if connErr != nil {
				exitStatus = false
				break EXP_BACKOFF
			}
		}

		// if the loop exit with a true the code below
		// must check the status code
		exitStatus = true
		break EXP_BACKOFF
	}

	if !exitStatus {
		panic(fmt.Sprintf("Registration Failed due to %s within %d retries", connErr, retries))
	}

	if res.StatusCode != 200 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var msg model.BasicPositiveAck
		json.Unmarshal(data, &msg)
		log.Println(msg.Succ)
	}
}

func prepareRegistrationRequest(seedAddress string, msg model.BasicJoinMessage) (*http.Request, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return http.NewRequest(
		http.MethodPost,
		generateFullHttpEndpoint(seedAddress),
		bytes.NewReader(data),
	)
}

func generateFullHttpEndpoint(seedAddress string) string {
	joinedURL, err := url.JoinPath("http://"+seedAddress, "/add-node")
	if err != nil {
		panic(err)
	}
	return joinedURL
}

func isNetError(err error) bool {
	_, ok := err.(net.Error)
	return ok
}

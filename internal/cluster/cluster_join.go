package cluster

import (
	"net/http"
	"log-b/model"
	"encoding/json"
	"bytes"
	"io"
	"time"
	"net"
	"errors"
	"net/url"
	"strconv"
)

const (
	backoffDeltaFactor int = 500
	backoffMaxRetries  int = 7
)

func RegisterToSeed(seedAddress, personalAddress string) {
	var (
		client   = http.Client{
			Timeout: 3 * time.Second,
		}
		msg      = model.BasicJoinMessage{
			NodeCompleteAddress: personalAddress,
		}
		req, err = prepareRegistrationRequest(seedAddress, msg) 
		res *http.Response
		connErr error
		timeSleeping = 500
		retries int = 0
		exitStatus bool
	)

	if err != nil {
		panic(err)
	}

EXP_BACKOFF:
	for {
		time.Sleep(time.Duration(timeSleeping) * time.Millisecond)
		res, connErr = client.Do(req)
		if nErr, ok := connErr.(net.Error); ok {
			if nErr.Timeout() {
				timeSleeping += backoffDeltaFactor
				if retries >= backoffMaxRetries {
					exitStatus = false
					break EXP_BACKOFF
				}
				retries += 1
				continue 
			}
		} else {
			exitStatus = false
			break EXP_BACKOFF
		}

		// if the loop exit with a true the code below
		// must check the status code
		exitStatus = true
		break EXP_BACKOFF
	}

	if !exitStatus {
		panic(errors.New("Registration Failed within retries: " + strconv.Itoa(retries)))
	}

	if res.StatusCode != 200 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		panic(data)
	}
}

func prepareRegistrationRequest(seedAddress string, msg model.BasicJoinMessage) (*http.Request, error) {
	data, err :=  json.Marshal(msg)
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
	joinedURL, err := url.JoinPath("http://" + seedAddress, "/add-node")
	if err != nil {
		panic(err)
	}
	return joinedURL
}

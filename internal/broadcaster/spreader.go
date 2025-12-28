package broadcaster

import (
	"bytes"
	"context"
	"errors"
	"log"
	"log-b/internal/cluster"
	"net/http"
	"strings"
	"time"
)

const (
	header      string = "application/json"
	ADD_NODE    string = "add-node"
	SET_DATA    string = "set-data"
	DELETE_DATA string = "delete-data"
)

func DoBroadcast(message []byte, methodRouter string) bool {
	memberlist := cluster.GetFullMembershipList()

	do := func() bool {
		var c ackCounter
		for _, node := range memberlist {
			endsystem := "http://" + node + "/" + methodRouter
			eval, err := send(endsystem, message, methodRouter)
			if err != nil {
				log.Fatal(err.Error())
			}

			if eval {
				c.inc()
			}
		}

		return c.isMajorityQuorumReached()
	}

	return do()
}

func send(addr string, msg []byte, methodRouter string) (bool, error) {
	var (
		res    *http.Response
		req    *http.Request
		err    error
		client = http.Client{Timeout: 5 * time.Second}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodRouter {
	case ADD_NODE, SET_DATA:
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, addr, bytes.NewBuffer(msg))
		if err != nil {
			log.Println(err.Error())
			return false, errors.New(makeError(err))
		}
		req.Header.Set("Content-Type", header)
		res, err = client.Do(req)
	case getPrefix(methodRouter):
		req, err = http.NewRequestWithContext(ctx, http.MethodDelete, addr, nil)
		if err != nil {
			log.Println(err.Error())
			return false, errors.New(makeError(err))
		}

		req.Header.Set("Content-Type", header)
		res, err = client.Do(req)
	default:
		return false, errors.New("invalid method router")
	}

	if err != nil {
		log.Fatal(err.Error())
		return false, errors.New(makeError(err))
	}

	return evaluateAck(res), nil
}

func evaluateAck(res *http.Response) bool {
	return res.StatusCode == 201 || res.StatusCode == 200
}

func makeError(err error) string {
	return "Error during Dial Broadcasting " + err.Error()
}

func getPrefix(endpoint string) string {
	if strings.Contains(endpoint, DELETE_DATA) {
		return endpoint
	}
	return ""
}

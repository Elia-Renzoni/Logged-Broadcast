package broadcaster

import (
	"bytes"
	"context"
	"log"
	"log-b/internal/cluster"
	"net/http"
	"time"
	"errors"
)

const (
	header      string = "application/json"
	ADD_NODE    string = "/add-node"
	SET_DATA    string = "/set-data"
	DELETE_DATA string = "/delete-data"
)

func DoBroadcast(message []byte, methodRouter string, addrWithEndpoints string) bool {
	var memberlist = cluster.GetFullMembershipList()

	do := func() bool {
		var c ackCounter
		for _, node := range memberlist {
			endsystem := node + addrWithEndpoints
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
		client = http.Client{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch methodRouter {
	case ADD_NODE, SET_DATA:
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, addr, bytes.NewBuffer(msg))
		if err != nil {
			return false, errors.New(errorMaker(err))
		}
		req.Header.Set("Content-Type", header)
		res, err = client.Do(req)

	case DELETE_DATA:
		req, err = http.NewRequestWithContext(ctx, http.MethodDelete, addr, nil)
		if err != nil {
			return false, errors.New(errorMaker(err))
		}

		req.Header.Set("Content-Type", header)
		res, err = client.Do(req)
	}

	if err != nil {
		return false, errors.New(errorMaker(err))
	}

	return evaluateAck(res), nil
}

func evaluateAck(res *http.Response) bool {
	return res.StatusCode == 201 || res.StatusCode == 200
}

func errorMaker(err error) string {
	return "Error during Dial Broadcasting " + err.Error()
}

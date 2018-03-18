package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func (w *Worker) handleWorkerResponse(conn net.Conn, rw http.ResponseWriter) error {
	r := bufio.NewReader(conn)
	jsonEncodedConfig, err := r.ReadBytes('\n')

	if err != nil {
		return err
	}

	workerResp := &WorkerResponse{}

	if err := json.Unmarshal(jsonEncodedConfig, workerResp); err != nil {
		return err
	}

	for key, value := range workerResp.Headers {
		rw.Header().Set(key, value)
	}

	rw.WriteHeader(workerResp.Status)

	if _, err := io.Copy(rw, conn); err != nil {
		return err
	}

	return nil
}

func (w *Worker) sendConfig(conn net.Conn, req *http.Request) error {
	serialized, jsonMarshalErr := json.Marshal(rackEnvBaseBy(req))

	if jsonMarshalErr != nil {
		return jsonMarshalErr
	}

	if _, err := conn.Write(serialized); err != nil {
		return err
	}

	if _, err := conn.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

func rackEnvBaseBy(req *http.Request) map[string]string {
	conf := make(map[string]string)

	if req.URL.Scheme == "" {
		conf["SCHEME"] = "HTTP"
	} else {
		conf["SCHEME"] = req.URL.Scheme
	}

	conf["HTTP_HOST"] = req.Host
	conf["HTTP_VERSION"] = req.Proto
	conf["SERVER_PROTOCOL"] = req.Proto

	conf["PATH_INFO"] = req.URL.Path
	conf["REQUEST_PATH"] = req.URL.Path
	conf["REQUEST_METHOD"] = req.Method
	conf["QUERY_STRING"] = req.URL.RawQuery
	conf["SERVER_NAME"] = req.URL.Hostname()
	conf["SERVER_ADDR"] = req.URL.Host
	conf["SERVER_PORT"] = req.URL.Port()
	conf["Content-Length"] = strconv.FormatInt(req.ContentLength, 10)
	conf["Content-Type"] = req.Header.Get("Content-type")
	conf["Transfer-Encoding"] = strings.Join(req.TransferEncoding, ",")
	conf["HTTP_COOKIE"] = req.Header.Get("Cookie")

	for key, values := range req.Header {
		conf[fmt.Sprintf("HTTP_%s", strings.ToUpper(key))] = strings.Join(values, ",")
	}

	return conf
}

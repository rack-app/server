package workers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (w *Worker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w.rp.ServeHTTP(rw, req)
}

func (w *Worker) ServeHTTPWithError(rw http.ResponseWriter, req *http.Request) error {

	conn, err := net.Dial("tcp", w.addr)

	if err != nil {
		return err
	}

	defer conn.Close()
	defer req.Body.Close()

	debug("---begin---")

	if err := forwardRequest(conn, req); err != nil {
		return err
	}

	if err := receiveResponse(conn, rw); err != nil {
		return err
	}

	debug("----end----")

	return nil

}

func forwardRequest(w io.Writer, req *http.Request) error {
	if err := sendEnv(w, req); err != nil {
		return err
	}

	go io.Copy(w, req.Body)

	return nil
}

func receiveResponse(r io.Reader, rw http.ResponseWriter) error {
	bufioReader := bufio.NewReader(r)

	rawStatus, err := fetchLine(bufioReader)

	if err != nil {
		return err
	}

	status, err := strconv.Atoi(string(rawStatus))

	if err != nil {
		return err
	}

	if err := populateHeader(bufioReader, rw); err != nil {
		return err
	}

	rw.WriteHeader(status)

	copied, err := io.Copy(rw, bufioReader)
	debug(fmt.Sprintf("-> %v", copied))

	if err != nil {
		return err
	}

	return nil
}

func populateHeader(r *bufio.Reader, rw http.ResponseWriter) error {
	h := rw.Header()

	for {

		line, err := fetchLine(r)

		if err != nil {
			return err
		}

		if len(line) == 0 {
			break
		}

		csvReader := csv.NewReader(bytes.NewReader(line))
		csvReader.Comma = '\t'

		row, err := csvReader.Read()

		if err != nil {
			return err
		}

		h.Add(row[0], row[1])

	}

	return nil
}

func fetchLine(reader *bufio.Reader) ([]byte, error) {
	wholeLine := []byte{}

	for {
		line, multipart, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		wholeLine = append(wholeLine, line...)

		if !multipart {
			break
		}
	}

	return wholeLine, nil
}

func sendEnv(wr io.Writer, req *http.Request) error {
	writer := csv.NewWriter(wr)
	writer.Comma = '\t'

	for key, value := range rackEnvFrom(req) {
		if err := writer.Write([]string{key, value}); err != nil {
			return err
		}
	}

	if err := writer.Write([]string{}); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

func rackEnvFrom(req *http.Request) map[string]string {
	conf := make(map[string]string)

	if req.URL.Scheme == "" {
		conf["SCHEME"] = "HTTP"
	} else {
		conf["SCHEME"] = req.URL.Scheme
	}

	conf["HTTP_HOST"] = req.Host
	conf["HTTP_VERSION"] = req.Proto

	conf["SERVER_NAME"] = "rack-app-server"
	conf["SERVER_PORT"] = req.URL.Port()
	conf["SERVER_PROTOCOL"] = req.Proto

	conf["PATH_INFO"] = req.URL.Path
	conf["QUERY_STRING"] = req.URL.RawQuery
	conf["REQUEST_PATH"] = req.URL.Path
	conf["REQUEST_METHOD"] = req.Method
	conf["REMOTE_ADDR"] = req.RemoteAddr
	conf["Content-Length"] = strconv.FormatInt(req.ContentLength, 10)
	conf["Transfer-Encoding"] = strings.Join(req.TransferEncoding, ",")

	if cookie := req.Header.Get("Cookie"); cookie != "" {
		conf["HTTP_COOKIE"] = cookie
	}

	if ctype := req.Header.Get("Content-type"); ctype != "" {
		conf["Content-Type"] = ctype
	}

	for key, values := range req.Header {

		if len(values) == 0 {
			continue
		}

		formattedKey := strings.Replace(strings.ToUpper(key), "-", "_", -1)

		conf[fmt.Sprintf("HTTP_%s", formattedKey)] = strings.Join(values, ",")

	}

	return conf
}

func debug(v ...interface{}) {
	if os.Getenv("RACK_APP_DEBUG") != "" {
		fmt.Println(v...)
	}
}

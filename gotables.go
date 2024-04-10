package gotables

import (
	"bytes"
	"encoding/json"
	"errors"
	"git.jereileu.ch/gotables/server/gt-server/fs"
	"git.jereileu.ch/gotables/server/gt-server/server"
	"git.jereileu.ch/gotables/server/gt-server/table"
	"io"
	"net/http"
)

type Config struct {
	fs.Conf
	Host string
}
type requestBody = server.Body
type Column = table.Column
type Table = table.Table
type TableU = table.TableU

func ConstructUrl(tbl string, db string, config Config) (string, error) {
	out := ""
	if config.HTTPSMode {
		out += "https://"
	} else {
		out += "http://"
	}
	if config.Host != "" {
		out += config.Host
	} else {
		return "", errors.New("invalid url: host needs to be set")
	}
	if config.Port != "" {
		out += ":"
		out += config.Port
	}
	if tbl != "" {
		if db != "" {
			out += "/"
			out += db
			out += "/"
			out += tbl
		} else {
			return "", errors.New("invalid url: table set but db is not")
		}
	} else {
		if db != "" {
			out += "/"
			out += db
		}
	}
	return out, nil
}

func ConstructRequest(body requestBody, url string) (*http.Request, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	return req, err
}

func DoRequest(req *http.Request) (Table, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Table{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Table{}, err
	}
	tblU := TableU{}
	err = tblU.FromJ(body)
	if err != nil {
		return Table{}, err
	}
	tbl, err := tblU.ToT()
	return tbl, err
}

func RunServer(config Config) {
	server.Run(config.Conf)
}

func TestServer(config Config) error {
	url, err := ConstructUrl("", "", config)
	if err != nil {
		return err
	}
	_, err = http.Get(url)
	return err
}

// Root

func ShowDBs(sessionId string, config Config) (Table, error) {
	query := "show"
	return request(query, "", "", sessionId, config)
}

// DB

func ShowTables() {}
func CreateDB()   {}
func SetDBName()  {}
func CopyDB()     {}
func DeleteDB()   {}

// Table

func ShowTable()    {}
func CreateTable()  {}
func SetTableName() {}
func CopyTable()    {}
func DeleteTable()  {}

// Column

func ShowColumn()       {}
func CreateColumn()     {}
func SetColumnName()    {}
func SetColumnDefault() {}
func CopyColumn()       {}
func DeleteColumn()     {}

// Row

func ShowRow()   {}
func CreateRow() {}
func SetRow()    {}
func CopyRow()   {}
func DeleteRow() {}

// User

// Backup

func request(query string, tbl string, db string, sessionId string, config Config) (Table, error) {
	url, err := ConstructUrl(tbl, db, config)
	if err != nil {
		return Table{}, err
	}
	body := requestBody{
		Query:     query,
		SessionId: sessionId,
	}
	req, err := ConstructRequest(body, url)
	if err != nil {
		return Table{}, err
	}
	return DoRequest(req)
}

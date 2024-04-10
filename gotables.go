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
	"strconv"
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

func Request(query string, tbl string, db string, sessionId string, config Config) (Table, error) {
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
	return Request(query, "", "", sessionId, config)
}

// DB

func ShowTables(db string, sessionId string, config Config) (Table, error) {
	query := "show"
	return Request(query, "", db, sessionId, config)
}

func CreateDB(db string, sessionId string, config Config) (Table, error) {
	query := "create"
	return Request(query, "", db, sessionId, config)
}

func SetDBName(name string, db string, sessionId string, config Config) (Table, error) {
	query := "set name " + name
	return Request(query, "", db, sessionId, config)
}

func CopyDB(name string, db string, sessionId string, config Config) (Table, error) {
	query := "copy " + name
	return Request(query, "", db, sessionId, config)
}

func DeleteDB(db string, sessionId string, config Config) (Table, error) {
	query := "delete"
	return Request(query, "", db, sessionId, config)
}

// Table

func ShowTable(tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "show"
	return Request(query, tbl, db, sessionId, config)
}

func ShowTableColumns(columns []string, tbl string, db string, sessionId string, config Config) (Table, error) {
	colNames := ""
	for i := 0; i < len(columns); i++ {
		colNames += columns[i]
		if i != len(columns)-1 {
			colNames += ":"
		}
	}
	query := "show " + colNames
	return Request(query, tbl, db, sessionId, config)
}

func ShowTableConditions(conditions []string, columns []string, tbl string, db string, sessionId string, config Config) (Table, error) {
	colNames := ""
	for i := 0; i < len(columns); i++ {
		colNames += columns[i]
		if i != len(columns)-1 {
			colNames += ":"
		}
	}
	conditionString := ""
	for i := 0; i < len(conditions); i++ {
		conditionString += conditions[i]
		if i != len(conditions)-1 {
			conditionString += " "
		}
	}
	query := "show " + colNames + " where " + conditionString
	return Request(query, tbl, db, sessionId, config)
}

func CreateTable(tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "create"
	return Request(query, tbl, db, sessionId, config)
}

func SetTableName(name string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "set name " + name
	return Request(query, tbl, db, sessionId, config)
}

func CopyTable(name string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "copy " + name
	return Request(query, tbl, db, sessionId, config)
}

func DeleteTable(tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "delete"
	return Request(query, tbl, db, sessionId, config)
}

// Column

func ShowColumn(column string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column show " + column
	return Request(query, tbl, db, sessionId, config)
}

func ShowColumns(columns []string, tbl string, db string, sessionId string, config Config) (Table, error) {
	colNames := ""
	for i := 0; i < len(columns); i++ {
		colNames += columns[i]
		if i != len(columns)-1 {
			colNames += ":"
		}
	}
	query := "column show " + colNames
	return Request(query, tbl, db, sessionId, config)
}

func CreateColumn(column Column, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column create " + column.Name + ":" + column.Type + ":" + column.Default
	return Request(query, tbl, db, sessionId, config)
}

func SetColumnName(name string, column string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column set name " + column + " " + name
	return Request(query, tbl, db, sessionId, config)
}

func SetColumnDefault(def string, column string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column set default " + column + " " + def
	return Request(query, tbl, db, sessionId, config)
}

func CopyColumn(name string, column string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column copy " + column + " " + name
	return Request(query, tbl, db, sessionId, config)
}

func DeleteColumn(column string, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "column delete " + column
	return Request(query, tbl, db, sessionId, config)
}

// Row

func ShowRow(row int, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "row show " + strconv.Itoa(row)
	return Request(query, tbl, db, sessionId, config)
}

func CreateRow(values [][2]string, tbl string, db string, sessionId string, config Config) (Table, error) {
	valueString := ""
	for i := 0; i < len(values); i++ {
		valueString += values[i][0] + ":" + values[i][1]
		if i != len(values)-1 {
			valueString += " "
		}
	}
	query := "row create " + valueString
	return Request(query, tbl, db, sessionId, config)
}

func SetRow(value string, colName string, row int, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "row set " + strconv.Itoa(row) + ":" + colName + " " + value
	return Request(query, tbl, db, sessionId, config)
}

func CopyRow(row int, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "row copy " + strconv.Itoa(row)
	return Request(query, tbl, db, sessionId, config)
}

func DeleteRow(row int, tbl string, db string, sessionId string, config Config) (Table, error) {
	query := "row delete " + strconv.Itoa(row)
	return Request(query, tbl, db, sessionId, config)
}

// User

// Backup

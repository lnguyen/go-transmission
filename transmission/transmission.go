package transmission

import (
	"encoding/json"
	"github.com/longnguyen11288/go-transmission/client"
)

var ()

type TransmissionClient struct {
	apiclient client.ApiClient
}

type command struct {
	Method    string    `json:"method,omitempty"`
	Arguments arguments `json:"arguments,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type arguments struct {
	Fields   []string  `json:"fields,omitempty"`
	Torrents []Torrent `json:"torrents,omitempty"`
	Ids      []int     `json:"ids,omitempty"`
}

type Torrent struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Status        int     `json:"status"`
	LeftUntilDone int     `json:"leftUntilDone"`
	Eta           int     `json:"eta"`
	UploadRatio   float64 `json:"uploadRatio"`
	RateDownload  int     `json:"rateDownload"`
	RateUpload    int     `json:"rateUpload"`
}

func New(url string,
	username string, password string) TransmissionClient {
	apiclient := client.NewClient(url, username, password)
	tc := TransmissionClient{apiclient: apiclient}
	return tc
}

func (ac *TransmissionClient) GetTorrents() ([]Torrent, error) {
	var getCommand command
	var outputCommand command
	getCommand.Method = "torrent-get"
	getCommand.Arguments.Fields = []string{"id", "name",
		"status", "leftUntilDone", "eta", "uploadRatio",
		"rateDownload", "rateUpload"}
	body, err := json.Marshal(getCommand)
	if err != nil {
		return []Torrent{}, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return []Torrent{}, err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return []Torrent{}, err
	}
	return outputCommand.Arguments.Torrents, nil
}

func (ac *TransmissionClient) RemoveTorrent(id int) (string, error) {
	var removeCommand command
	var outputCommand command

	removeCommand.Method = "torrent-remove"
	removeCommand.Arguments.Ids = []int{id}
	body, err := json.Marshal(removeCommand)
	if err != nil {
		return "", err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return "", err
	}
	return outputCommand.Result, nil
}

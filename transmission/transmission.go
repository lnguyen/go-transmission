package transmission

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/longnguyen11288/go-transmission/client"
)

//TransmissionClient to talk to transmission
type TransmissionClient struct {
	apiclient client.ApiClient
}

type command struct {
	Method    string    `json:"method,omitempty"`
	Arguments arguments `json:"arguments,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type arguments struct {
	Fields       []string     `json:"fields,omitempty"`
	Torrents     []Torrent    `json:"torrents,omitempty"`
	Ids          []int        `json:"ids,omitempty"`
	DeleteData   bool         `json:"delete-local-data,omitempty"`
	DownloadDir  string       `json:"download-dir,omitempty"`
	MetaInfo     string       `json:"metainfo,omitempty"`
	Filename     string       `json:"filename,omitempty"`
	TorrentAdded TorrentAdded `json:"torrent-added"`
}

//Torrent struct for torrents
type Torrent struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Status        int     `json:"status"`
	LeftUntilDone int     `json:"leftUntilDone"`
	Eta           int     `json:"eta"`
	UploadRatio   float64 `json:"uploadRatio"`
	RateDownload  int     `json:"rateDownload"`
	RateUpload    int     `json:"rateUpload"`
	DownloadDir   string  `json:"downloadDir"`
	IsFinished    bool    `json:"isFinished"`
	PercentDone   float64 `json:"percentDone"`
	SeedRatioMode int     `json:"seedRatioMode"`
}

//TorrentAdded data returning
type TorrentAdded struct {
	HashString string `json:"hashString"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
}

//New create new transmission torrent
func New(url string, username string, password string) TransmissionClient {
	apiclient := client.NewClient(url, username, password)
	tc := TransmissionClient{apiclient: apiclient}
	return tc
}

//GetTorrents get a list of torrents
func (ac *TransmissionClient) GetTorrents() ([]Torrent, error) {
	var getCommand command
	var outputCommand command
	getCommand.Method = "torrent-get"
	getCommand.Arguments.Fields = []string{"id", "name",
		"status", "leftUntilDone", "eta", "uploadRatio",
		"rateDownload", "rateUpload", "downloadDir",
		"isFinished", "percentDone", "seedRatioMode"}
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

//StartTorrent start the torrent
func (ac *TransmissionClient) StartTorrent(id int) (string, error) {
	return ac.sendSimpleCommand("torrent-start", id)
}

//StopTorrent start the torrent
func (ac *TransmissionClient) StopTorrent(id int) (string, error) {
	return ac.sendSimpleCommand("torrent-stop", id)
}

//RemoveTorrent remove the torrents
func (ac *TransmissionClient) RemoveTorrent(id int, removeFile bool) (string, error) {
	var removeCommand command
	var outputCommand command

	removeCommand.Method = "torrent-remove"
	removeCommand.Arguments.Ids = []int{id}
	removeCommand.Arguments.DeleteData = removeFile
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

//AddTorrentByURL add torrent by url
func (ac *TransmissionClient) AddTorrentByURL(url string, downloadDir string) (TorrentAdded, error) {
	resp, err := http.Get(url)
	if err != nil {
		return TorrentAdded{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TorrentAdded{}, err
	}
	return ac.addTorrent(body, downloadDir)
}

//AddTorrentByFilename add torrent by filename
func (ac *TransmissionClient) AddTorrentByFilename(filename string, downloadDir string) (TorrentAdded, error) {
	var addCommand command
	var outputCommand command

	addCommand.Method = "torrent-add"
	addCommand.Arguments.Filename = filename
	addCommand.Arguments.DownloadDir = downloadDir

	body, err := json.Marshal(addCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return TorrentAdded{}, err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	return outputCommand.Arguments.TorrentAdded, nil
}

//AddTorrentByFile add torrent by file
func (ac *TransmissionClient) AddTorrentByFile(file string, downloadDir string) (TorrentAdded, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return TorrentAdded{}, err
	}
	return ac.addTorrent(fileData, downloadDir)
}

func (ac *TransmissionClient) addTorrent(data []byte, downloadDir string) (TorrentAdded, error) {
	var addCommand command
	var outputCommand command

	addCommand.Method = "torrent-add"
	addCommand.Arguments.MetaInfo = base64.StdEncoding.EncodeToString(data)
	addCommand.Arguments.DownloadDir = downloadDir

	body, err := json.Marshal(addCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return TorrentAdded{}, err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	return outputCommand.Arguments.TorrentAdded, nil
}

func encodeFile(file string) (string, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}

func (ac *TransmissionClient) sendSimpleCommand(method string, id int) (result string, err error) {
	cmd := command{Method: method}
	cmd.Arguments.Ids = []int{id}
	resp, err := ac.sendCommand(cmd)
	return resp.Result, err
}

func (ac *TransmissionClient) sendCommand(cmd command) (response command, err error) {
	body, err := json.Marshal(cmd)
	if err != nil {
		return
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return
	}
	err = json.Unmarshal(output, &response)
	if err != nil {
		return
	}
	return response, nil
}

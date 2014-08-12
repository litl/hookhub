package fogbugz

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Session struct {
	host  string
	token string
}

type Config interface {
	GetEmail() string
	GetPassword() string
	GetHost() string
}

type errorCodeTag struct {
	ErrorCode int    `xml:"code,attr"`
	ErrorDesc string `xml:",chardata"`
}

type caseTag struct {
	BugNumber  string `xml:"ixBug,attr"`
	Operations string
}

func (session *Session) ResolveBug(caseNumber string, comment string) error {
	//$ http --form POST https://litl.fogbugz.com/api.asp cmd=resolve token=igl1mpdt5gu4q3qod50f4gb3nbnedb ixBug=54393 sEvent=sandydidit
	values := url.Values{}
	values.Set("cmd", "resolve")
	values.Set("token", session.token)
	values.Set("ixBug", caseNumber)
	values.Set("sEvent", comment) // TODO: Sanitize?
	url := &url.URL{"https", "", nil, session.host, "/api.asp", "", ""}

	resp, err := http.PostForm(url.String(), values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (session *Session) FileBug(project string, area string, title string, content string) (string, error) {
	values := url.Values{}
	values.Set("cmd", "new")
	values.Set("token", session.token)
	values.Set("sProject", project)
	values.Set("sArea", area)
	values.Set("sTitle", title)
	values.Set("sEvent", content)
	url := &url.URL{"https", "", nil, session.host, "/api.asp", "", ""}

	resp, err := http.PostForm(url.String(), values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Failed to create bug report: %s", resp.Status))
		return "", err
	}

	type Response struct {
		XMLName xml.Name     `xml:"response"`
		Case    caseTag      `xml:"case"`
		Error   errorCodeTag `xml:"error"`
	}

	var r Response
	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&r)
	if err != nil {
		return "", err
	}

	if r.Error.ErrorCode != 0 {
		return "", errors.New(fmt.Sprintf("Fogbugz error %d: %s", r.Error.ErrorCode, r.Error.ErrorDesc))
	}

	return fmt.Sprintf("https://%s/default.asp?%s", session.host, r.Case.BugNumber), nil
}

func fecthAuthToken(config Config) (string, error) {
	values := url.Values{}
	values.Set("cmd", "logon")
	values.Set("email", config.GetEmail())
	values.Set("password", config.GetPassword())
	url := &url.URL{"https", "", nil, config.GetHost(), "/api.asp", values.Encode(), ""}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	type Response struct {
		XMLName xml.Name     `xml:"response"`
		Token   string       `xml:"token"`
		Error   errorCodeTag `xml:"error"`
	}

	var r Response
	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&r)
	if err != nil {
		return "", err
	}

	if r.Error.ErrorCode != 0 {
		return "", errors.New(fmt.Sprintf("Fogbugz error %d: %s", r.Error.ErrorCode, r.Error.ErrorDesc))
	}

	return r.Token, nil
}

func NewSession(config Config) (*Session, error) {
	session := new(Session)

	var err error
	if session.token, err = fecthAuthToken(config); err != nil {
		return nil, err
	}

	session.host = config.GetHost()

	return session, nil
}

func (session *Session) String() string {
	return fmt.Sprintf("FogBugzSession for token %s", session.token)
}

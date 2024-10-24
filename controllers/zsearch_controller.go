package controllers

import (
	"net/http"
	"net/url"
	//"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"fmt"
	"strings"
	"bytes"
	"errors"

	"apiwholesale/models"
	s "apiwholesale/system"

	"github.com/goji/param"
	"github.com/zenazn/goji/web"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type TokenResponse struct {
	Token_type   string  `json:"token_type"`
	AccessToken  string  `json:"access_token"`
	Expires_in   string  `json:"expires_in"`
	Scope        string  `json:"scope"`
}

var CurrentMSToken string

func renewMSToken(force bool) error {
	if len( CurrentMSToken ) == 0 || force == true {
		tokenUrl := "https://api.cognitive.microsoft.com/sts/v1.0/issueToken?Subscription-Key=HIDDEN"

		data := url.Values{}
		resp, err := http.PostForm(tokenUrl, data)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		token, _ := ioutil.ReadAll(resp.Body)
		CurrentMSToken = string(token)
	}

	return nil
}

func requestMSTranslate( search string ) (string, error) {

	if err := renewMSToken(false); err != nil {
		s.Log.Err( "[error] Token request failed" )
		return "", err
	}

	data := url.Values{}
	data.Set("text", search)
	// data.Add("from", "ru") // let's Bing detect our language
	data.Add("to", "zh-CHT")

	transUrl := "https://api.microsofttranslator.com/v2/http.svc/Translate"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", transUrl, nil)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "Bearer " + CurrentMSToken)

	result, err := client.Do(req)

	if err != nil || result.StatusCode != 200 {
		// here we have to force request new token
		s.Log.Info( "[info] Token renewed!" )
		if err := renewMSToken(true); err != nil {
			s.Log.Err( "[error] Token request failed" )
			return "", err
		}

		req.Header.Set("Authorization", "Bearer " + CurrentMSToken)
		result, err = client.Do(req)

		if err != nil || result.StatusCode != 200 {
			s.Log.Err( "[error] MS Translation limit" )
			return "", errors.New("MS Translation limit")
		}
	}
	defer result.Body.Close()

	res, e := ioutil.ReadAll(result.Body)
	if e != nil {
		return "", e
	}

	// s.DEBUG( string(res) )

	type Translated struct {
		Result string `xml:",chardata"`
	}
	v := Translated{}
	err = xml.Unmarshal(res, &v)

	return v.Result, err
}

type SearchQuery struct {
	Q string `param:"q"`
}

func Search(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		result string
		err error
		query SearchQuery
	)

	if r.ContentLength < 1 {
		err := errors.New("Empty query string")
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] empty query string"))
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] empty query string"))
		return
	}

	err = param.Parse(r.Form, &query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] empty query string"))
		return
	}

	search := strings.TrimSpace( query.Q )
	if len( search ) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] empty query string"))
		return
	}

	result, err = requestMSTranslate( search )

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.Log.Err(fmt.Sprintf("[error] request MS translation %s", err.Error()))
		return
	}

	gbk, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(result)), simplifiedchinese.GBK.NewEncoder()))

	translated := models.SearchResponse{}
	translated.QueryRu      = search
	translated.ResultZh     = result
	translated.ResultZhGBK  = url.QueryEscape(string(gbk))
	translated.Source       = "bing"

	models.SearchInsert( translated )

	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder.Encode(translated)
}

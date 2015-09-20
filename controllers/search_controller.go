package controllers

import (
    "net/http"
    "net/url"
    "crypto/tls"
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
    "fmt"
    "strings"
    "apiwholesale/models"

    "github.com/goji/param"
    "github.com/zenazn/goji/web"
    s "apiwholesale/system"
)

type TokenResponse struct {
    Token_type   string  `json:"token_type"`
    AccessToken  string  `json:"access_token"`
    Expires_in   string  `json:"expires_in"`
    Scope        string  `json:"scope"`
}

var CurrentMSToken string

func getMSToken(force bool) error {
    if len( CurrentMSToken ) == 0 || force == true {
        apiUrl := "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13/"
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        client := &http.Client{Transport: tr}

        data := url.Values{}
        data.Set("client_id", "jurbat-alibaba24")
        data.Add("client_secret", "54Sko8XnhS/5b4Gt/YKG39AknOhE2OljChb3tZNOZIc=")
        data.Add("scope", "http://api.microsofttranslator.com")
        data.Add("grant_type", "client_credentials")

        resp, err := client.PostForm(apiUrl, data)
        if err != nil {
            return err
        }

        token := TokenResponse{}
        if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
            return err
        }

        CurrentMSToken = token.AccessToken
    }

    return nil
}

func requestMSTranslate( search string ) (string, error) {

    if err := getMSToken(false); err != nil {
        return "", err
    }

    data := url.Values{}
    data.Set("text", search)
    data.Add("from", "ru")
    data.Add("to", "zh-CHT")

    transUrl := "http://api.microsofttranslator.com/v2/Http.svc/Translate"

    client := &http.Client{}
    req, _ := http.NewRequest("GET", transUrl, nil)
    req.URL.RawQuery = data.Encode()
    req.Header.Set("Authorization", "Bearer " + CurrentMSToken)

    result, err := client.Do(req)

    if err != nil {
        if err = getMSToken(true); err != nil {
            s.Log.Err(fmt.Sprintf("[error] MS token %s", err.Error()))
            return "", err
        }

        // here we have new token
        req.Header.Set("Authorization", "Bearer " + CurrentMSToken)
        result, err = client.Do(req)
        if err != nil {
            s.Log.Err(fmt.Sprintf("[error] MS Translation limit %s", err.Error()))
            return "", err
        }
    }

    defer result.Body.Close()

    res, e := ioutil.ReadAll(result.Body)
    if e != nil {
        return "", e
    }

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


    r.ParseForm()
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

    translated := models.SearchResponse{}
    translated.QueryRu  = search
    translated.ResultZh = result
    translated.Source   = "bing"

    go func(){
        models.SearchInsert( translated )
    }()


    encoder := json.NewEncoder(w)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    encoder.Encode(translated)
}

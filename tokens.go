package npmgo

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TokenInfo struct {
	Objects []Objects `json:"objects"`
}
type Objects struct {
	Id             string         `json:"key"`
	Token          string         `json:"token"`
	Created        string         `json:"created"`
	Readonly       bool           `json:"readonly"`
	CidrWhitelist  []string	      `json:"cidr_whitelist"`
}

//type Objects struct {
//	TokenInfo string `json:"objects"`
//}

// Settings used to create users. Tags must be comma-separated.
type TokenSettings struct {
	//Password is required again to create token
	Password      string         `json:"password"`
	Readonly      bool           `json:"readonly"`
	CidrWhitelist []string       `json:"cidr_whitelist"`
}

//
// GET /-/npm/v1/tokens
//

// Example response:
// "objects":[{"token":"9d6cef","key":"a8430c124cc6968da1f5030f0fa1574529e2ada1be352a615703e34def3a624bc1f6b6f95485e4c0253065b1323c40713ee9308228dd557a874f778c8275bce1","cidr_whitelist":null,"readonly":false,"created":"2019-01-23T18:53:40.800Z","updated":"2019-01-23T18:53:40.800Z"},{"token":"28faa0","key":"7784edd1e237e24ebbe8e25da4c629e581e566ca1b53714c2eda35c9c9b5e040cbecae0a04ca760aea7215c47d51a7c93b42dbba603e64315f0c1df47b0b26a0","cidr_whitelist":[],"readonly":false,"created":"2019-01-23T18:46:06.554Z","updated":"2019-01-23T18:46:06.554Z"}],"total":2,"urls":{}}

// Returns a list of all users in a cluster.
func (c *Client) ListTokens() (rec []Objects, err error) {
	req, err := newGETRequest(c, "npm/v1/tokens")
	if err != nil {
		return []Objects{}, err
	}

	var tokeninfo TokenInfo

	if err = executeAndParseRequest(c, req, &tokeninfo); err != nil {
		return []Objects{}, err
	}
	rec = tokeninfo.Objects
	
	return rec, nil
}

//
// GET /-/tokens/ (specific token)
//

// Returns information about individual user.
func (c *Client) GetToken(id string) (rec Objects, err error) {
	req, err := newGETRequest(c, "npm/v1/tokens")
	if err != nil {
		return Objects{}, err
	}

	var tokeninfo TokenInfo

	if err = executeAndParseRequest(c, req, &tokeninfo); err != nil {
		return Objects{}, err
	}
	objects := tokeninfo.Objects
        for _, i := range objects {
                if id == i.Id {
                        rec = i
                        break
                }
        }
	return rec, nil
}

//
// POST /-/npm/v1/tokens
//

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
    b, err := json.Marshal(v)

    if safeEncoding {
        b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
        b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
        b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
    }
    return b, err
}

// Creates a new token.
func (c *Client) CreateToken(info TokenSettings) (res Objects, err error) {
	if info.CidrWhitelist == nil {
		info.CidrWhitelist = []string{}
	}
	body, err := JSONMarshal(info, true)
	if err != nil {
		return Objects{}, err
	}

	req, err := newRequestWithBody(c, "POST", "npm/v1/tokens", body)
	if err != nil {
		return Objects{}, err
	}

	//var res Objects
	//res, err = executeRequest(c, req)
	if err = executeAndParseRequest(c, req, &res); err != nil {
		return Objects{}, err
	}
	if err != nil {
		return Objects{}, err
	}

	return res, nil
}

//
// DELETE /-/npm/v1/tokens/token/{id}
//

// Deletes token.
func (c *Client) DeleteToken(id string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "npm/v1/tokens/token/"+PathEscape(id), nil)
	if err != nil {
		return nil, err
	}

	res, err = executeRequest(c, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

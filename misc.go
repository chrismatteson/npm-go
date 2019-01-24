package npmgo

//
// GET /-/whoami
//

type WhoamiInfo struct {
	Username        string `json:"username"`
}

func (c *Client) Whoami() (rec *WhoamiInfo, err error) {
	req, err := newGETRequest(c, "whoami")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

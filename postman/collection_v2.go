package postman

type collectionV2 struct {
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"info"`
	Items []itemV2 `json:"item"`
}

type itemV2 struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Request     requestV2    `json:"request"`
	Response    []responseV2 `json:"response"`
	Items       []itemV2     `json:"item"`
}

type requestV2 struct {
	Method      string            `json:"method"`
	Headers     []requestHeaderV2 `json:"header"`
	Body        requestBodyV2     `json:"body"`
	URL         interface{}       `json:"url"`
	Description string            `json:"description"`
}

type requestHeaderV2 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type requestBodyV2 struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type requestURLV2 struct {
	Raw string `json:"raw"`
}

type responseV2 struct {
	Name            string            `json:"name"`
	OriginalRequest requestV2         `json:"originalRequest"`
	Status          string            `json:"status"`
	Code            int               `json:"code"`
	Headers         []requestHeaderV2 `json:"header"`
	Body            string            `json:"body"`
}

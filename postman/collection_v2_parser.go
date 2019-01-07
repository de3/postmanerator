package postman

import (
	"encoding/json"
	"html"
	"strconv"
)

type CollectionV2Parser struct {
	PostmanVersion string
}

func (p *CollectionV2Parser) CanParse(contents []byte) bool {
	if p.PostmanVersion == "v2" {
		return true
	}
	return false
}

func (p *CollectionV2Parser) Parse(contents []byte, options BuilderOptions) (Collection, error) {
	src := collectionV2{}
	if err := json.Unmarshal(contents, &src); err != nil {
		return Collection{}, err
	}
	return p.buildCollectionFromV2(src, options)
}

func (p *CollectionV2Parser) buildCollectionFromV2(src collectionV2, options BuilderOptions) (Collection, error) {
	collection := Collection{
		Name:        src.Info.Name,
		Description: src.Info.Description,
		Requests:    make([]Request, 0),
		Folders:     make([]Folder, 0),
		Structures:  make([]StructureDefinition, 0),
	}

	// build folder based on item data
	for _, item := range src.Items {
		newFolder := Folder{}
		newRequest := Request{}
		newFolder, newRequest, _ = p.buildFolderFromItem(item, newFolder)

		collection.Folders = append(collection.Folders, newFolder)
		if newRequest.Name != "" {
			collection.Requests = append(collection.Requests, newRequest)
		}
	}

	return collection, nil
}

func (p *CollectionV2Parser) buildFolderFromItemRecurse(item itemV2, folder Folder) (Folder, error) {
	if len(item.Items) > 0 {
		for idx := range item.Items {
			subfolder := Folder{}
			subfolder, _ = p.buildFolderFromItemRecurse(item.Items[idx], subfolder)
			folder.Folders = append(folder.Folders, subfolder)
		}
	}

	folder.ID = item.ID
	folder.Name = item.Name
	folder.Description = item.Description

	var req Request
	if item.Request.Method != "" {
		req, _ = p.buildRequest(item)
		folder.Requests = append(folder.Requests, req)
	}

	return folder, nil
}

func (p *CollectionV2Parser) buildFolderFromItem(item itemV2, folder Folder) (Folder, Request, error) {
	var req Request

	if item.Request.Method != "" {
		req, _ = p.buildRequest(item)
	}

	if len(item.Items) == 0 {
		return Folder{}, req, nil
	}

	folder.ID = item.ID
	folder.Name = item.Name
	folder.Description = item.Description

	for _, subitem := range item.Items {
		subreq, _ := p.buildRequest(subitem)
		folder.Requests = append(folder.Requests, subreq)
	}

	return folder, req, nil
}

func (p *CollectionV2Parser) buildRequest(item itemV2) (Request, error) {
	var url string
	switch v := item.Request.URL.(type) {
	case map[string]interface{}:
		urlV2, _ := v["raw"].(string)
		url = urlV2
	case string:
		url, _ = item.Request.URL.(string)
	}

	request := Request{
		Name:        item.Name,
		Description: item.Request.Description,
		Method:      item.Request.Method,
		URL:         url,
		PayloadType: item.Request.Body.Mode,
		PayloadRaw:  item.Request.Body.Raw,
		Headers:     p.buildHeaders(item.Request.Headers),
		Responses:   p.buildResponse(item.Response),
	}

	return request, nil
}

func (p *CollectionV2Parser) buildHeaders(headersV2 []requestHeaderV2) []KeyValuePair {
	headers := make([]KeyValuePair, 0)
	for _, h := range headersV2 {
		headers = append(headers, KeyValuePair{
			Name:  h.Key,
			Key:   h.Key,
			Value: html.EscapeString(h.Value),
		})
	}
	return headers
}

func (p *CollectionV2Parser) buildResponse(responseV2 []responseV2) []Response {
	responses := make([]Response, 0)
	for idx, res := range responseV2 {
		responses = append(responses, Response{
			ID:         strconv.Itoa(idx),
			Name:       res.Name,
			Status:     res.Status,
			StatusCode: res.Code,
			Body:       res.Body,
			Headers:    p.buildHeaders(res.Headers),
		})
	}

	return responses
}

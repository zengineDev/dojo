package dojo

import "encoding/json"

type JsonResponseBody struct {
	Data interface{} `json:"data"`
}

func (app Application) JSON(ctx Context, data interface{}) error {
	resBody := JsonResponseBody{
		Data: data,
	}
	respData, err := json.Marshal(resBody)
	if err != nil {
		return err
	}
	_, err = ctx.Response().Write(respData)
	if err != nil {
		return err
	}
	return nil
}

package dropletgenomics

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	PumpSetTargetVolume clientInvocation = iota
	PumpReset
	PumpToggleWithdrawInfuse
	PumpSetVolume
	PumpToggle
	PumpRefresh
	PumpPurge
)

type dataPack struct {
	DataEscaped string `json:"data_pack"`
	Success     int    `json:"success"`
}

type pumpNameAndValues map[string]pump

type pump struct {
	VolumeTarget  float64 `json:"volumeTarget"`
	PurgeRate     float64 `json:"purge_rate"`
	PumpID        float64 `json:"pump_id"`
	RateW         float64 `json:"rateW"`
	Volume        float64 `json:"volume"`
	Status        bool    `json:"status"`
	Name          string  `json:"name"`
	Direction     bool    `json:"direction"`
	Syringe       float64 `json:"syringe"`
	Used          bool    `json:"used"`
	VolumeTargetW float64 `json:"volumeTargetW"`
	VolumeW       float64 `json:"volumeW"`
	Rate          float64 `json:"rate"`
	Stalled       bool    `json:"stalled"`
	Force         float64 `json:"force"`
	initialized   bool
}

type requestBody struct {
	Par   interface{} `json:"par"`
	Pump  interface{} `json:"pump"`
	Value interface{} `json:"value"`
}

type response struct {
	Success int         `json:"success"`
	Data    interface{} `json:"data"`
}

func (p pump) Invoke(invoke clientInvocation, data interface{}) error {
	const pumpBaseAddr = "http://192.168.1.100:8764"

	var (
		endpoint    string
		payloadData interface{}
	)

	switch invoke {
	case PumpSetTargetVolume:
		payloadData = requestBody{Par: "volumeTargetW", Pump: p.PumpID, Value: data}
		endpoint = pumpBaseAddr + "/update"
	case PumpReset:
		payloadData = requestBody{Pump: p.PumpID}
		endpoint = pumpBaseAddr + "/update"
	case PumpToggleWithdrawInfuse:
		payloadData = requestBody{Par: "direction", Pump: p.PumpID, Value: data}
		endpoint = pumpBaseAddr + "/update"
	case PumpSetVolume:
		payloadData = requestBody{Par: "rate", Pump: p.PumpID, Value: data}
		endpoint = pumpBaseAddr + "/update"
	case PumpToggle:
		payloadData = requestBody{Par: "status", Pump: p.PumpID, Value: data}
		endpoint = pumpBaseAddr + "/update"
	case PumpRefresh:
		payloadData = requestBody{Par: "status", Pump: p.PumpID, Value: data}
		endpoint = pumpBaseAddr + "/refresh"
	case PumpPurge:
		// TODO : collect data in GMC
	default:
		panic("incorrect invoke operation of pump client")
	}

	var httpResponse *http.Response
	if err := makePost(endpoint, "application/json", payloadData, httpResponse); err != nil {
		return err
	}

	var responseData response

	switch invoke {
	case PumpRefresh:
		var doubleJson dataPack
		if err := json.NewDecoder(httpResponse.Body).Decode(&doubleJson); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(doubleJson.DataEscaped), &p); err != nil {
			return err
		}
	default:
		if err := json.NewDecoder(httpResponse.Body).Decode(&responseData); err != nil {
			return err
		}

		if responseData.Success == 1 {
			return errors.New("camera device failed to process the request")
		}
	}
	return nil
}

func NewPump(pumpID int) pump {
	if pumpID > 0 && pumpID < 4 {
		newPump := pump{}
		newPump.PumpID = float64(pumpID)
		return newPump
	}
	newPump := pump{}
	newPump.PumpID = -1
	return newPump

}

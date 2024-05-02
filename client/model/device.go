package model

import (
	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
)

type GetDeviceResponse struct {
	Device            *api.Device
	DeviceStatus      *api.DeviceStatus
	ClassEnabled      common.DeviceClass
	DeviceActivation  *api.DeviceActivation
	JoinServerContext *common.JoinServerContext
	DeviceKeys        *api.DeviceKeys
}

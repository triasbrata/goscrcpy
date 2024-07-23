package clientManager

import "github.com/triasbrata/gadb"

type cm struct {
	client           gadb.Client
	avaliableDevices []gadb.Device
}

func (c *cm) RefreshDevices() error {
	devices, err := c.client.DeviceList()
	if err != nil {
		return err
	}
	c.avaliableDevices = devices
	return nil
}

func New() (*cm, error) {
	client, err := gadb.NewClient()
	if err != nil {
		return nil, err
	}
	return &cm{
		client: client,
	}, nil
}

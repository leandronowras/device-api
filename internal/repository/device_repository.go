package repository

import (
	"context"

	"github.com/leandronowras/device-api/internal/device"
)

type DeviceRepository interface {
	Save(ctx context.Context, d *device.Device) (*device.Device, error)
	FindByID(ctx context.Context, id string) (*device.Device, error)
	FindAll(ctx context.Context, brand, state *string) ([]*device.Device, error)
	Update(ctx context.Context, d *device.Device) (*device.Device, error)
	Delete(ctx context.Context, id string) error
}

package create

import (
	"fmt"

	"github.com/docker/go-units"
	"github.com/rancher/rio/types/client/rio/v1beta1"
)

func populateMemory(c *Create, service *client.Service) error {
	var err error

	if c.M_Memory != "" {
		service.MemoryBytes, err = units.RAMInBytes(c.M_Memory)
		if err != nil {
			return fmt.Errorf("failed to parse --memory: %v", err)
		}
	}

	if c.MemoryReservation != "" {
		service.MemoryReservationBytes, err = units.RAMInBytes(c.MemoryReservation)
		if err != nil {
			return fmt.Errorf("failed to parse --memory-reservation: %v", err)
		}
	}

	return nil
}

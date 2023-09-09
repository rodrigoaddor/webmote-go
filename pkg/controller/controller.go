package controller

import (
	"github.com/71/stadiacontroller"
	"log"
)

type Controller struct {
	controller *stadiacontroller.Xbox360Controller
}

var emulator *stadiacontroller.Emulator

func init() {
	var err error
	emulator, err = stadiacontroller.NewEmulator(func(_ stadiacontroller.Vibration) {})
	if err != nil {
		log.Fatalf("failed to init emulator: %v", err)
	}
}

func NewController() (*Controller, error) {
	x360controller, err := emulator.CreateXbox360Controller()
	if err != nil {
		_ = x360controller.Close()
		return nil, err
	}

	if err = x360controller.Connect(); err != nil {
		_ = x360controller.Close()
		return nil, err
	}

	controller := &Controller{
		controller: x360controller,
	}

	return controller, nil
}

func (c *Controller) Close() error {
	return c.controller.Close()
}

func (c *Controller) SetButton(key int, value bool) error {
	report := stadiacontroller.NewXbox360ControllerReport()
	report.MaybeSetButton(key, value)
	return c.controller.Send(&report)
}

func (c *Controller) SetLeftAxis(x, y int16) error {
	report := stadiacontroller.NewXbox360ControllerReport()
	report.SetLeftThumb(x, y)
	return c.controller.Send(&report)
}

func (c *Controller) SetRightAxis(x, y int16) error {
	report := stadiacontroller.NewXbox360ControllerReport()
	report.SetRightThumb(x, y)
	return c.controller.Send(&report)
}

func (c *Controller) SetLeftTrigger(value byte) error {
	report := stadiacontroller.NewXbox360ControllerReport()
	report.SetLeftTrigger(value)
	return c.controller.Send(&report)
}

func (c *Controller) SetRightTrigger(value byte) error {
	report := stadiacontroller.NewXbox360ControllerReport()
	report.SetRightTrigger(value)
	return c.controller.Send(&report)
}

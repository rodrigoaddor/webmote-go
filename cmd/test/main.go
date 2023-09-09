package main

import (
	"github.com/71/stadiacontroller"
	"log"
	"time"
)

func main() {
	emulator, err := stadiacontroller.NewEmulator(func(vibration stadiacontroller.Vibration) {})
	if err != nil {
		log.Fatalln("failed to start ViGEm client: %w", err)
	}

	defer func() { _ = emulator.Close() }()

	controller, err := emulator.CreateXbox360Controller()
	if err != nil {
		log.Fatalf("failed to create controller: %v", err)
	}

	defer func() { _ = controller.Close() }()

	if err = controller.Connect(); err != nil {
		log.Fatalf("unable to connect to emulated Xbox 360 controller: %v", err)
	}

	isPressed := false

	for {
		isPressed = !isPressed
		report := stadiacontroller.NewXbox360ControllerReport()
		report.MaybeSetButton(stadiacontroller.Xbox360ControllerButtonA, isPressed)
		if err := controller.Send(&report); err != nil {
			log.Fatalf("failed to send report: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}

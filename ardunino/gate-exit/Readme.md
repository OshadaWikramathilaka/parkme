# Gate Exit Module

This module controls the exit gate of the parking lot. It manages the hardware and logic for allowing vehicles to leave.

## Wiring Instructions

### ESP32-CAM (AI Thinker)
- VCC → 5V
- GND → GND
- SERVO (Gate Servo) → GPIO2 (may require external pullup resistor)
- LED (Status) → GPIO4 (onboard LED)

### NodeMCU
- VCC → 5V
- GND → GND
- TRIG (Ultrasonic Trigger) → D1 (GPIO5)
- ECHO (Ultrasonic Echo) → D2 (GPIO4)
- SERVO (Gate Servo) → D3 (GPIO0)

## Purpose
Controls the exit barrier, reads ultrasonic sensor input, and communicates with the central system to allow vehicles to exit the parking area. The module supports both ESP32-CAM and NodeMCU variants, with pin assignments as specified above for correct hardware setup.
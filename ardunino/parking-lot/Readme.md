# Parking Lot NodeMCU System

This project uses a NodeMCU (ESP8266) to monitor parking spot availability using TCRT5000 IR sensors and displays the status on a 1604A LCD with I2C interface. The system also sends real-time status updates to a backend server.

## Hardware Connections

### TCRT5000 IR Sensors (for 3 parking spots)
- **VCC** → 5V (NodeMCU)
- **GND** → GND (NodeMCU)
- **D0 (digital output)** → NodeMCU pins D5 (GPIO14), D6 (GPIO12), D7 (GPIO13)
- **A0 (analog output)** → Not used in current implementation

### 1604A LCD with I2C Converter
- **SDA** → NodeMCU D2 (GPIO4)
- **SCL** → NodeMCU D1 (GPIO5)
- **VCC** → 5V (NodeMCU)
- **GND** → GND (NodeMCU)
- **I2C Address**: 0x27 (default for most modules)

## Features
- Displays the number of available parking spots on the LCD.
- Sends each spot's occupancy status to a backend server via HTTP POST.
- WiFi credentials, backend server address, and API key are configured in the code.

## Usage
1. Connect the sensors and LCD as described above.
2. Flash the `parking-lot.ino` sketch to your NodeMCU.
3. Power on the NodeMCU. The LCD will show the system status and available spots.
4. The system will automatically connect to WiFi and start sending updates to the backend.

## Notes
- Adjust the `SENSOR_THRESHOLD` in the code if needed for your TCRT5000 calibration.
- Ensure your backend server is reachable from the NodeMCU's network.
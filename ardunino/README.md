# ParkMe Arduino

This project creates a gate entry system using a NodeMCU with ultrasonic sensor and an ESP32-CAM for taking photos. The system controls a servo motor to open and close a gate for authorized vehicles.

## Components
- NodeMCU (ESP8266)
- ESP32-CAM
- HC-SR04 Ultrasonic Sensor
- SG90 Servo Motor (or similar)

## Wiring
### NodeMCU Connections
- HC-SR04 Trigger pin: D1
- HC-SR04 Echo pin: D2
- Servo motor signal pin: D3
- Serial communication to ESP32-CAM:
  - NodeMCU RX pin → ESP32-CAM U0TXD pin
  - NodeMCU TX pin → ESP32-CAM U0RXD pin
  - NodeMCU GND → ESP32-CAM GND
- Optional debug output: D6 (RX), D7 (TX)

### ESP32-CAM Connections
- Serial communication to NodeMCU:
  - ESP32-CAM U0TXD pin → NodeMCU RX pin
  - ESP32-CAM U0RXD pin → NodeMCU TX pin
  - ESP32-CAM GND → NodeMCU GND
- Optional debug output: pins 12(RX), 13(TX)

### Power Supply
- ESP32-CAM: Use an external 5V power supply capable of providing at least 500mA
- NodeMCU: Can be powered via USB or an external 5V source

**IMPORTANT NOTE:** When flashing the ESP32-CAM, you'll need to disconnect the RX/TX connections to the NodeMCU. After flashing, reconnect them.

## Setup Instructions

1. Edit the ESP32-CAM code (`gate_enter_esp32cam.ino`):
   - Enter your WiFi credentials (`ssid` and `password`)
   - Enter your server details:
     - `serverHost`: Your backend server hostname (e.g., "your-server.com" or "192.168.1.100")
     - `serverPath`: The API endpoint path (default: "/api/arduino/gate/enter/upload")
     - `serverPort`: The port number (default: 80 for HTTP, 443 for HTTPS)
   - Enter your API key: `apiKey`
   - Enter your location ID: `locationId`

2. Flash the ESP32-CAM with `gate_enter_esp32cam.ino`:
   - Disconnect RX/TX connections to NodeMCU first
   - Connect the ESP32-CAM to your computer via an FTDI programmer or ESP32-CAM development board
   - Select "AI Thinker ESP32-CAM" in the Arduino IDE board manager
   - Upload the code
   - Reconnect the RX/TX connections to NodeMCU after flashing is complete

3. Flash the NodeMCU with `gate_enter_nodemcu.ino`:
   - Disconnect RX/TX connections to ESP32-CAM first
   - Connect the NodeMCU to your computer via USB
   - Select your NodeMCU board in Arduino IDE (e.g., "NodeMCU 1.0 (ESP-12E Module)")
   - Upload the code
   - Reconnect the RX/TX connections to ESP32-CAM after flashing is complete

4. Connect the servo motor to the gate mechanism

## How It Works
1. The NodeMCU continuously monitors distance using the ultrasonic sensor
2. When an object is detected at approximately 4.5cm:
   - The NodeMCU sends a "TAKE_PHOTO" command to the ESP32-CAM
   - The ESP32-CAM captures an image
   - The ESP32-CAM sends the image to your backend server
   - The backend server processes the image and responds with vehicle information
   - The ESP32-CAM parses the response and forwards relevant information to the NodeMCU
   - If access is granted, the NodeMCU opens the gate using the servo motor
   - The gate automatically closes after 10 seconds (configurable)

## Debugging
Both devices support debug output:
- ESP32-CAM: Connect a USB-to-TTL converter to pins 12(RX), 13(TX) to view detailed debug messages
- NodeMCU: Connect a USB-to-TTL converter to pins D6(RX), D7(TX) for debug output

## Backend Server API
Your backend server should:
1. Accept a POST request to the configured endpoint (default: `/api/arduino/gate/enter/upload`)
2. Expect form-data with the following fields:
   - `image`: The JPEG image file from the camera
   - `location_id`: Your location identifier
3. Require a header `X-API-Key` with your API key
4. Process the image (e.g., license plate recognition, facial recognition)
5. Return a JSON response with:
   - `success`: boolean indicating whether access is granted
   - `message`: description of the result
   - `data`: object containing booking, vehicle, and owner information

## Response Format Example
```json
{
    "success": true,
    "message": "Vehicle processed and booking created successfully",
    "data": {
        "booking": {
            "id": "000000000000000000000000",
            "vehicleId": "677cbb9817a0737c5329f968",
            "userId": "6778eff38d6d57cdd58075ff",
            "locationId": "678bd59982059b655da1098e",
            "startTime": "2025-03-19T19:08:17.2840851+05:30",
            "status": "pending",
            "spotNumber": "A1",
            "bookingType": "on_site",
            "createdAt": "2025-03-19T19:08:17.861994+05:30",
            "updatedAt": "2025-03-19T19:08:17.861994+05:30"
        },
        "owner": {
            "id": "6778eff38d6d57cdd58075ff",
            "name": "Leo Felcianas",
            "email": "example@email.com",
            "role": "admin",
            "status": "active"
        },
        "vehicle": {
            "id": "677cbb9817a0737c5329f968",
            "plate_number": "PH3392",
            "brand": "Honda",
            "model": "Civic",
            "owner": "6778eff38d6d57cdd58075ff"
        }
    }
}
```

## Customization
- Adjust `GATE_CLOSED_POSITION` and `GATE_OPEN_POSITION` in the NodeMCU code to match your servo motor's requirements
- Modify `GATE_OPEN_TIME` to change how long the gate stays open (default: 10 seconds)

## Troubleshooting
- **Connection Failed Error**: Check that your `serverHost` is correct and doesn't include "http://" or the path
- **Timeout Errors**: Make sure your server is accessible from the ESP32-CAM's network
- **Parsing Errors**: Check that your server is returning properly formatted JSON
- **Camera Errors**: If the camera fails to initialize, check that you're using the correct board selection
- **Communication Issues**: If the NodeMCU and ESP32-CAM aren't communicating, ensure the RX/TX connections are correct (NodeMCU RX → ESP32-CAM TX, NodeMCU TX → ESP32-CAM RX)
- Make sure both devices are powered adequately (ESP32-CAM needs at least 5V with stable power)
- If the ESP32-CAM keeps resetting, make sure your power supply can deliver enough current
- Verify your API key and location ID are correctly configured
- Test the servo motor independently to ensure it has sufficient power
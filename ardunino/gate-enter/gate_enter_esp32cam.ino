#include "esp_camera.h"
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include <ESP32Servo.h>

// WiFi credentials
const char* ssid = "Dialog 4G 405";
const char* password = "64CBCCC7";

// Backend server settings - separate host and path
const char* serverHost = "192.168.8.194"; // Replace with your actual server hostname
const char* serverPath = "/api/arduino/gate/enter/upload";
const int serverPort = 8080; // Default HTTP port, use 443 for HTTPS
const char* apiKey = "533a269d-e866-43be-93af-e65a43905762";
const char* locationId = "678bd59982059b655da1098e";

// Pin definitions for CAMERA_MODEL_AI_THINKER
#define PWDN_GPIO_NUM     32
#define RESET_GPIO_NUM    -1
#define XCLK_GPIO_NUM      0
#define SIOD_GPIO_NUM     26
#define SIOC_GPIO_NUM     27
#define Y9_GPIO_NUM       35
#define Y8_GPIO_NUM       34
#define Y7_GPIO_NUM       39
#define Y6_GPIO_NUM       36
#define Y5_GPIO_NUM       21
#define Y4_GPIO_NUM       19
#define Y3_GPIO_NUM       18
#define Y2_GPIO_NUM        5
#define VSYNC_GPIO_NUM    25
#define HREF_GPIO_NUM     23
#define PCLK_GPIO_NUM     22

// LED pin for status indication
#define LED_PIN 4  // GPIO4 is the onboard LED on many ESP32-CAM boards

// GPIO pin for Servo control
#define SERVO_PIN 2  // GPIO2 is available on ESP32-CAM (may need external pullup resistor)

// Gate control parameters
// Servo gateServo; // REMOVED: No longer controlling servo from ESP32-CAM
// const int GATE_CLOSED_POSITION = 0;    // Angle for closed gate (adjust as needed)
// const int GATE_OPEN_POSITION = 90;     // Angle for open gate (adjust as needed)
// const unsigned long GATE_OPEN_TIME = 10000; // Time to keep gate open in ms (10 seconds)
// unsigned long gateOpenTime = 0;
// bool gateIsOpen = false;

// Ultrasonic sensor pins (if connected directly to ESP32-CAM)
#define TRIGGER_PIN 12  // Can use GPIO12 on ESP32-CAM
#define ECHO_PIN 13     // Can use GPIO13 on ESP32-CAM
bool processingPicture = false;
long duration;
float distance;

void setup() {
  // Start serial communication
  Serial.begin(115200); // Use default UART for both debug and NodeMCU communication
  
  // Initialize status LED
  pinMode(LED_PIN, OUTPUT);
  digitalWrite(LED_PIN, LOW); // Turn on LED (LOW is ON for most ESP32-CAM boards)
  
  // Setup ultrasonic sensor (if connected directly to ESP32-CAM)
  pinMode(TRIGGER_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);
  digitalWrite(TRIGGER_PIN, LOW);
  
  // Initial startup message
  delay(1000); // Wait for serial to be ready
  Serial.println("\n\n");
  Serial.println("====================================");
  Serial.println("ESP32-CAM Gate System Initializing");
  Serial.println("====================================");
  Serial.println("Firmware compiled: " __DATE__ " " __TIME__);
  Serial.printf("ESP32 SDK Version: %s\n", ESP.getSdkVersion());
  Serial.printf("ESP32 CPU Freq: %d MHz\n", ESP.getCpuFreqMHz());
  Serial.printf("ESP32 Flash Size: %d bytes\n", ESP.getFlashChipSize());
  Serial.printf("Free Heap: %d bytes\n", ESP.getFreeHeap());
  
  // Setup servo control (allocate timer before camera initialization)
  // ESP32PWM::allocateTimer(0);  // Allocate timer 0
  // gateServo.setPeriodHertz(50); // Standard 50Hz servo
  // gateServo.attach(SERVO_PIN, 500, 2400); // Attach with min/max pulse width
  // gateServo.write(GATE_CLOSED_POSITION); // Start with gate closed
  // Serial.println("Servo initialized");
  
  // Initialize camera
  Serial.println("Initializing camera...");
  camera_config_t config;
  config.ledc_channel = LEDC_CHANNEL_0;
  config.ledc_timer = LEDC_TIMER_0;
  config.pin_d0 = Y2_GPIO_NUM;
  config.pin_d1 = Y3_GPIO_NUM;
  config.pin_d2 = Y4_GPIO_NUM;
  config.pin_d3 = Y5_GPIO_NUM;
  config.pin_d4 = Y6_GPIO_NUM;
  config.pin_d5 = Y7_GPIO_NUM;
  config.pin_d6 = Y8_GPIO_NUM;
  config.pin_d7 = Y9_GPIO_NUM;
  config.pin_xclk = XCLK_GPIO_NUM;
  config.pin_pclk = PCLK_GPIO_NUM;
  config.pin_vsync = VSYNC_GPIO_NUM;
  config.pin_href = HREF_GPIO_NUM;
  config.pin_sscb_sda = SIOD_GPIO_NUM;
  config.pin_sscb_scl = SIOC_GPIO_NUM;
  config.pin_pwdn = PWDN_GPIO_NUM;
  config.pin_reset = RESET_GPIO_NUM;
  config.xclk_freq_hz = 20000000;
  config.pixel_format = PIXFORMAT_JPEG;
  
  // Check for PSRAM
  bool hasPSRAM = psramFound();
  Serial.printf("PSRAM found: %s\n", hasPSRAM ? "YES" : "NO");
  
  // Use a much smaller image size to avoid memory issues
  if (hasPSRAM) {
    config.frame_size = FRAMESIZE_VGA; // Reduced from UXGA to VGA (640x480)
    config.jpeg_quality = 12; // Lower quality (higher number = lower quality)
    config.fb_count = 1;
    Serial.println("Camera config: VGA resolution (640x480)");
  } else {
    config.frame_size = FRAMESIZE_QVGA; // Even smaller (320x240)
    config.jpeg_quality = 12;
    config.fb_count = 1;
    Serial.println("Camera config: QVGA resolution (320x240)");
  }

  // Initialize the camera
  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK) {
    Serial.printf("Camera initialization FAILED with error 0x%x\n", err);
    // Blink LED rapidly to indicate camera error
    for (int i = 0; i < 10; i++) {
      digitalWrite(LED_PIN, HIGH);
      delay(100);
      digitalWrite(LED_PIN, LOW);
      delay(100);
    }
    Serial.println("System halted due to camera initialization failure");
    while (1) {
      // Blink SOS pattern
      for (int i = 0; i < 3; i++) { // Short blinks (S)
        digitalWrite(LED_PIN, LOW); delay(200); digitalWrite(LED_PIN, HIGH); delay(200);
      }
      delay(300);
      for (int i = 0; i < 3; i++) { // Long blinks (O)
        digitalWrite(LED_PIN, LOW); delay(600); digitalWrite(LED_PIN, HIGH); delay(200);
      }
      delay(300);
      for (int i = 0; i < 3; i++) { // Short blinks (S)
        digitalWrite(LED_PIN, LOW); delay(200); digitalWrite(LED_PIN, HIGH); delay(200);
      }
      delay(1000);
    }
  }
  Serial.println("Camera initialized successfully!");
  
  // Test camera by taking a snapshot
  Serial.println("Taking test image...");
  camera_fb_t * fb = esp_camera_fb_get();
  if (!fb) {
    Serial.println("Test image capture failed!");
  } else {
    Serial.printf("Test image captured successfully! Size: %d bytes\n", fb->len);
    esp_camera_fb_return(fb);
  }

  // Connect to WiFi
  Serial.printf("Connecting to WiFi network: %s\n", ssid);
  WiFi.begin(ssid, password);
  
  // Blink LED while connecting to WiFi
  int wifiAttempts = 0;
  while (WiFi.status() != WL_CONNECTED) {
    digitalWrite(LED_PIN, !digitalRead(LED_PIN)); // Toggle LED
    delay(500);
    Serial.print(".");
    wifiAttempts++;
    if (wifiAttempts > 20) {
      Serial.println("\nWiFi connection timeout! Will continue but WiFi is not connected.");
      break;
    }
  }
  
  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("\nWiFi connected successfully!");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());
    // Blink LED twice to indicate WiFi connected
    for (int i = 0; i < 2; i++) {
      digitalWrite(LED_PIN, HIGH);
      delay(100);
      digitalWrite(LED_PIN, LOW);
      delay(100);
    }
  }
  
  // Turn off LED to indicate all systems ready
  digitalWrite(LED_PIN, HIGH); // LED off (HIGH is OFF for most ESP32-CAM boards)
  
  // All systems ready
  Serial.println("\n====================================");
  Serial.println("ESP32-CAM Gate System READY");
  Serial.println("====================================");
  Serial.println("Waiting for vehicle detection...");
  
  // Send ready signal to serial
  Serial.println("ESP32CAM:READY");
}

void loop() {
  // Process any commands from serial (if still connected to NodeMCU)
  if (Serial.available()) {
    String command = Serial.readStringUntil('\n');
    command.trim();
    
    Serial.printf("Command received: %s\n", command.c_str());
    
    if (command == "TAKE_PHOTO") {
      Serial.println("Photo request received");
      // Blink LED once
      digitalWrite(LED_PIN, LOW);
      delay(100);
      digitalWrite(LED_PIN, HIGH);
      
      takePhotoAndSend();
    }
  }
  
  // Check ultrasonic sensor for vehicle detection
  checkUltrasonicSensor();
  
  // Check if it's time to close the gate
  // if (gateIsOpen && (millis() - gateOpenTime > GATE_OPEN_TIME)) {
  //   closeGate();
  // }
  
  delay(100);
}

void checkUltrasonicSensor() {
  // Clear trigger pin
  digitalWrite(TRIGGER_PIN, LOW);
  delayMicroseconds(2);
  
  // Send trigger pulse
  digitalWrite(TRIGGER_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIGGER_PIN, LOW);
  
  // Get duration
  duration = pulseIn(ECHO_PIN, HIGH);
  
  // Calculate distance
  distance = (duration * 0.0343) / 2;
  
  // Only print distance every 2 seconds to reduce serial output
  static unsigned long lastPrintTime = 0;
  if (millis() - lastPrintTime > 2000) {
    Serial.print("Distance: ");
    Serial.print(distance);
    Serial.println(" cm");
    lastPrintTime = millis();
  }
  
  // Check if object is at approximately 4.5cm (with 0.5cm tolerance)
  // and we're not currently processing a picture
  if (distance >= 4.0 && distance <= 5.0 && !processingPicture) {
    Serial.println("Detected vehicle at recognition distance!");
    digitalWrite(LED_PIN, LOW);  // Turn on LED
    processingPicture = true;
    takePhotoAndSend();
  }
}

void takePhotoAndSend() {
  Serial.println("Taking photo...");
  
  // Take a photo
  camera_fb_t * fb = esp_camera_fb_get();
  if (!fb) {
    Serial.println("Camera capture failed");
    digitalWrite(LED_PIN, HIGH);  // Turn off LED
    processingPicture = false;
    return;
  }
  
  Serial.println("Photo taken! Sending to server...");
  Serial.printf("Image size: %d bytes\n", fb->len);
  
  // Check WiFi connection
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi disconnected. Reconnecting...");
    WiFi.reconnect();
    int attempts = 0;
    while (WiFi.status() != WL_CONNECTED && attempts < 20) {
      delay(500);
      Serial.print(".");
      attempts++;
    }
    if (WiFi.status() != WL_CONNECTED) {
      Serial.println("WiFi reconnection failed");
      esp_camera_fb_return(fb);
      digitalWrite(LED_PIN, HIGH);  // Turn off LED
      processingPicture = false;
      return;
    }
    Serial.println("WiFi reconnected");
  }
  
  Serial.printf("Connecting to server: %s:%d%s\n", serverHost, serverPort, serverPath);
  
  // Create multipart form data manually in memory
  // Prepare headers
  String boundary = "ParkMeArduinoBoundary";
  String head = "--" + boundary + "\r\n";
  head += "Content-Disposition: form-data; name=\"image\"; filename=\"image.jpg\"\r\n";
  head += "Content-Type: image/jpeg\r\n\r\n";
  
  String tail = "\r\n--" + boundary + "\r\n";
  tail += "Content-Disposition: form-data; name=\"location_id\"\r\n\r\n";
  tail += locationId;
  tail += "\r\n--" + boundary + "--\r\n";
  
  // Calculate content length
  size_t contentLength = head.length() + fb->len + tail.length();
  
  // Use direct WiFi client for more control - create TCP connection manually
  WiFiClient client;
  
  Serial.println("Creating connection to server...");
  
  // Set connection timeout
  client.setTimeout(10000); // 10 seconds timeout
  
  if (!client.connect(serverHost, serverPort)) {
    Serial.println("Connection failed! Check your server address and port.");
    esp_camera_fb_return(fb);
    digitalWrite(LED_PIN, HIGH);  // Turn off LED
    processingPicture = false;
    return;
  }
  
  Serial.println("Connection established. Sending HTTP request headers...");
  
  // Send the HTTP POST request headers manually
  client.print("POST ");
  client.print(serverPath);
  client.println(" HTTP/1.1");
  client.print("Host: ");
  client.println(serverHost);
  client.print("X-API-Key: ");
  client.println(apiKey);
  client.print("Content-Type: multipart/form-data; boundary=");
  client.println(boundary);
  client.print("Content-Length: ");
  client.println(contentLength);
  client.println("Connection: close");
  client.println(); // End of headers
  
  // Check if we can still write to client
  if (!client.connected()) {
    Serial.println("Connection lost after sending headers!");
    client.stop();
    esp_camera_fb_return(fb);
    digitalWrite(LED_PIN, HIGH);  // Turn off LED
    processingPicture = false;
    return;
  }
  
  // Send the multipart form data
  // Send boundary and file header
  client.print(head);
  
  // Send file data in chunks
  Serial.println("Sending image data in chunks...");
  const size_t bufferSize = 1024;
  uint8_t *buffer = new uint8_t[bufferSize];
  if (!buffer) {
    Serial.println("Failed to allocate buffer");
    client.stop();
    esp_camera_fb_return(fb);
    digitalWrite(LED_PIN, HIGH);  // Turn off LED
    processingPicture = false;
    return;
  }
  
  // Send image data
  bool success = true;
  for (size_t i = 0; i < fb->len; i += bufferSize) {
    // Check connection status periodically
    if (!client.connected()) {
      Serial.println("Connection lost while sending image data!");
      success = false;
      break;
    }
    
    size_t chunkSize = (i + bufferSize < fb->len) ? bufferSize : (fb->len - i);
    
    // Copy chunk to buffer
    memcpy(buffer, fb->buf + i, chunkSize);
    
    // Write chunk
    if (client.write(buffer, chunkSize) != chunkSize) {
      Serial.printf("Failed to send data chunk at position %d\n", i);
      success = false;
      break;
    }
    
    // Print progress every 10KB
    if (i % 10240 == 0) {
      Serial.printf("Sent %d bytes (%.1f%%)\n", i, (100.0 * i) / fb->len);
    }
  }
  
  delete[] buffer;
  
  if (!success) {
    Serial.println("Failed to send image data");
    client.stop();
    esp_camera_fb_return(fb);
    digitalWrite(LED_PIN, HIGH);  // Turn off LED
    processingPicture = false;
    return;
  }
  
  // Send tail with location_id
  client.print(tail);
  
  // Free the camera frame buffer
  esp_camera_fb_return(fb);
  
  Serial.println("HTTP request sent, waiting for response...");

  // Read the HTTP response body directly (skip header parsing)
  String responseBody = "";
  unsigned long responseStart = millis();
  bool headersEnded = false;
  while (client.connected() && millis() - responseStart < 10000) {
    while (client.available()) {
      String line = client.readStringUntil('\n');
      line.trim();
      if (!headersEnded) {
        if (line.length() == 0) {
          headersEnded = true;
        }
        continue;
      }
      // After headers, accumulate JSON body
      responseBody += line;
    }
  }
  client.stop();
  Serial.print("Response body: ");
  Serial.println(responseBody);

  // Parse JSON and check 'success' field
  bool accessGranted = false;
  if (responseBody.length() > 0) {
    StaticJsonDocument<256> doc;
    DeserializationError error = deserializeJson(doc, responseBody);
    if (!error) {
      if (doc["success"] == true) {
        accessGranted = true;
      }
    } else {
      Serial.print("JSON parse error: ");
      Serial.println(error.c_str());
    }
  }

  // Handle gate control based on JSON response
  if (accessGranted) {
    Serial.println("ACCESS GRANTED - Notifying NodeMCU");
    Serial.println("ACCESS:GRANTED"); // Send signal to NodeMCU
  } else {
    Serial.println("ACCESS DENIED - Not opening gate");
    Serial.println("ACCESS:DENIED"); // Send denial to NodeMCU
  }
  digitalWrite(LED_PIN, HIGH);  // Turn off LED
  processingPicture = false;
  Serial.println("Request completed");
}

// Function to open the gate
// void openGate() {
//   Serial.println("Opening gate...");
//   Serial.print("DEBUG: Writing servo position to ");
//   Serial.println(GATE_OPEN_POSITION);
//   gateServo.write(GATE_OPEN_POSITION);
//   Serial.println("DEBUG: Servo write command sent for opening gate");
//   gateIsOpen = true;
//   gateOpenTime = millis(); // Record the time when gate was opened
//   // Flash LED to indicate gate opening
//   for (int i = 0; i < 3; i++) {
//     digitalWrite(LED_PIN, LOW);  // LED on
//     delay(100);
//     digitalWrite(LED_PIN, HIGH); // LED off
//     delay(100);
//   }
// }

// Function to close the gate
// void closeGate() {
//   Serial.println("Closing gate...");
//   Serial.print("DEBUG: Writing servo position to ");
//   Serial.println(GATE_CLOSED_POSITION);
//   gateServo.write(GATE_CLOSED_POSITION);
//   Serial.println("DEBUG: Servo write command sent for closing gate");
//   gateIsOpen = false;
//   // Flash LED to indicate gate closing
//   for (int i = 0; i < 2; i++) {
//     digitalWrite(LED_PIN, LOW);  // LED on
//     delay(200);
//     digitalWrite(LED_PIN, HIGH); // LED off
//     delay(200);
//   }
// }

#include <Servo.h>

#define TRIGGER_PIN 5 //D1
#define ECHO_PIN    4 //D2
#define SERVO_PIN   0  // Servo control pin (GPIO0) //D3

// Serial communication will be used for both debugging and ESP32-CAM communication
// Note: This will mix debug output with communication to the ESP32-CAM

long duration;
float distance;
bool processingPicture = false;
unsigned long pictureRequestTime = 0;
const unsigned long PICTURE_TIMEOUT = 30000; // 30 seconds timeout for picture processing

// Gate control parameters
Servo gateServo;
const int GATE_CLOSED_POSITION = 0;    // Angle for closed gate (adjust as needed)
const int GATE_OPEN_POSITION = 90;     // Angle for open gate (adjust as needed)
const unsigned long GATE_OPEN_TIME = 10000; // Time to keep gate open in ms (10 seconds)
unsigned long gateOpenTime = 0;
bool gateIsOpen = false;

void setup() {
  // Begin serial communication with ESP32-CAM and for debug
  Serial.begin(115200);
  
  pinMode(TRIGGER_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);
  digitalWrite(TRIGGER_PIN, LOW);
  
  // Initialize the servo motor
  gateServo.attach(SERVO_PIN);
  gateServo.write(GATE_CLOSED_POSITION); // Start with gate closed
  
  // Wait for the serial connection
  delay(1000);
  Serial.println("NodeMCU Gate Entry System Started");
  
  // Wait for ESP32-CAM to be ready
  Serial.println("Waiting for ESP32-CAM to initialize...");
  bool espReady = false;
  unsigned long startTime = millis();
  while (!espReady && millis() - startTime < 20000) { // Wait up to 20 seconds
    if (Serial.available()) {
      String message = Serial.readStringUntil('\n');
      message.trim();
      
      if (message == "ESP32CAM:READY") {
        espReady = true;
        Serial.println("ESP32-CAM is ready!");
      } else {
        Serial.print("Received from ESP32-CAM: ");
        Serial.println(message);
      }
    }
    delay(100);
  }
  
  if (!espReady) {
    Serial.println("ESP32-CAM did not respond in time. Continuing anyway...");
  }
  
  Serial.println("System is ready!");
}

void loop() {
  // Process incoming messages from ESP32-CAM
  if (Serial.available()) {
    String response = Serial.readStringUntil('\n');
    response.trim();
    
    Serial.print("Received: ");
    Serial.println(response);
    
    // Reset processing flag for ANY response received from ESP32-CAM
    if (processingPicture) {
      processingPicture = false;
      Serial.println("Picture processing completed.");
    }
    
    if (response.startsWith("ACCESS GRANTED - Notifying NodeMCU")) {
      // Access granted - open the gate
      Serial.println("Access granted!");
      
      // Open the gate
      openGate();
    } 
    else if (response.startsWith("ACCESS:DENIED")) {
      String reason = "Access denied";
      if (response.length() > 14) {
        reason = response.substring(14); // Get reason if available
      }
      Serial.println(reason);
    } 
    else if (response.startsWith("ERROR:")) {
      String error = response.substring(6); // Remove "ERROR:" prefix
      Serial.print("Error: ");
      Serial.println(error);
    }
  }
  
  // Check for picture processing timeout
  if (processingPicture && (millis() - pictureRequestTime > PICTURE_TIMEOUT)) {
    Serial.println("Picture processing timeout. Ready for new detection.");
    processingPicture = false;
  }
  
  // Check if it's time to close the gate
  if (gateIsOpen && (millis() - gateOpenTime > GATE_OPEN_TIME)) {
    closeGate();
  }
  
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
  
  Serial.print("Distance: ");
  Serial.print(distance);
  Serial.println(" cm");
  
  // Check if object is at approximately 4.5cm (with 0.5cm tolerance)
  // and we're not currently processing a picture
  if (distance >= 4.0 && distance <= 5.0 && !processingPicture) {
    Serial.println("Detected object at recognition distance!");
    Serial.println("TAKE_PHOTO");  // Send command to ESP32-CAM
    processingPicture = true;
    pictureRequestTime = millis(); // Start timeout timer
    Serial.println("Picture requested, waiting for response...");
  }
  
  delay(500);
}

// Function to open the gate
void openGate() {
  Serial.println("Opening gate...");
  gateServo.write(GATE_OPEN_POSITION);
  gateIsOpen = true;
  gateOpenTime = millis(); // Record the time when gate was opened
}

// Function to close the gate
void closeGate() {
  Serial.println("Closing gate...");
  gateServo.write(GATE_CLOSED_POSITION);
  gateIsOpen = false;
}
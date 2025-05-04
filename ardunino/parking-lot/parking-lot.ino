#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <ArduinoJson.h>
#include <Wire.h>
#include <hd44780.h>
#include <hd44780ioClass/hd44780_I2Cexp.h> // Include the I2C expander I/O class

hd44780_I2Cexp lcd; // Create hd44780 object for I2C backpack, auto-detects address

// TCRT5000 sensor settings
#define SENSOR_THRESHOLD 500 // Adjust based on your calibration

// WiFi credentials
const char* ssid = "Dialog 4G 405";
const char* password = "64CBCCC7";

// Backend server settings
const char* serverHost = "192.168.8.194";
const char* serverPath = "/api/arduino/spot/status";
const int serverPort = 8080;
const char* locationId = "678bd59982059b655da1098e";
const char* apiKey = "533a269d-e866-43be-93af-e65a43905762";

// Parking spot configuration
const String spotNumbers[3] = {"A1", "A2", "A3"};
const int sensorPins[3] = {14, 12, 13}; // Using GPIO numbers for ESP8266 (D5=14, D6=12, D7=13)

// Status check interval (milliseconds)
const unsigned long checkInterval = 5000;
unsigned long lastCheckTime = 0;

void setup() {
  Serial.begin(115200);

  // Initialize I2C with custom SDA (D2=GPIO4) and SCL (D1=GPIO5)
  Wire.begin(4, 5);  // For ESP8266
  delay(100);

  lcd.begin(16, 2); // Initialize LCD with 16 columns and 2 rows for hd44780
  lcd.backlight();
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Park Me");
  lcd.setCursor(0, 1);
  int availableCount = 0;
  for (int i = 0; i < 3; i++) {
    if (digitalRead(sensorPins[i]) != LOW) {
      availableCount++;
    }
  }
  lcd.print("Available: ");
  lcd.print(availableCount);
  
  // Initialize sensor pins
  for (int i = 0; i < 3; i++) {
    pinMode(sensorPins[i], INPUT);
  }
  
  // Connect to WiFi
  WiFi.begin((char*)ssid, password);
  Serial.print("Connecting to WiFi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  if (millis() - lastCheckTime >= checkInterval) {
    lastCheckTime = millis();
    int availableCount = 0;
    // Check each sensor and send status
    for (int i = 0; i < 3; i++) {
      bool isOccupied = digitalRead(sensorPins[i]) == LOW; // TCRT5000 outputs LOW when object is detected
      if (!isOccupied) {
        availableCount++;
      }
      sendSpotStatus(spotNumbers[i], isOccupied);
    }
    lcd.setCursor(0, 1);
    lcd.print("Available: ");
    lcd.print(availableCount);
    lcd.print("   "); // Clear any leftover digits
  }
}

void sendSpotStatus(String spotNumber, bool isOccupied) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi not connected, skipping status update");
    return;
  }
  
  WiFiClient client;
  HTTPClient http;
  String url = "http://" + String(serverHost) + ":" + String(serverPort) + String(serverPath);
  
  // Create JSON payload
  DynamicJsonDocument doc(1024);
  doc["location_id"] = locationId;
  doc["spot_number"] = spotNumber;
  doc["is_occupied"] = isOccupied;
  
  String payload;
  serializeJson(doc, payload);
  
  http.begin(client, url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("X-API-Key", apiKey);
  
  int httpCode = http.POST(payload);
  if (httpCode > 0) {
    Serial.printf("[HTTP] POST... code: %d\n", httpCode);
    if (httpCode == HTTP_CODE_OK) {
      String response = http.getString();
      Serial.println(response);
    }
  } else {
    Serial.printf("[HTTP] POST... failed, error: %s\n", http.errorToString(httpCode).c_str());
  }
  
  http.end();
}
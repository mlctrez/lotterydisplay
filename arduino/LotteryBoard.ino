#include "LedControl.h"
#include <ESP8266WiFi.h>
#include <WiFiUdp.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>

// The MCU is a ESP8266 ESP-12F Mini (Wemos D1 mini)
// A TXS0108E level shifter is used for the clock/data/sel pins
// Display driver is a MAX7219CNG (https://www.adafruit.com/product/453)

const char* ssid = "secret";
const char* password = "secret";

WiFiUDP Udp;
unsigned int localUdpPort = 12345;  // local port to listen on
byte packetBuffer[255];             // buffer for incoming packets

LedControl lc = LedControl(D1, D2, D3, 1);
unsigned int connectCount = 0;

void setup() {
  lc.shutdown(0, false);
  lc.setIntensity(0, 15);

  lc.setDigit(0, 3, connectCount / 10, false);
  lc.setDigit(0, 4, connectCount % 10, false);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    connectCount++;
    lc.setDigit(0, 3, connectCount / 10, false);
    lc.setDigit(0, 4, connectCount % 10, false);
    delay(100);
  }

  Udp.begin(localUdpPort);
}

void loop() {
  delay(10);

  if (Udp.parsePacket()) {
    int len = Udp.read(packetBuffer, 255);

    if (len == 9) {
      for (int i = 0; i < 8; i++) {
        lc.setRow(0, i, packetBuffer[i]);
      }

      lc.setIntensity(0, packetBuffer[8]);
    }
  }
}

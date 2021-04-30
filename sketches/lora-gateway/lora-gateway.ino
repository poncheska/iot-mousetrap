#include <LoRa.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>

const char* ssid = "<WIFI_NAME>";
const char* password = "<WIFI_PASS>";


const String org_name = "<ORG_NAME>";
const String org_pass = "<ORG_PASS>";

String host = "https://smart-mousetrap.herokuapp.com";
String token = "";

#define ss 5
#define rst 14
#define dio0 2

void setup() {
  Serial.begin(115200);
  Serial.println("LoRa Receiver");

  LoRa.setPins(ss, rst, dio0);
  while (!LoRa.begin(433E6)) {
    Serial.println(".");
    delay(500);
  }
  LoRa.setSyncWord(0x13);
  Serial.println("LoRa Initializing OK!");

  WiFi.begin(ssid, password);
  Serial.println("Connecting");
  while(WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.print("Connected to WiFi network with IP Address: ");
  Serial.println(WiFi.localIP());

  HTTPClient http;
  String serverPath = host + "/org/sign-in";
  http.begin(serverPath.c_str());
  int httpResponseCode = -1;
  String response = "";
  while (httpResponseCode != 200){
    String httpRequestData = "{\"name\": \"" + org_name + "\", \"pass\": \"" + org_pass + "\"}";
    httpResponseCode = http.POST(httpRequestData);
    Serial.print("Sign-in response status code: ");
    Serial.println(httpResponseCode);
    response = http.getString();
    Serial.print("Sign-in response: ");
    Serial.println(response);
    if (httpResponseCode != 200){
      delay(10000);
    }
  }
  StaticJsonDocument<200> doc;
  DeserializationError error = deserializeJson(doc, response);
  if (error) {
    Serial.print(F("deserializeJson() failed: "));
    Serial.println(error.f_str());
    return;
  }
  const char* tokenCh = doc["token"];
  token = String(tokenCh);
  Serial.print("token: ");
  Serial.println(token);
  Serial.println();
  http.end();
}

void loop() {
  // try to parse packet
  int packetSize = LoRa.parsePacket();
  if (packetSize) {
    Serial.print("Received packet '");

    String LoRaData = "";
    while (LoRa.available()) {
      LoRaData = LoRa.readString();
      Serial.print(LoRaData);
    }
    Serial.print("' with RSSI ");
    Serial.println(LoRa.packetRssi());

    if (LoRaData != "")  {
      if(WiFi.status()== WL_CONNECTED){
        HTTPClient http;


        String serverPath = host + "/trigger/"+ LoRaData;

        http.begin(serverPath.c_str());
        http.addHeader("Authorization", "Bearer " + token);

        int httpResponseCode = http.GET();

        if (httpResponseCode>0) {
          Serial.print("HTTP Response code: ");
          Serial.println(httpResponseCode);
          String payload = http.getString();
          Serial.println(payload);
        }
        else {
          Serial.print("Error code: ");
          Serial.println(httpResponseCode);
        }
        http.end();
      }
      else {
        Serial.println("WiFi Disconnected");
      }
    }
    else {
      Serial.println("LoRa packet is empty");
    }
  }
}
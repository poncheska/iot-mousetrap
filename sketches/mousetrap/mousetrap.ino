#include <LoRa.h>

//define the pins used by the transceiver module
#define ss 5
#define rst 14
#define dio0 2

const String mousetrap_name = "<MT_NAME>";
const int buttonPin = 25;
int lastButtonState = HIGH;

void setup() {
  Serial.begin(115200);
  Serial.println("LoRa Sender");

  pinMode(buttonPin, INPUT_PULLUP);
  Serial.println("button pin input mode set");

  LoRa.setPins(ss, rst, dio0);

  while (!LoRa.begin(433E6)) {
    delay(500);
  }

  LoRa.setSyncWord(0x13);
  Serial.println("LoRa Initializing OK!");
}

void loop() {
  int reading = digitalRead(buttonPin);
  if (reading != lastButtonState)  {
    String sts = "";
    if (reading == HIGH){
      Serial.println("Mousetrap OFF");
      sts = "0";
    } else {
      Serial.println("Mousetrap ON");
      sts = "1";
    }
    LoRa.beginPacket();
    LoRa.print(mousetrap_name+"/");
    LoRa.print(sts);
    LoRa.endPacket();
    lastButtonState = reading;
  }
  delay(1000);
}
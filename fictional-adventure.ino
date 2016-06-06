#include <ESP8266WiFi.h>
#include <DHT.h>
#include <WiFiUdp.h>


#define DHTPIN 7
#define DHTTYPE DHT22

DHT dht(DHTPIN, DHTTYPE);

WiFiUDP udp;

void setup() {
  Serial.begin(115200);
  delay(100);

  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());

  pinMode(16, OUTPUT);
  digitalWrite(16, LOW);

  //dht.begin();
  digitalWrite(16, HIGH);

}

void loop() {
  digitalWrite(16, LOW);
  // Reading temperature or humidity takes about 250 milliseconds!
  // Sensor readings may also be up to 2 seconds 'old' (its a very slow sensor)
  //float humidity = dht.readHumidity();
  float humidity = 75.2;
  delay(250);
  // Read temperature as Fahrenheit (isFahrenheit = true)
  //float temperature = dht.readTemperature(true);
  float temperature = 75.8;
  delay(250);

  // Check if any reads failed and exit early (to try again)
  if (isnan(humidity) || isnan(temperature)) {
    Serial.println("Failed to read from DHT sensor!");
  } else {
    String string_temperature = String(temperature, 2);
    String string_humidity = String(humidity, 2);
    String line = String("crawl_space temperature=" + string_temperature + ",humidity=" + string_humidity);
    Serial.println(line);
    udp.beginPacket(influxdb_host, influxdb_udp_port);
    udp.print(line);
    udp.endPacket();
  }
  digitalWrite(16, HIGH);
  delay(2000);
}

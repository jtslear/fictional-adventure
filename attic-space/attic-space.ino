#include <ESP8266WiFi.h>
#include <DHT.h>
#include <WiFiUdp.h>

#define DHTPIN 5
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

  dht.begin();
  Serial.println("Started DHT");
  digitalWrite(16, HIGH);

}

void loop() {
  digitalWrite(16, LOW);
  float humidity = dht.readHumidity();
  // Read temperature as Fahrenheit (isFahrenheit = true)
  float temperature = dht.readTemperature(true);

  if (isnan(humidity) || isnan(temperature)) {
    Serial.println("Failed to read from DHT sensor!");
  } else {
    String string_temperature = String(temperature, 2);
    String string_humidity = String(humidity, 2);
    String line = String("attic_space temperature=" + string_temperature + ",humidity=" + string_humidity);
    Serial.println(line);
    udp.beginPacket(influxdb_host, influxdb_udp_port);
    udp.print(line);
    udp.endPacket();
  }
  digitalWrite(16, HIGH);
  // Reading temperature or humidity takes about 250 milliseconds!
  // Sensor readings may also be up to 2 seconds 'old' (its a very slow sensor)
  delay(2000);
}

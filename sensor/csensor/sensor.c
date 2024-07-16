#include "sensor.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>

double previous_temperature = 20.0; // Initial temperature

void generate_sensor_data(double *temperature, double *humidity) {
  // Simulate temperature between 20°C and 25°C with consecutive difference
  // limit
  double temp_diff =
      (rand() % 100) / 200.0 - 0.25; // Random difference between -0.25 and 0.25
  *temperature = previous_temperature + temp_diff;

  // Ensure temperature stays within range
  if (*temperature < 20.0)
    *temperature = 20.0;
  if (*temperature > 25.0)
    *temperature = 25.0;

  previous_temperature = *temperature;

  // Simulate humidity between 40% and 60%
  *humidity = (rand() % 200) / 10.0 + 40.0;
}

void send_data(double temperature, double humidity) {
  printf("[C] Temperature: %.2f°C, Humidity: %.2f%%\n", temperature, humidity);
}

void run_sensor(void) {
  srand(time(NULL));
  while (true) {
    double temperature, humidity;
    generate_sensor_data(&temperature, &humidity);
    send_data(temperature, humidity);
    sleep(2);
  }
}

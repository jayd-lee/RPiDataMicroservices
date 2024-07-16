#ifndef SENSOR_H
#define SENSOR_H

void generate_sensor_data(double *temperature, double *humidity);
void send_data(double temperature, double humidity);
void run_sensor(void);

#endif

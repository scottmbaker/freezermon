# freezermon
By Scott Baker, https://medium.com/@smbaker

Monitor the temperature of your freezer with a raspberry pi and Prometheus/Grafana

## why?

I've had my fridge and freezer go bad. It's not a pleasant experience. I wanted a simple
temperature monitor that will report the temperatures to prometheus and grafana, where
I can configure prometheus-alert-manager to generate alerts.

# Hardware

You'll need:

* A raspberry pi. I suggest going with a `raspberry pi zero 2 w`.

* A DS18B20 temperature sensor. Search for "ds18b20 waterproof" on amazon to get some
  plausible options.

* A 4.7 K resistor. Approximate value is fine. It's just a pullup.

* A case for the raspberry pi, if you want to make it nice.

# Wiring

* Connect the GND of the DS18B20 sensor to a GND on the pi.

* Connect the V+ of the DS18B20 sensor to +3.3V on the pi.

* Connect the data line of the DS18B20 to GPIO4 on the pi.

* Connect the 4.7K resistor from GPIO4 to +3.3V.
  (yes, we've already attached wires to GPIO4 and +3.3V)

# Building

* You can `make build` and `make install` directly on the pi.

* ... or you can `make release` somewhere else, and copy the built binary
  to your pi.

See INSTALL.md for further instructions.

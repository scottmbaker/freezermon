# Installation instructions

First, you will need to edit your /boot/config.txt and enable the 1wire driver
To do this, open /boot/config.txt in your favorite editor, and add the following
line at the end:

```
dtoverlay=w1-gpio
```

After this, it is wise to reboot and check to see that the 1wire driver is
active:

```
ls /sys/bus/w1/devices
```

In my case the output looks like this:

```
   28-0301979408ef  w1_bus_master1
```

The 28- device is the DS18B20 temperature sensor. If you don't have a sensor
attached to GPIO4, then obviously you won't see anything.

To autostart, run `sudo crontab -e` and add the following line:

```
@reboot /usr/local/freezermon/freezermon
```

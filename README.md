# Labrador - MQTT broker for router lab environment

Labrador is a simple MQTT broker for a router lab environment. It's at the core of the lab environment. 

![Labrador](labrador.webp)

## Features

It has a dashboard where it'll display the state of the devices in the lab. It gets this state from
snooping on the MQTT messages that the devices publish.

## Usage

Routerunner, when setting physical routers, are given two options:
 - What power device to use - POWER_DEVICE
 - What storage device to use - STORAGE_DEVICE

Both of these are device IDs. In the labrador MQTT topic hierarchy, these correspond with the following topics:
 - `lab/power/$POWER_DEVICE`
 - `lab/storage/$STORAGE_DEVICE`

### Built in power devices

Labrador comes with built-in power devices. These are based on Philips Hue smart plugs and work through the Hue bridge.
The package gohue, which is in turn generated out of the openhue api project yaml, is used to control the smart plugs.

Labrador tries to use the `HUE_BRIDGE` environment variable to connect to the bridge. If it's not set, it'll try to 
discover the bridge on the local network through mDNS.



## MQTT Topics

## General status topisc

lab/status
- anyone is free to post here to update the status of the lab
```json
{
  "status": "message"
}
```

### Power devices

lab/power/$device_name
- status - on state change a message is published to this topic
```json
{
    "power": true
}
```
- control - messages are published here to control the power state of the device. 
```json
{
  "power": true,
  "error": "error message"
}
```
Note that *requests* might be posted to this topic. The request itself will then specify where the
response should be posted. I recommend `lab/power/$device_name/response` for this.

### Magic storage gadgets - 
lab/storage/$device_name
- status - on state change a message is published to this topic
```json
{
  "active": true,
  "images": [
    {
      "source": "image_reference",
      "lun": 0,
      "size": 123456
    }
  ]
}

```
- control - messages are published here to control the state of the storage gadget
```json
{
  "active": true,
  "images": [
    {
      "reference": "image_reference"
    }
  ]
}
```
Note that *requests* might be posted to this topic. The request itself will then specify where the
response should be posted. I recommend `lab/storage/$device_name/response` for this.



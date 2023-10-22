![Logo](http://svg.wiersma.co.za/glasslabs/module?title=SOLAR&tag=a%20solar%20inverter%20module)

Solar inverter module for [looking glass](http://github.com/glasslabs/looking-glass)

## Usage

```yaml
modules:
 - name: simple-solar
    url:  https://github.com/glasslabs/solar/releases/download/v1.0.0/solar.wasm
    position: top:right
    config:
      url: http://my-hass-instance:8123
      token: <your-hass-token>
      sensorIds:
        load: sensor.deye_sunsynk_sol_ark_load_power
        pv: sensor.deye_sunsynk_sol_ark_pv_power
        battery: sensor.deye_sunsynk_sol_ark_battery_power
        batterySoC: sensor.deye_sunsynk_sol_ark_battery_state_of_charge
        grid: sensor.deye_sunsynk_sol_ark_generator_power
        gridFrequency: sensor.deye_sunsynk_sol_ark_grid_frequency
      battery:
        warning: 50
        low: 30
      maxWatts: 6000
```

## Configuration

### Load Sensor ID (sensorIds.load)

The Home Assistant load power sensor ID.

### PV Sensor ID (sensorIds.pv)

The Home Assistant PV power sensor ID.

### Battery Sensor ID (sensorIds.battery)

The Home Assistant Battery power sensor ID.

### Battery SoC Sensor ID (sensorIds.batterySoC)

The Home Assistant Battery State of Charge sensor ID.

### Grid Sensor ID (sensorIds.grid)

The Home Assistant Grid power sensor ID.

### Grid Frequency Sensor ID (sensorIds.gridFrequency)

The Home Assistant grid frequency sensor ID. This is used to determine if there is a grid connection.

### Battery Warning Percentage (battery.warning)

The Battery percentage for the battery bar to display in warning style.

### Battery Low Percentage (battery.low)

The Battery percentage for the battery bar to display in low style.

### Maximum Watts (maxWatts)

The maximum Watts used to scale the other sensors.

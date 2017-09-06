# Install Notes

## Networking
* Ensure that the Pi has connectivity. You can check this by pinging:
        ping HOSTNAME.byu.edu
  Or by checking on the elastic search for heartbeat events at `kibana-dev.byu.edu`
* If the pi is not responding check the VLANs and the IPaddress of the pi. you should be able to see the assigned IP address on the pi via the z pattern, or during startup. 
* If the pi comes up and is responsive but the icon doesn't stop spinning OR the wheel is not populated with a title/devices/anything the pi cannot connect with the devices in the room. --> Check VLAN/Physical connections

## Projector
* The Sony projectors must have PJLINK enabled with the password set. 
* The Sony projectors must have ADCP enabled, currently without authentication. 
* You can enable both in the Web UI. 

## TV
* The TV must have IP control enabled, authentication set to 'Normal and Pre-Shared Key'. The Pre-Shared key must be set. 
    * This is done in Nework Settings -> Home Network -> IP Control
* It is recommended that you disable automatic updates, as it will turn off ip control every time it updates. 

# Other Pi issues
* If the 'cannot connect' screen is up for long periods of time (> 2 minutes) the docker containers are weird. Reboot the pi. 
> At this point you can continue troubleshooting below, but it will require a little ssh knowledge. You can also call Engineering. 

* If it still does not come up - check the docker status with 
        docker ps 
  after ssh-ing into the pi. You should see something like: 
  
          CONTAINER ID        IMAGE                                                          COMMAND                  CREATED             STATUS              PORTS               NAMES
          87c0f9023d4f        byuoitav/rpi-pjlink-microservice:development                   "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_pjlink-microservice_1
          daf8eaa33c14        byuoitav/rpi-adcp-control-microservice:development             "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_adcp-control-microservice_1
          28dfe18c112d        byuoitav/rpi-av-api-rpc:development                            "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_av-api-rpc_1
          7dca36ce6f4c        byuoitav/rpi-touchpanel-ui-microservice:development            "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_touchpanel-ui-microservice_1
          24ba66b5c3cd        byuoitav/rpi-event-router-microservice:development             "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_event-router-microservice_1
          38ff17e2a988        byuoitav/rpi-pulse-eight-neo-microservice:development          "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_pulse-eight-neo-microservice_1
          a01e6ac0f3bc        byuoitav/rpi-event-translator-microservice:development         "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_event-translator-microservice_1
          3ee749d104a0        byuoitav/rpi-device-monitoring-microservice:development        "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_device-monitoring-microservice_1
          25182feacdd5        byuoitav/rpi-sony-control-microservice:development             "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_sony-control-microservice_1
          f5bf488e5137        byuoitav/rpi-configuration-database-microservice:development   "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_configuration-database-microservice_1
          d2fd9da5459e        byuoitav/rpi-av-api:development                                "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_av-api_1
          868da10e4ec6        byuoitav/rpi-london-audio-microservice:development             "/usr/bin/entry.sh..."   2 hours ago         Up 2 hours                              tmp_london-audio-microservice_1


* If nothing looks weird in the docker command (like a restarting container or something), contact engineering. 
* If there is a container that is having issues (denoted by the `status` column) Try a redeployment
* After SSH:
        sudo ./update_docker_containers.sh
* Wait 2 or so minutes and then check the `docker ps` command again - the `Created` column should show the docker containers got rebuilt. 
* At this point, reboot the pi - if there are still issues, call Engineering. 

### Bad UI Config
![Screen with stopped messages] (https://github.com/byuoitav/team/raw/master/images/bad-ui-config.jpg)
If you see this screen where the errors DO NOT disapper, it is usually caused by a bad UI config. Traditionally that there is a device in the UI config not found in the room configuration returned from the Database. 

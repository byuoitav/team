# Install Salt Master
> Note that if you see an error like `Cannot open: https://repo.saltstack.com/yum/redhat/salt-repo-2016.11-1.el7.noarch.rpm. Skipping.` when adding the repo you'll need to talk to the platform team to get the repo added. 
    1. Follow instructions found [here](https://repo.saltstack.com/#rhel)
    1. Install `salt1.master` and `salt1.minion`
    1. Run `systemctl start salt1.master`
    1. Run `systemctl start salt1.minion`
    1. Based on instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html)
    1. Run `salt1.key -F master`
    1. Copy the master.pub value and set it as the value of `master_finger` in `/etc/salt/minion`
    1. In `/etc/salt/minion` set `master` to the hostname of the master server.
    1. Be sure to restart the minion using `systemctl restart salt1.minion`
    1. Follow the instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html#key1.management) to accept the minion key.
    1. Validate by following instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html#sending-commands)


# Set up Eauth
You have to set up Eaut (External-Authentication) so that we can run salt commands as non-root users.
Based on Instructions found [here](https://docs.saltstack.com/en/latest/topics/eauth/index.html)
    1. Set up the `aveng` user on the server.
        - Instructions to add a user to RHEL7 found [here](https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/7/html/System_Administrators_Guide/s1-users-tools.html)
    1. Set up the PAM external_auth section in the `/etc/salt/master` file by adding the following to the end:
```
external_auth:
  pam:
      aveng:
            - .*
      skill_av_eng%:
            - .*
      av_eng%:
            - .*
```
    1. Restart salt-master service `systemctl restart salt-master`

Now that Eauth is set up you can run salt commands with `salt -a pam <command>` and Salt will prompt you for a password.
To create a token (for sequential commands) add the `-T` flag. `salt -T -a pam <command>`

#Setting up a heartbeat Beacon 
See documentation [here](https://docs.saltstack.com/en/latest/topics/beacons/)
A heartbeat beacon will send an event to the salt master every 10 seconds
    1. Add the following code into the `etc/salt/minion` file
```
beacons:
    status:
        - interval: 10
        - time:
            - all
        - cpustats:
            - all
        - meminfo:
            - all
```

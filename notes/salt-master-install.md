# Install Salt Master
> Note that if you see an error like `Cannot open: https://repo.saltstack.com/yum/redhat/salt-repo-2016.11-1.el7.noarch.rpm. Skipping.` when adding the repo you'll need to talk to the platform team to get the repo added. 
1. Talk to the newtorking team and let them know that you'll need your server to be reachable from the AV VRF.
1. Add the following lines into the `/etc/sysconfig/iptables` file to allow salt connections through the firewall:
    #Salt Stack Allow
    -A INPUT -m state --state new -m tcp -p tcp --dport 4505 -j ACCEPT
    -A INPUT -m state --state new -m tcp -p tcp --dport 4506 -j ACCEPT
1. Reload the iptables with `systemctl reload iptables`
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
3. Restart salt-master service `systemctl restart salt-master`

Now that Eauth is set up you can run salt commands with `salt -a pam <command>` and Salt will prompt you for a password.
To create a token (for sequential commands) add the `-T` flag. `salt -T -a pam <command>`

# Setting up a heartbeat Beacon 
See documentation [here](https://docs.saltstack.com/en/latest/topics/beacons/)
A heartbeat beacon will send an event to the salt master every 10 seconds
1. Add the following code into the `etc/salt/minion` file
```
beacons:
    status:
        - interval: 10
        - time:
            - all
        - loadavg:
            - all
```
> You may see a depreciation error like `[WARNING ] /usr/lib/python2.7/site-packages/salt/beacons/__init__.py:56: DeprecationWarning: Beacon configuration should be a list instead of a dictionary.` This is a known bug with the version of SALT we're using. See [here](https://github.com/saltstack/salt/issues/38121) - it appears as though it will be fixed in the next release.

# Installing Salt on a Pi
Follow instructions found [here](https://repo.saltstack.com/#raspbian). Note that we're using the 2016.11 release currently. And you will need to select the 'Pin to Major Version option' The version specific URLs are:

- `wget -O - https://repo.saltstack.com/apt/debian/8/armhf/2016.11/SALTSTACK-GPG-KEY.pub | sudo apt-key add -`
- `deb http://repo.saltstack.com/apt/debian/8/armhf/2016.11 jessie main`

Only install `salt-minion`

### Configure the pi
1. Based on instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html)
1. On the salt master run `salt1.key -F master` and copy the master.pub value.
1. Set the copied value as the value of `master_finger` in `/etc/salt/minion` on the minion
1. In `/etc/salt/minion` on the minion set `master` to the hostname of the master server.
1. Be sure to restart the minion using `systemctl restart salt-minion`
1. Follow the instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html#key1.management) to accept the minion key on the master.
1. Validate by following instructions found [here](https://docs.saltstack.com/en/latest/ref/configuration/index.html#sending-commands)




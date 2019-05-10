#!/usr/bin/env bash

echo -n "Proceed to reboot the Switchers? [y/n]: "

# This should reboot the ITB
echo "Rebooting ITB-1106B-SW1."
curl http://localhost:8011/ITB-1106B-SW1/reboot

echo "Reboot started, waiting 10 seconds to reboot next system..."
sleep 10 # Wait 10 seconds for the first switch to finish rebooting

echo "Rebooting EB-321-SW1."
curl http://localhost:8011/EB-321-SW1/reboot

echo "Reboot started, waiting 10 seconds to reboot next system"
sleep 10 # Wait another 10 seconds

echo "Rebooting EB-325-SW1."

curl http://localhost:8011/EB-325-SW1/reboot
echo "Reboot started, waiting 10 seconds to finish process"

sleep 10

echo "All done, See you again same time tomorrow!"

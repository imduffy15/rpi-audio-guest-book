#!/bin/bash

dmesg --follow | while read -r line; do
  if [[ "$line" == *"xHCI host controller not responding"* ]]; then
    echo -n "0000:01:00.0" > /sys/bus/pci/drivers/xhci_hcd/unbind
    echo -n "0000:01:00.0" > /sys/bus/pci/drivers/xhci_hcd/bind
  fi
done

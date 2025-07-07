#!/bin/bash

# File: active-service-manager.sh
# Description: Simple interactive service manager for running (active) systemd services

while true; do
  clear
  echo "==============================================="
  echo "🔧 Active Systemd Service Manager"
  echo "==============================================="
  echo
  echo "📋 List of currently running (active) services:"
  echo

  # Get list of active services
  mapfile -t SERVICES < <(systemctl list-units --type=service --state=running | awk 'NR>1 && NF {print $1}' | sort)

  if [[ ${#SERVICES[@]} -eq 0 ]]; then
    echo "⚠️  No running services found."
    read -rp "Press Enter to exit..."
    exit 0
  fi

  for i in "${!SERVICES[@]}"; do
    printf "%2d. %s\n" "$((i + 1))" "${SERVICES[$i]}"
  done

  echo
  read -rp "👉 Select a service by number (or 0 to exit): " choice

  if [[ "$choice" == "0" ]]; then
    echo "👋 Exiting..."
    exit 0
  fi

  if ! [[ "$choice" =~ ^[0-9]+$ ]] || ((choice < 1 || choice > ${#SERVICES[@]})); then
    echo "❌ Invalid selection. Press Enter to try again..."
    read
    continue
  fi

  selected_service="${SERVICES[$((choice - 1))]}"
  echo
  echo "🔍 Selected service: $selected_service"
  echo "-----------------------------------------------"
  systemctl status "$selected_service" --no-pager | grep -E "Loaded:|Active:"
  echo "-----------------------------------------------"
  echo "Available actions:"
  echo "  1. Stop     - Stop the service"
  echo "  2. Disable  - Prevent it from starting on boot"
  echo "  3. Back     - Return to the main menu"
  echo

  read -rp "🔧 Choose an action [1-3]: " action

  case "$action" in
  1)
    sudo systemctl stop "$selected_service" &&
      echo "✅ Service stopped: $selected_service"
    ;;
  2)
    sudo systemctl disable "$selected_service" &&
      echo "✅ Service disabled: $selected_service"
    ;;
  3)
    continue
    ;;
  *)
    echo "❌ Invalid action."
    ;;
  esac

  echo
  read -rp "🔁 Press Enter to return to the main menu..."
done

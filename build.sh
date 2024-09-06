#!/bin/bash

# Build the Go application
go build -o CityHallClock

# Create the app bundle structure
mkdir -p CityHallClock.app/Contents/MacOS
mkdir -p CityHallClock.app/Contents/Resources

# Move the executable to the app bundle
mv CityHallClock CityHallClock.app/Contents/MacOS/

# Copy the Info.plist
cp Info.plist CityHallClock.app/Contents/

# Copy the chime sound to Resources
cp chime.mp3 CityHallClock.app/Contents/Resources/

# Copy the menu bar icon
cp icon.png CityHallClock.app/Contents/Resources/

# Copy the application icon
cp AppIcon.icns CityHallClock.app/Contents/Resources/

echo "Application bundle created: CityHallClock.app"
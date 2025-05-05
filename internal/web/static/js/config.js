// config.js - Handles the configuration page functionality

// Global variables
let map;
let notificationCircle;
let scanCircle;
let locationMarker;
let configCoordinates = { latitude: null, longitude: null }; // Store coordinates for map link

document.addEventListener('DOMContentLoaded', () => {
    // Initialize theme toggle 
    initThemeToggle();

    // Fetch version information
    fetchVersionInfo();

    // Fetch configuration data
    fetchConfigData();
    
    // Update map link with coordinates
    updateMapLink();
    
    // Initialize user dropdown menu
    initUserDropdown();
});

// Initialize the theme toggle functionality
function initThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    const storedTheme = localStorage.getItem('theme') || 'light';
    
    // Apply the stored theme
    document.documentElement.setAttribute('data-theme', storedTheme);
    
    // Add event listener for theme toggle
    themeToggle.addEventListener('click', () => {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
    });
}

// Fetch version information from the API
async function fetchVersionInfo() {
    try {
        const response = await fetch('/api/version');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const versionData = await response.json();
        const appVersion = versionData.version || 'dev';
        
        // Update the version display in the footer
        const appVersionElement = document.getElementById('appVersion');
        if (appVersionElement) {
            appVersionElement.textContent = appVersion;
        }
    } catch (error) {
        console.error('Error fetching version information:', error);
        // Fall back to default 'dev' version if there's an error
        document.getElementById('appVersion').textContent = 'dev';
    }
}

// Fetch configuration data from the API
async function fetchConfigData() {
    try {
        const response = await fetch('/api/config');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const configData = await response.json();
        displayConfigData(configData);
    } catch (error) {
        console.error('Error fetching configuration data:', error);
        document.getElementById('loading-section').innerHTML = `
            <div class="error-message">
                Failed to load configuration data. Please refresh the page to try again.
            </div>
        `;
    }
}

// Display the configuration data on the page
function displayConfigData(config) {
    // Hide loading section
    document.getElementById('loading-section').style.display = 'none';
    
    // Show all sections
    document.getElementById('location-section').style.display = 'block';
    document.getElementById('tracking-section').style.display = 'block';
    document.getElementById('notification-filters-section').style.display = 'block';
    document.getElementById('notification-section').style.display = 'block';
    
    // Location data
    const latitude = config.Location?.Lat || 0;
    const longitude = config.Location?.Lon || 0;
    document.getElementById('latitude-value').textContent = latitude;
    document.getElementById('longitude-value').textContent = longitude;
    
    // Tracking settings
    const notificationRange = config.MaxRangeKilometers || 30;
    const scanRange = config.MaxScanRangeKilometers || config.MaxRangeKilometers || 30;
    
    document.getElementById('range-value').textContent = notificationRange;
    document.getElementById('scan-range-value').textContent = scanRange;
    document.getElementById('altitude-value').textContent = config.MaxAltitudeFeet === 0 ? 'No limit' : config.MaxAltitudeFeet || '-';
    document.getElementById('interval-value').textContent = config.FetchInterval || '-';
    
    // Initialize map if we have location data
    if (config.Location?.Lat && config.Location?.Lon) {
        initMap(latitude, longitude, notificationRange, scanRange);
    }
    
    // Aircraft types
    const aircraftTypes = config.AircraftTypes || [];
    const typesElement = document.getElementById('aircraft-types-value');
    
    if (aircraftTypes.length === 0) {
        typesElement.textContent = 'No specific types configured';
    } else if (aircraftTypes.includes('ALL')) {
        typesElement.textContent = 'All aircraft types';
    } else if (aircraftTypes.includes('MILITARY')) {
        const otherTypes = aircraftTypes.filter(type => type !== 'MILITARY');
        if (otherTypes.length === 0) {
            typesElement.textContent = 'Military aircraft only';
        } else {
            typesElement.textContent = `Military aircraft and: ${otherTypes.join(', ')}`;
        }
    } else {
        typesElement.textContent = aircraftTypes.join(', ');
    }
    
    // Notification services
    configureNotificationCard('discord', config.DiscordWebHookURL);
    configureNotificationCard('slack', config.SlackWebHookURL);
    configureNotificationCard('gotify', config.GotifyURL && config.GotifyToken);
    configureNotificationCard('ntfy', config.NtfyTopic);

    // Additional notification configs
    if (config.DiscordWebHookURL) {
        document.getElementById('discord-color-value').textContent = 
            config.DiscordColorAltitude === "true" ? "Yes" : "No";
    }
    
    if (config.NtfyTopic) {
        document.getElementById('ntfy-topic-value').textContent = config.NtfyTopic;
    }
}

// Configure notification card based on service status
function configureNotificationCard(service, isConfigured) {
    const card = document.getElementById(`${service}-card`);
    const statusIndicator = card.querySelector('.status-indicator');
    const statusText = card.querySelector('.status-text');
    const configSection = document.getElementById(`${service}-config`);
    
    if (isConfigured) {
        statusIndicator.classList.add('active');
        statusText.textContent = 'Active';
        
        if (configSection) {
            configSection.style.display = 'block';
        }
    } else {
        statusIndicator.classList.remove('active');
        statusText.textContent = 'Not configured';
        
        if (configSection) {
            configSection.style.display = 'none';
        }
    }
}

// Initialize Leaflet Map
function initMap(latitude, longitude, notificationRange, scanRange) {
    try {
        console.log("Initializing map with coordinates:", latitude, longitude);
        
        // Parse coordinates to ensure they're valid numbers
        const lat = parseFloat(latitude);
        const lng = parseFloat(longitude);
        
        if (isNaN(lat) || isNaN(lng)) {
            throw new Error("Invalid coordinates");
        }
        
        // Create the map
        map = L.map('location-map').setView([lat, lng], 10);
        
        // Add the base map layer (OpenStreetMap)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
        
        // Add a marker for the center location
        locationMarker = L.marker([lat, lng]).addTo(map);
        locationMarker.bindPopup(`
            <b>Tracking Location</b><br>
            Latitude: ${lat.toFixed(6)}<br>
            Longitude: ${lng.toFixed(6)}
        `);
        
        // Convert kilometers to meters for the circles
        const notificationRadiusMeters = notificationRange * 1000;
        const scanRadiusMeters = scanRange * 1000;
        
        // Add notification range circle
        notificationCircle = L.circle([lat, lng], {
            color: '#E74C3C',
            fillColor: '#E74C3C',
            fillOpacity: 0.2,
            radius: notificationRadiusMeters
        }).addTo(map);
        
        // Add scan range circle (only if different from notification range)
        if (scanRange > notificationRange) {
            scanCircle = L.circle([lat, lng], {
                color: '#3498DB',
                fillColor: '#3498DB',
                fillOpacity: 0.1,
                radius: scanRadiusMeters
            }).addTo(map);
            
            // Set the view to fit the larger circle
            map.fitBounds(scanCircle.getBounds());
        } else {
            // Set the view to fit the notification circle
            map.fitBounds(notificationCircle.getBounds());
        }
        
        // Apply dark mode if needed
        const currentTheme = document.documentElement.getAttribute('data-theme');
        if (currentTheme === 'dark') {
            applyDarkThemeToMap();
        }
        
    } catch (error) {
        console.error("Error initializing map:", error);
        document.getElementById('location-map').innerHTML = `
            <div style="padding: 20px; text-align: center; color: #e74c3c;">
                Failed to load the map. Please refresh the page to try again.
            </div>
        `;
    }
}

// Apply dark theme to map
function applyDarkThemeToMap() {
    // If using a dark theme-compatible tile server, you could switch here
    // For now, we'll just add a dark overlay to simulate dark mode
    if (map) {
        L.tileLayer('https://cartodb-basemaps-{s}.global.ssl.fastly.net/dark_all/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        }).addTo(map);
    }
}

// Initialize the user dropdown menu functionality
function initUserDropdown() {
    const userDropdownButton = document.getElementById('userDropdownButton');
    if (userDropdownButton) {
        userDropdownButton.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            const dropdown = this.closest('.user-dropdown');
            dropdown.classList.toggle('active');
            
            // Close dropdown when clicking outside
            document.addEventListener('click', function closeDropdown(e) {
                if (!dropdown.contains(e.target)) {
                    dropdown.classList.remove('active');
                    document.removeEventListener('click', closeDropdown);
                }
            });
        });
    }
}

// Initialize map link with coordinates
async function initializeMapLink() {
    try {
        const response = await fetch('/api/config');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const config = await response.json();
        
        // Store coordinates for future use
        if (config.Location?.Lat && config.Location?.Lon) {
            configCoordinates.latitude = config.Location.Lat;
            configCoordinates.longitude = config.Location.Lon;
            
            // Update map link
            updateMapLink();
        }
    } catch (error) {
        console.error('Error fetching configuration for map link:', error);
    }
}

// Update the map link with the current coordinates
function updateMapLink() {
    const mapLink = document.getElementById('mapLink');
    if (!mapLink) return;
    
    // Use coordinates from the template variables
    if (typeof SITE_LATITUDE !== 'undefined' && typeof SITE_LONGITUDE !== 'undefined') {
        const url = `https://globe.airplanes.live/?lat=${SITE_LATITUDE}&lon=${SITE_LONGITUDE}&SiteLat=${SITE_LATITUDE}&SiteLon=${SITE_LONGITUDE}&zoom=11&enableLabels&extendedLabels=1&hideSidebar`;
        mapLink.href = url;
    } else {
        mapLink.href = 'https://globe.airplanes.live/';
    }
}

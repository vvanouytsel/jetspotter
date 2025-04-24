// Global variables
let allAircraft = [];
let aircraftTypes = new Set(['ALL']);
let currentFilters = {
    type: 'ALL',
    military: false
};
let currentSort = 'distance';
let currentSortOrder = 'asc';
let lastUpdateTime = null;
let countdownInterval = null;

// DOM elements
document.addEventListener('DOMContentLoaded', () => {
    // Set up event listeners
    document.getElementById('filter-type').addEventListener('change', (e) => {
        currentFilters.type = e.target.value;
        renderAircraftGrid();
    });

    document.getElementById('filter-military').addEventListener('change', (e) => {
        currentFilters.military = e.target.checked;
        renderAircraftGrid();
    });

    document.getElementById('sort-by').addEventListener('change', (e) => {
        currentSort = e.target.value;
        renderAircraftGrid();
    });
    
    document.getElementById('sort-order').addEventListener('change', (e) => {
        currentSortOrder = e.target.value;
        renderAircraftGrid();
    });

    // Start fetching data
    fetchData();
    setInterval(fetchData, REFRESH_PERIOD * 1000);
    
    // Initialize countdown timer
    startCountdownTimer();
});

// Fetch aircraft data from the API
async function fetchData() {
    try {
        const response = await fetch('/api/aircraft');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        allAircraft = await response.json();
        
        // Update last update timestamp
        document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();
        
        // Update aircraft type filter options
        updateAircraftTypes();
        
        // Render the grid with the new data
        renderAircraftGrid();
        
        // Update dashboard stats
        updateStats();
        
        // Reset countdown timer
        startCountdownTimer();
    } catch (error) {
        console.error('Error fetching aircraft data:', error);
    }
}

// Update available aircraft types for filtering
function updateAircraftTypes() {
    // Extract all unique aircraft types
    allAircraft.forEach(aircraft => {
        if (aircraft.Type) {
            aircraftTypes.add(aircraft.Type);
        }
    });
    
    // Update the filter dropdown
    const typeFilter = document.getElementById('filter-type');
    
    // Save current selection
    const currentSelection = typeFilter.value;
    
    // Clear existing options (except ALL)
    while (typeFilter.options.length > 1) {
        typeFilter.remove(1);
    }
    
    // Add all aircraft types as options
    Array.from(aircraftTypes).sort().forEach(type => {
        if (type === 'ALL') return; // Skip the ALL option as it's already there
        
        const option = document.createElement('option');
        option.value = type;
        option.textContent = type;
        typeFilter.appendChild(option);
    });
    
    // Restore selection if possible
    if (Array.from(typeFilter.options).some(opt => opt.value === currentSelection)) {
        typeFilter.value = currentSelection;
    }
}

// Get altitude color class based on altitude
function getAltitudeColorClass(altitude) {
    if (!altitude) return "altitude-unknown";
    if (altitude < 1000) return "altitude-lt-1000";
    if (altitude >= 1000 && altitude < 2000) return "altitude-1000-2000";
    if (altitude >= 2000 && altitude < 3000) return "altitude-2000-3000";
    if (altitude >= 3000 && altitude < 5000) return "altitude-3000-5000";
    if (altitude >= 5000 && altitude < 7000) return "altitude-5000-7000";
    if (altitude >= 7000 && altitude < 10000) return "altitude-7000-10000";
    if (altitude >= 10000 && altitude < 15000) return "altitude-10000-15000";
    if (altitude >= 15000 && altitude < 20000) return "altitude-15000-20000";
    if (altitude >= 20000 && altitude < 30000) return "altitude-20000-30000";
    if (altitude >= 30000) return "altitude-gt-30000";
    return "altitude-unknown";
}

// Get altitude color header class based on altitude
function getAltitudeHeaderClass(altitude) {
    if (!altitude) return "aircraft-header-unknown";
    if (altitude < 1000) return "aircraft-header-lt-1000";
    if (altitude >= 1000 && altitude < 2000) return "aircraft-header-1000-2000";
    if (altitude >= 2000 && altitude < 3000) return "aircraft-header-2000-3000";
    if (altitude >= 3000 && altitude < 5000) return "aircraft-header-3000-5000";
    if (altitude >= 5000 && altitude < 7000) return "aircraft-header-5000-7000";
    if (altitude >= 7000 && altitude < 10000) return "aircraft-header-7000-10000";
    if (altitude >= 10000 && altitude < 15000) return "aircraft-header-10000-15000";
    if (altitude >= 15000 && altitude < 20000) return "aircraft-header-15000-20000";
    if (altitude >= 20000 && altitude < 30000) return "aircraft-header-20000-30000";
    if (altitude >= 30000) return "aircraft-header-gt-30000";
    return "aircraft-header-unknown";
}

// Update dashboard statistics
function updateStats() {
    const totalAircraft = allAircraft.length;
    const militaryAircraft = allAircraft.filter(a => a.Military).length;
    
    // Find closest and highest aircraft
    let closestAircraft = "-";
    let highestAircraft = "-";
    
    if (totalAircraft > 0) {
        // Sort by distance
        const sortedByDistance = [...allAircraft].sort((a, b) => (a.Distance || 0) - (b.Distance || 0));
        const closest = sortedByDistance[0];
        closestAircraft = `${closest.Callsign || 'Unknown'} (${closest.Distance}km)`;
        
        // Sort by altitude
        const sortedByAltitude = [...allAircraft].sort((a, b) => (b.Altitude || 0) - (a.Altitude || 0));
        const highest = sortedByAltitude[0];
        highestAircraft = `${highest.Callsign || 'Unknown'} (${Math.round(highest.Altitude).toLocaleString()}ft)`;
    }
    
    // Update DOM
    document.getElementById('totalAircraft').textContent = totalAircraft;
    document.getElementById('militaryAircraft').textContent = militaryAircraft;
    document.getElementById('closestAircraft').textContent = closestAircraft;
    document.getElementById('highestAircraft').textContent = highestAircraft;
}

// Filter and sort aircraft based on current settings
function getFilteredAndSortedAircraft() {
    // Apply filters
    let filtered = allAircraft;
    
    if (currentFilters.type !== 'ALL') {
        filtered = filtered.filter(aircraft => aircraft.Type === currentFilters.type);
    }
    
    if (currentFilters.military) {
        filtered = filtered.filter(aircraft => aircraft.Military);
    }
    
    // Apply sorting
    filtered.sort((a, b) => {
        switch (currentSort) {
            case 'distance':
                return (a.Distance || 0) - (b.Distance || 0);
            case 'altitude':
                return (b.Altitude || 0) - (a.Altitude || 0);
            case 'speed':
                return (b.Speed || 0) - (a.Speed || 0);
            case 'type':
                return (a.Type || '').localeCompare(b.Type || '');
            default:
                return 0;
        }
    });
    
    if (currentSortOrder === 'desc') {
        filtered.reverse();
    }
    
    return filtered;
}

// Render the aircraft grid with filtered and sorted data
function renderAircraftGrid() {
    const gridElement = document.getElementById('aircraftGrid');
    const noAircraftMessage = document.getElementById('noAircraftMessage');
    const filteredAircraft = getFilteredAndSortedAircraft();
    
    // Clear the grid except for the no-aircraft message
    Array.from(gridElement.children).forEach(child => {
        if (!child.classList.contains('no-aircraft')) {
            gridElement.removeChild(child);
        }
    });
    
    // Show/hide the no-aircraft message
    if (filteredAircraft.length === 0) {
        noAircraftMessage.style.display = 'block';
    } else {
        noAircraftMessage.style.display = 'none';
        
        // Render each aircraft card
        filteredAircraft.forEach(aircraft => {
            const card = createAircraftCard(aircraft);
            gridElement.appendChild(card);
        });
    }
}

// Create an aircraft card element
function createAircraftCard(aircraft) {
    const template = document.getElementById('aircraft-card-template');
    const card = document.importNode(template.content, true).querySelector('.aircraft-card');
    
    if (aircraft.Military) {
        card.classList.add('is-military');
    }
    
    // Set the callsign
    card.querySelector('.aircraft-callsign').textContent = aircraft.Callsign || 'Unknown';
    
    // Apply altitude-based color to the header
    const aircraftHeader = card.querySelector('.aircraft-header');
    const altitude = aircraft.Altitude || 0;
    aircraftHeader.classList.add(getAltitudeHeaderClass(altitude));
    
    // Set the image - use ImageURL as fallback if thumbnail is not available
    const imgElement = card.querySelector('.aircraft-image img');
    if (aircraft.ImageThumbnailURL) {
        imgElement.src = aircraft.ImageThumbnailURL;
        imgElement.alt = `${aircraft.Type || 'Aircraft'} - ${aircraft.Registration || ''}`;
    } else if (aircraft.ImageURL) {
        imgElement.src = aircraft.ImageURL;
        imgElement.alt = `${aircraft.Type || 'Aircraft'} - ${aircraft.Registration || ''}`;
    }
    
    // Add notification icons
    const notificationsContainer = card.querySelector('.aircraft-notifications');
    addNotificationIcons(notificationsContainer, aircraft);
    
    // Set basic info
    card.querySelector('.aircraft-type').textContent = aircraft.Type || 'Unknown';
    card.querySelector('.aircraft-registration').textContent = aircraft.Registration || 'Unknown';
    card.querySelector('.aircraft-icao').textContent = aircraft.ICAO || 'Unknown';
    
    // Apply altitude color coding and formatting
    const altitudeElement = card.querySelector('.aircraft-altitude');
    altitudeElement.textContent = aircraft.Altitude ? Math.round(aircraft.Altitude).toLocaleString() : 'Unknown';
    altitudeElement.classList.add(getAltitudeColorClass(altitude));
    
    card.querySelector('.aircraft-speed').textContent = aircraft.Speed || 'Unknown';
    card.querySelector('.aircraft-distance').textContent = aircraft.Distance || 'Unknown';
    card.querySelector('.aircraft-heading').textContent = aircraft.Heading ? Math.round(aircraft.Heading) : 'Unknown';
    
    // Set links (using default blue color)
    const trackerLink = card.querySelector('.tracker-link');
    if (aircraft.TrackerURL) {
        trackerLink.href = aircraft.TrackerURL;
    } else {
        trackerLink.style.display = 'none';
    }
    
    const imageLink = card.querySelector('.image-link');
    if (aircraft.ImageURL) {
        imageLink.href = aircraft.ImageURL;
    } else {
        imageLink.style.display = 'none';
    }
    
    return card;
}

// Add notification icons based on aircraft notification status
function addNotificationIcons(container, aircraft) {
    // If no notifications, hide the container
    if (!aircraft.NotifiedDiscord && !aircraft.NotifiedSlack && 
        !aircraft.NotifiedGotify && !aircraft.NotifiedNtfy && !aircraft.NotifiedTerminal) {
        container.style.display = 'none';
        return;
    }
    
    container.style.display = 'flex';
    
    // Create notification banner
    const bannerElement = document.createElement('div');
    bannerElement.className = 'notification-banner';
    
    // Create text for the banner based on which notifications were sent
    let notificationServices = [];
    if (aircraft.NotifiedDiscord) notificationServices.push('Discord');
    if (aircraft.NotifiedSlack) notificationServices.push('Slack');
    if (aircraft.NotifiedGotify) notificationServices.push('Gotify');
    if (aircraft.NotifiedNtfy) notificationServices.push('Ntfy');
    if (aircraft.NotifiedTerminal) notificationServices.push('Terminal');
    
    // Format the notification text
    let bannerText = '';
    if (notificationServices.length === 1) {
        bannerText = `Notification sent to ${notificationServices[0]}`;
    } else if (notificationServices.length === 2) {
        bannerText = `Notifications sent to ${notificationServices[0]} and ${notificationServices[1]}`;
    } else if (notificationServices.length > 2) {
        const lastService = notificationServices.pop();
        bannerText = `Notifications sent to ${notificationServices.join(', ')} and ${lastService}`;
    }
    
    bannerElement.textContent = bannerText;
    container.appendChild(bannerElement);
    
    // Create container for icons
    const iconsContainer = document.createElement('div');
    iconsContainer.className = 'notification-icons';
    container.appendChild(iconsContainer);
    
    // Add Discord icon
    if (aircraft.NotifiedDiscord) {
        const discordIcon = document.createElement('div');
        discordIcon.className = 'notification-icon notification-discord';
        discordIcon.title = 'Sent to Discord';
        discordIcon.textContent = 'D';
        iconsContainer.appendChild(discordIcon);
    }
    
    // Add Slack icon
    if (aircraft.NotifiedSlack) {
        const slackIcon = document.createElement('div');
        slackIcon.className = 'notification-icon notification-slack';
        slackIcon.title = 'Sent to Slack';
        slackIcon.textContent = 'S';
        iconsContainer.appendChild(slackIcon);
    }
    
    // Add Gotify icon
    if (aircraft.NotifiedGotify) {
        const gotifyIcon = document.createElement('div');
        gotifyIcon.className = 'notification-icon notification-gotify';
        gotifyIcon.title = 'Sent to Gotify';
        gotifyIcon.textContent = 'G';
        iconsContainer.appendChild(gotifyIcon);
    }
    
    // Add Ntfy icon
    if (aircraft.NotifiedNtfy) {
        const ntfyIcon = document.createElement('div');
        ntfyIcon.className = 'notification-icon notification-ntfy';
        ntfyIcon.title = 'Sent to Ntfy';
        ntfyIcon.textContent = 'N';
        iconsContainer.appendChild(ntfyIcon);
    }
    
    // Add Terminal icon
    if (aircraft.NotifiedTerminal) {
        const terminalIcon = document.createElement('div');
        terminalIcon.className = 'notification-icon notification-terminal';
        terminalIcon.title = 'Sent to Terminal';
        terminalIcon.textContent = 'T';
        iconsContainer.appendChild(terminalIcon);
    }
}

// Start the countdown timer
function startCountdownTimer() {
    // Clear any existing interval
    if (countdownInterval) {
        clearInterval(countdownInterval);
    }
    
    // Set the last update time to now
    lastUpdateTime = new Date();
    
    // Update the countdown immediately
    updateCountdown();
    
    // Set up the countdown interval (update every second)
    countdownInterval = setInterval(updateCountdown, 1000);
}

// Update the countdown display
function updateCountdown() {
    const nextUpdateCountdown = document.getElementById('nextUpdateCountdown');
    
    if (!lastUpdateTime) {
        nextUpdateCountdown.textContent = '-';
        return;
    }
    
    // Calculate time passed since last update
    const now = new Date();
    const elapsedSeconds = Math.floor((now - lastUpdateTime) / 1000);
    const secondsLeft = REFRESH_PERIOD - elapsedSeconds;
    
    // Format the countdown
    if (secondsLeft <= 0) {
        nextUpdateCountdown.textContent = 'Refreshing...';
    } else {
        const minutes = Math.floor(secondsLeft / 60);
        const seconds = secondsLeft % 60;
        nextUpdateCountdown.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;
    }
}

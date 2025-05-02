// Global variables
let allAircraft = [];
let aircraftDescriptions = new Set();
let currentFilters = {
    description: '',
    statuses: {
        military: false,
        inbound: false,
        hideGround: false
    }
};
let currentSort = 'distance';
let currentSortOrder = 'asc';
let lastUpdateTime = null;
let countdownInterval = null;
let isLoading = true; // New variable to track loading state
let appVersion = 'dev'; // Variable to store app version
let configCoordinates = { latitude: null, longitude: null }; // Store coordinates for map link

// DOM elements
document.addEventListener('DOMContentLoaded', () => {
    // Set up event listeners
    document.getElementById('filter-description').addEventListener('change', (e) => {
        currentFilters.description = e.target.value;
        renderAircraftGrid();
    });

    document.getElementById('filter-military').addEventListener('change', (e) => {
        currentFilters.statuses.military = e.target.checked;
        renderAircraftGrid();
    });

    document.getElementById('filter-inbound').addEventListener('change', (e) => {
        currentFilters.statuses.inbound = e.target.checked;
        renderAircraftGrid();
    });
    
    document.getElementById('filter-hide-ground').addEventListener('change', (e) => {
        currentFilters.statuses.hideGround = e.target.checked;
        renderAircraftGrid();
    });

    // Update map link with coordinates
    updateMapLink();
    
    // Add click handler for the total aircraft stat box
    document.getElementById('totalAircraftStat').addEventListener('click', () => {
        // Reset all filters
        currentFilters.description = '';
        currentFilters.statuses.military = false;
        currentFilters.statuses.inbound = false;
        
        // Update checkbox states in the UI
        document.getElementById('filter-description').value = '';
        document.getElementById('filter-military').checked = false;
        document.getElementById('filter-inbound').checked = false;
        
        // Refresh the grid
        renderAircraftGrid();
    });
    
    // Add click handler for the military aircraft stat box
    document.getElementById('militaryAircraftStat').addEventListener('click', () => {
        // Toggle the military filter checkbox
        const militaryCheckbox = document.getElementById('filter-military');
        militaryCheckbox.checked = !militaryCheckbox.checked;
        
        // Update filter state and refresh the grid
        currentFilters.statuses.military = militaryCheckbox.checked;
        renderAircraftGrid();
    });
    
    // Add click handler for the inbound aircraft stat box
    document.getElementById('inboundAircraftStat').addEventListener('click', () => {
        // Toggle the inbound filter checkbox
        const inboundCheckbox = document.getElementById('filter-inbound');
        inboundCheckbox.checked = !inboundCheckbox.checked;
        
        // Update filter state and refresh the grid
        currentFilters.statuses.inbound = inboundCheckbox.checked;
        renderAircraftGrid();
    });

    document.getElementById('sort-by').addEventListener('change', (e) => {
        currentSort = e.target.value;
        renderAircraftGrid();
    });
    
    // Initialize and set up the sort direction toggle
    const sortDirectionToggle = document.getElementById('sort-direction-toggle');
    sortDirectionToggle.classList.add(currentSortOrder);
    
    sortDirectionToggle.addEventListener('click', () => {
        // Toggle between asc and desc
        currentSortOrder = currentSortOrder === 'asc' ? 'desc' : 'asc';
        
        // Update the button appearance
        sortDirectionToggle.classList.remove('asc', 'desc');
        sortDirectionToggle.classList.add(currentSortOrder);
        
        // Update the label text
        sortDirectionToggle.querySelector('.sort-label').textContent = 
            currentSortOrder === 'asc' ? 'Ascending' : 'Descending';
        
        // Refresh the grid with the new sort order
        renderAircraftGrid();
    });

    // Initialize dark mode
    initializeTheme();
    
    // Theme toggle button event listener
    document.getElementById('themeToggle').addEventListener('click', toggleTheme);

    // Add event listener to the user dropdown button
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

    // Fetch version information
    fetchVersionInfo();
    
    // Start fetching data
    fetchData();
    setInterval(fetchData, REFRESH_PERIOD * 1000);
    
    // Initialize countdown timer
    startCountdownTimer();
});

// Initialize theme from localStorage or system preference
function initializeTheme() {
    // Check if theme is stored in localStorage
    const storedTheme = localStorage.getItem('theme');
    
    if (storedTheme) {
        // Apply stored theme
        document.documentElement.setAttribute('data-theme', storedTheme);
    } else {
        // Check for system preference
        const prefersDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const initialTheme = prefersDarkMode ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', initialTheme);
        localStorage.setItem('theme', initialTheme);
    }
}

// Toggle between light and dark themes
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    
    // Set the theme attribute on the html element
    document.documentElement.setAttribute('data-theme', newTheme);
    
    // Save theme preference to localStorage
    localStorage.setItem('theme', newTheme);
}

// Fetch aircraft data from the API
async function fetchData() {
    try {
        // Only show loading spinner on initial load, not during refreshes
        const initialLoad = allAircraft.length === 0;
        if (initialLoad) {
            isLoading = true;
            renderAircraftGrid(); // Show loading indicator
        }
        
        // Check for test mode in URL parameters
        const urlParams = new URLSearchParams(window.location.search);
        const testMode = urlParams.get('test');
        
        if (testMode === 'true') {
            // Use demo data instead of fetching from API
            const demoData = generateDemoData();
            allAircraft = demoData;
            
            // Update UI with demo data
            lastUpdateTime = new Date();
            updateAircraftDescriptions();
            updateStats();
            isLoading = false;
            renderAircraftGrid();
            return;
        }
        
        const response = await fetch('/api/aircraft');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const newAircraft = await response.json();
        
        // Always update last update timestamp - this is key to fixing the countdown timer
        lastUpdateTime = new Date();
        
        // Only update the UI if data has changed
        if (JSON.stringify(allAircraft) !== JSON.stringify(newAircraft)) {
            allAircraft = newAircraft;
            
            // Update aircraft description dropdown
            updateAircraftDescriptions();
            
            // Update dashboard stats
            updateStats();
        }
        
        // Set loading to false now that we have data
        isLoading = false;
        
        // Render the grid with the new data
        renderAircraftGrid();
    } catch (error) {
        console.error('Error fetching aircraft data:', error);
        // Even on error, update timestamp and reset the timer
        lastUpdateTime = new Date();
        // Stop showing loading
        isLoading = false;
        renderAircraftGrid();
    }
}

// Generate demo aircraft data for testing
function generateDemoData() {
    return [
        // Military aircraft example
        {
            Callsign: "ARMY01",
            Country: "Unknown",
            Military: true,
            ICAO: "AE01FF",
            Altitude: 25000,
            Speed: 450,
            Distance: 42.3,
            Heading: 225,
            BearingFromLocation: 170,
            Inbound: false,
            Type: "C-17A Globemaster III",
            Description: "Military Transport Aircraft",
            Registration: "02-1111",
            TrackerURL: "javascript:void(0)",
            NotifiedDiscord: true,
            NotifiedSlack: true
        },
        // Military inbound aircraft example
        {
            Callsign: "NAVY42",
            Country: "Unknown",
            Military: true,
            ICAO: "AE1234",
            Altitude: 18000,
            Speed: 380,
            Distance: 25.7,
            Heading: 90,
            BearingFromLocation: 45,
            Inbound: true,
            Type: "F/A-18E Super Hornet",
            Description: "Military Fighter Aircraft",
            Registration: "165667",
            TrackerURL: "javascript:void(0)",
            NotifiedDiscord: true,
            NotifiedTerminal: true
        },
        // Civilian aircraft example
        {
            Callsign: "UAL123",
            Country: "United States",
            Military: false,
            ICAO: "A12345",
            Altitude: 35000,
            Speed: 520,
            Distance: 15.2,
            Heading: 45,
            BearingFromLocation: 280,
            Inbound: false,
            Type: "Boeing 777-300ER",
            Description: "Twin-Engine Passenger Jet",
            Registration: "N12345",
            TrackerURL: "javascript:void(0)",
            ImageURL: "https://picsum.photos/600/400",
            ImageThumbnailURL: "https://picsum.photos/300/200"
        },
        // Civilian aircraft invalid countryexample
        {
            Callsign: "ABC123",
            Country: "Does not exist",
            Military: false,
            ICAO: "A12345",
            Altitude: 35000,
            Speed: 520,
            Distance: 15.2,
            Heading: 45,
            BearingFromLocation: 280,
            Inbound: false,
            Type: "Boeing 777-300ER",
            Description: "Twin-Engine Passenger Jet",
            Registration: "ABCD",
            TrackerURL: "javascript:void(0)",
            ImageURL: "https://picsum.photos/600/400",
            ImageThumbnailURL: "https://picsum.photos/300/200"
        },
        // Civilian inbound aircraft example
        {
            Callsign: "DAL456",
            Country: "United States",
            Military: false,
            ICAO: "DAL456",
            Altitude: 10000,
            Speed: 320,
            Distance: 8.5,
            Heading: 180,
            BearingFromLocation: 90,
            Inbound: true,
            Type: "Airbus A320",
            Description: "Twin-Engine Passenger Jet",
            Registration: "N987AA",
            TrackerURL: "javascript:void(0)",
            ImageURL: "https://picsum.photos/600/400",
            ImageThumbnailURL: "https://picsum.photos/300/200",
            NotifiedSlack: true
        },
        // Military aircraft on ground example
        {
            Callsign: "AF101",
            Country: "United States",
            Military: true,
            ICAO: "AE2468",
            Altitude: 0,
            Speed: 0,
            Distance: 5.1,
            Heading: 90,
            BearingFromLocation: 120,
            Inbound: false,
            OnGround: true,
            Type: "F-16 Fighting Falcon",
            Description: "Military Fighter Aircraft",
            Registration: "86-0241",
            TrackerURL: "javascript:void(0)",
            ImageURL: "https://picsum.photos/600/400",
            ImageThumbnailURL: "https://picsum.photos/300/200",
            NotifiedGotify: true
        },
        // Civilian aircraft on ground example
        {
            Callsign: "UAL789",
            Country: "United States",
            Military: false,
            ICAO: "A54321",
            Altitude: 0,
            Speed: 0,
            Distance: 3.2,
            Heading: 270,
            BearingFromLocation: 200,
            Inbound: false,
            OnGround: true,
            Type: "Boeing 737-800",
            Description: "Twin-Engine Passenger Jet",
            Registration: "N12346",
            TrackerURL: "javascript:void(0)",
            ImageURL: "https://picsum.photos/600/400",
            ImageThumbnailURL: "https://picsum.photos/300/200",
            NotifiedNtfy: true
        }
    ];
}

// Update available aircraft descriptions for filtering
function updateAircraftDescriptions() {
    // Clear previous set
    aircraftDescriptions = new Set();
    
    // Extract all unique aircraft descriptions
    allAircraft.forEach(aircraft => {
        const description = aircraft.Description || aircraft.Type || '';
        if (description) {
            aircraftDescriptions.add(description);
        }
    });
    
    // Update the filter dropdown
    const descriptionFilter = document.getElementById('filter-description');
    
    // Save current selection
    const currentSelection = descriptionFilter.value;
    
    // Clear existing options (except the "All Aircraft" option)
    while (descriptionFilter.options.length > 1) {
        descriptionFilter.remove(1);
    }
    
    // Add all aircraft descriptions as options
    Array.from(aircraftDescriptions).sort().forEach(description => {
        const option = document.createElement('option');
        option.value = description;
        option.textContent = description;
        descriptionFilter.appendChild(option);
    });
    
    // Restore selection if possible
    if (Array.from(descriptionFilter.options).some(opt => opt.value === currentSelection)) {
        descriptionFilter.value = currentSelection;
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
    const inboundAircraft = allAircraft.filter(a => a.Inbound).length;
    
    // Update DOM
    document.getElementById('totalAircraft').textContent = totalAircraft;
    document.getElementById('militaryAircraft').textContent = militaryAircraft;
    document.getElementById('inboundAircraft').textContent = inboundAircraft;
    
    // Update the active state of the filter tiles based on current filters
    document.getElementById('inboundAircraftStat').classList.toggle('active', currentFilters.statuses.inbound);
    document.getElementById('militaryAircraftStat').classList.toggle('active', currentFilters.statuses.military);
}

// Filter and sort aircraft based on current settings
function getFilteredAndSortedAircraft() {
    // Apply filters
    let filtered = allAircraft;
    
    // Filter by description
    if (currentFilters.description && currentFilters.description !== '') {
        filtered = filtered.filter(aircraft => {
            const description = aircraft.Description || aircraft.Type || '';
            return description === currentFilters.description;
        });
    }
    
    if (currentFilters.statuses.military) {
        filtered = filtered.filter(aircraft => aircraft.Military);
    }
    
    if (currentFilters.statuses.inbound) {
        filtered = filtered.filter(aircraft => aircraft.Inbound);
    }
    
    if (currentFilters.statuses.hideGround) {
        filtered = filtered.filter(aircraft => !aircraft.OnGround);
    }
    
    // Apply sorting
    filtered.sort((a, b) => {
        switch (currentSort) {
            case 'distance':
                return (a.Distance || 0) - (b.Distance || 0);
            case 'altitude':
                // Fix sorting so ascending means lower altitudes first
                return (a.Altitude || 0) - (b.Altitude || 0);
            case 'speed':
                // Fix sorting so ascending means lower speeds first
                return (a.Speed || 0) - (b.Speed || 0);
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
    
    // Update the active state of the filter tiles based on current filters
    document.getElementById('inboundAircraftStat').classList.toggle('active', currentFilters.statuses.inbound);
    document.getElementById('militaryAircraftStat').classList.toggle('active', currentFilters.statuses.military);
    
    // Show/hide the no-aircraft message or loading spinner
    if (isLoading) {
        noAircraftMessage.style.display = 'flex';
        noAircraftMessage.innerHTML = `
            <div class="scanning-animation">
                <div class="radar-circle"></div>
                <div class="radar-sweep"></div>
            </div>
            <div class="scanning-text">
                <div class="scanning-title">Scanning the skies...</div>
                <div class="scanning-subtitle">Looking for aircraft in your vicinity</div>
            </div>
        `;
        noAircraftMessage.classList.add('loading', 'scanning');
    } else if (filteredAircraft.length === 0) {
        noAircraftMessage.style.display = 'flex';
        // Use the same message for consistency
        noAircraftMessage.innerHTML = `
            <div class="scanning-animation">
                <div class="radar-circle"></div>
                <div class="radar-sweep"></div>
            </div>
            <div class="scanning-text">
                <div class="scanning-title">Scanning the skies...</div>
                <div class="scanning-subtitle">Looking for aircraft in your vicinity</div>
            </div>
        `;
        noAircraftMessage.classList.add('scanning');
        noAircraftMessage.classList.remove('loading');
    } else {
        noAircraftMessage.style.display = 'none';
        noAircraftMessage.classList.remove('loading', 'scanning');
        
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
    
    // Set country info
    const countryName = aircraft.Country || 'Unknown';
    card.querySelector('.aircraft-country-info').textContent = countryName;
    
    // Add a flag emoji in the header only if country is known
    const flagElement = card.querySelector('.aircraft-country-flag');
    if (countryName !== 'Unknown') {
        flagElement.textContent = getCountryFlagEmoji(countryName);
        flagElement.style.display = 'inline-block';
    } else {
        // Don't show any flag for unknown countries
        flagElement.textContent = '';
        flagElement.style.display = 'none';
    }
    
    // Apply altitude-based color to the header
    const aircraftHeader = card.querySelector('.aircraft-header');
    const altitude = aircraft.Altitude || 0;
    aircraftHeader.classList.add(getAltitudeHeaderClass(altitude));
    
    // Show/hide military badge based on aircraft status
    const militaryBadge = card.querySelector('.aircraft-military-badge');
    militaryBadge.style.display = aircraft.Military ? 'block' : 'none';
    
    // Handle inbound status display
    const approachBadge = card.querySelector('.aircraft-approach-badge');
    if (aircraft.Inbound) {
        approachBadge.style.display = 'block';
        card.classList.add('is-inbound');
    } else {
        approachBadge.style.display = 'none';
    }
    
    // Handle on ground status display
    const groundBadge = card.querySelector('.aircraft-ground-badge');
    if (aircraft.OnGround) {
        groundBadge.style.display = 'block';
        card.classList.add('is-on-ground');
    } else {
        groundBadge.style.display = 'none';
    }
    
    // Get the image container
    const imageContainer = card.querySelector('.aircraft-image');
    // Clear the container first
    imageContainer.innerHTML = '';
    
    // Create link element that will wrap the image
    const imageLink = document.createElement('a');
    imageLink.href = aircraft.ImageURL || 'javascript:void(0)'; 
    imageLink.target = '_blank';
    imageLink.style.display = 'block';
    imageLink.style.width = '100%';
    imageLink.style.height = '100%';
    
    // If the URL isn't available, make it non-clickable
    if (!aircraft.ImageURL) {
        imageLink.style.cursor = 'default';
        imageLink.onclick = (e) => e.preventDefault();
    } else {
        imageLink.title = 'Click to view more images of this aircraft';
    }
    
    // Create the image element
    const imgElement = document.createElement('img');
    if (aircraft.ImageThumbnailURL) {
        imgElement.src = aircraft.ImageThumbnailURL;
        imgElement.alt = `${aircraft.Type || 'Aircraft'} - ${aircraft.Registration || ''}`;
        
        // Add photographer information as tooltip if available
        if (aircraft.Photographer) {
            imgElement.title = `Photo by: ${aircraft.Photographer}`;
        }
    } else if (aircraft.ImageURL) {
        imgElement.src = aircraft.ImageURL;
        imgElement.alt = `${aircraft.Type || 'Aircraft'} - ${aircraft.Registration || ''}`;
        
        // Add photographer information as tooltip if available
        if (aircraft.Photographer) {
            imgElement.title = `Photo by: ${aircraft.Photographer}`;
        }
    } else {
        imgElement.src = '/static/images/aircraft_not_found.png';
        imgElement.alt = 'No image available';
        imgElement.classList.add('fallback-image');
    }

    // Add the image to the link
    imageLink.appendChild(imgElement);
    // Add the link to the image container
    imageContainer.appendChild(imageLink);
    
    // Add notification icons
    const notificationsContainer = card.querySelector('.aircraft-notifications');
    addNotificationIcons(notificationsContainer, aircraft);
    
    // Set basic info
    card.querySelector('.aircraft-description').textContent = aircraft.Description || aircraft.Type || 'Unknown';
    card.querySelector('.aircraft-registration').textContent = aircraft.Registration || 'Unknown';
    card.querySelector('.aircraft-icao').textContent = aircraft.ICAO || 'Unknown';
    
    // Apply altitude color coding and formatting
    const altitudeElement = card.querySelector('.aircraft-altitude');
    if (aircraft.OnGround) {
        altitudeElement.textContent = 'On ground';
        altitudeElement.classList.add('altitude-on-ground');
    } else {
        altitudeElement.textContent = aircraft.Altitude ? Math.round(aircraft.Altitude).toLocaleString() : 'Unknown';
        altitudeElement.classList.add(getAltitudeColorClass(altitude));
    }
    
    // Set speed with special handling for ground aircraft
    const speedElement = card.querySelector('.aircraft-speed');
    if (aircraft.OnGround) {
        speedElement.textContent = 'N/A';
        speedElement.classList.add('value-na');
    } else {
        speedElement.textContent = aircraft.Speed || 'Unknown';
        speedElement.classList.remove('value-na');
    }
    
    // Set distance with special handling for ground aircraft
    const distanceElement = card.querySelector('.aircraft-distance');
    if (aircraft.OnGround) {
        distanceElement.textContent = 'N/A';
        distanceElement.classList.add('value-na');
    } else {
        distanceElement.textContent = aircraft.Distance || 'Unknown';
        distanceElement.classList.remove('value-na');
    }
    
    // Fix: Use aircraft.Heading instead of the undefined 'heading' variable
    const heading = aircraft.Heading;
    const headingElement = card.querySelector('.aircraft-heading');
    const headingIndicator = card.querySelector('.heading-indicator');
    
    if (aircraft.OnGround) {
        // Hide heading information for grounded aircraft
        headingElement.textContent = 'N/A';
        headingElement.classList.add('value-na');
        headingIndicator.style.display = 'none';
    } else {
        // Show heading for airborne aircraft
        headingElement.textContent = heading ? Math.round(heading) : 'Unknown';
        headingElement.classList.remove('value-na');
        headingIndicator.style.display = 'flex';
        
        // Set aircraft heading direction indicator
        if (heading !== undefined && heading !== null) {
            // Rotate the SVG to match the aircraft's heading
            headingIndicator.style.transform = `rotate(${heading}deg)`;
            headingIndicator.title = `Aircraft heading direction: ${Math.round(heading)}°`;
        }
    }
    
    // Set bearing and rotate the compass arrow
    const bearing = aircraft.BearingFromLocation;
    if (bearing !== undefined && bearing !== null) {
        const directionIndicator = card.querySelector('.direction-indicator');
        const directionArrow = directionIndicator.querySelector('svg');
        
        // Rotate the SVG for the compass arrow
        directionArrow.style.transform = `rotate(${bearing}deg)`;
        
        // Update the tooltip to show the bearing
        const bearingRounded = Math.round(bearing);
        directionIndicator.title = `Direction to aircraft from your location: ${bearingRounded}°`;
    }
    
    // Set up links section with appropriate disabled states for unavailable resources
    const trackerLink = card.querySelector('.tracker-link');
    const bottomImageLink = card.querySelector('.image-link');
    
    // Track Aircraft link
    if (aircraft.TrackerURL) {
        trackerLink.href = aircraft.TrackerURL;
        trackerLink.classList.remove('disabled');
    } else {
        trackerLink.href = 'javascript:void(0)';
        trackerLink.classList.add('disabled');
        trackerLink.title = 'Tracking not available for this aircraft';
    }
    
    // More Images link (bottom toolbar)
    if (aircraft.ImageURL) {
        bottomImageLink.href = aircraft.ImageURL;
        bottomImageLink.classList.remove('disabled');
        bottomImageLink.title = 'More images of this aircraft';
        bottomImageLink.target = '_blank';
    } else {
        bottomImageLink.href = 'javascript:void(0)';
        bottomImageLink.classList.add('disabled');
        bottomImageLink.title = 'No additional images available';
        bottomImageLink.target = '_self';
    }

    // Remove any existing click handlers for the flight details button
    bottomImageLink.removeEventListener('click', showFlightDetails);
    
    // Check if this aircraft has flight details (origin/destination)
    const hasFlightDetails = aircraft.Origin && (
        (typeof aircraft.Origin === 'string' && aircraft.Origin !== 'null' && aircraft.Origin.length > 2) || 
        (typeof aircraft.Origin === 'object' && aircraft.Origin !== null)
    );
    
    // Convert from flight image button to flight details button if flight information exists
    if (hasFlightDetails) {
        bottomImageLink.title = 'Show flight details';
        bottomImageLink.classList.remove('disabled');
        bottomImageLink.href = 'javascript:void(0)';
        bottomImageLink.target = '_self';
        
        // Add click event listener for flight details
        bottomImageLink.addEventListener('click', (e) => {
            e.preventDefault();
            showFlightDetails(aircraft);
        });
    } else {
        // If the aircraft has no flight details but has images, keep it as an image button
        if (!aircraft.ImageURL) {
            bottomImageLink.href = 'javascript:void(0)';
            bottomImageLink.classList.add('disabled');
            bottomImageLink.title = 'No flight details available';
            bottomImageLink.target = '_self';
        }
    }
    
    return card;
}

// Get emoji flag based on country name
function getCountryFlagEmoji(countryName) {
    // Map country names to two-letter ISO country codes for emoji flags
    const countryToISOCode = {
        'United States': 'US',
        'United Kingdom': 'GB',
        'Canada': 'CA',
        'Mexico': 'MX',
        'France': 'FR',
        'Germany': 'DE',
        'Italy': 'IT',
        'Spain': 'ES',
        'Portugal': 'PT',
        'Ireland': 'IE',
        'Belgium': 'BE',
        'Netherlands': 'NL',
        'Sweden': 'SE',
        'Denmark': 'DK',
        'Finland': 'FI',
        'Norway': 'NO',
        'Romania': 'RO',
        'Poland': 'PL',
        'Czech Republic': 'CZ',
        'Hungary': 'HU',
        'Serbia': 'RS',
        'Greece': 'GR',
        'Malta': 'MT',
        'Japan': 'JP',
        'China': 'CN',
        'India': 'IN',
        'Thailand': 'TH',
        'Indonesia': 'ID',
        'Malaysia': 'MY',
        'Singapore': 'SG',
        'Australia': 'AU',
        'New Zealand': 'NZ',
        'Argentina': 'AR',
        'Brazil': 'BR',
        'Chile': 'CL',
        'Colombia': 'CO',
        'Israel': 'IL',
        'Turkey': 'TR',
        'Egypt': 'EG',
        'South Africa': 'ZA',
        'Ethiopia': 'ET',
        'Nigeria': 'NG',
        'Algeria': 'DZ',
        'Morocco': 'MA',
        'Saudi Arabia': 'SA',
        'United Arab Emirates': 'AE',
        'Qatar': 'QA',
        'Bahrain': 'BH',
        'Iran': 'IR',
        'Iraq': 'IQ',
        'Kuwait': 'KW',
        // Adding more countries from the Wikipedia list of aircraft registration prefixes
        'Tunisia': 'TN',
        'Libya': 'LY',
        'Jordan': 'JO',
        'Lebanon': 'LB',
        'Syria': 'SY',
        'Yemen': 'YE',
        'Oman': 'OM',
        'Afghanistan': 'AF',
        'Pakistan': 'PK',
        'Bangladesh': 'BD',
        'Nepal': 'NP',
        'Sri Lanka': 'LK',
        'Myanmar': 'MM',
        'Vietnam': 'VN',
        'Philippines': 'PH',
        'Taiwan': 'TW',
        'South Korea': 'KR',
        'North Korea': 'KP',
        'Mongolia': 'MN',
        'Kazakhstan': 'KZ',
        'Uzbekistan': 'UZ',
        'Turkmenistan': 'TM',
        'Tajikistan': 'TJ',
        'Kyrgyzstan': 'KG',
        'Azerbaijan': 'AZ',
        'Armenia': 'AM',
        'Georgia': 'GE',
        'Ukraine': 'UA',
        'Belarus': 'BY',
        'Russia': 'RU',
        'Moldova': 'MD',
        'San Marino': 'SM',
        'Bulgaria': 'BG',
        'Slovakia': 'SK',
        'Austria': 'AT',
        'Switzerland': 'CH',
        'Luxembourg': 'LU',
        'Latvia': 'LV',
        'Lithuania': 'LT',
        'Estonia': 'EE',
        'Iceland': 'IS',
        'Croatia': 'HR',
        'Slovenia': 'SI',
        'Bosnia and Herzegovina': 'BA',
        'North Macedonia': 'MK',
        'Albania': 'AL',
        'Montenegro': 'ME',
        'Cyprus': 'CY',
        'Peru': 'PE',
        'Ecuador': 'EC',
        'Bolivia': 'BO',
        'Paraguay': 'PY',
        'Uruguay': 'UY',
        'Venezuela': 'VE',
        'Panama': 'PA',
        'Costa Rica': 'CR',
        'Guatemala': 'GT',
        'El Salvador': 'SV',
        'Honduras': 'HN',
        'Nicaragua': 'NI',
        'Belize': 'BZ',
        'Jamaica': 'JM',
        'Cuba': 'CU',
        'Dominican Republic': 'DO',
        'Kenya': 'KE',
        'Tanzania': 'TZ',
        'Uganda': 'UG',
        'Rwanda': 'RW',
        'Ghana': 'GH',
        'Senegal': 'SN',
        'Ivory Coast': 'CI',
        'Cameroon': 'CM',
        'Zimbabwe': 'ZW',
        'Mozambique': 'MZ',
        'Angola': 'AO',
        'Mauritius': 'MU',
        'Seychelles': 'SC',
        'Madagascar': 'MG',
        'Fiji': 'FJ',
        'Papua New Guinea': 'PG',
        'Solomon Islands': 'SB'
    };
    
    // Extract country code from "Unknown (X)" format
    if (countryName.startsWith('Unknown (') && countryName.endsWith(')')) {
        return '';  // Return empty string instead of generic flag
    }
    
    const countryCode = countryToISOCode[countryName];
    if (!countryCode) return '';  // Return empty string for unknown countries
    
    // Convert ISO country code to country flag emoji
    // Regional Indicator Symbols are Unicode characters U+1F1E6 to U+1F1FF
    // which represent
    const offset = 127397; // Offset to convert ASCII letter to Regional Indicator Symbol
    const codePoints = [...countryCode].map(char => char.charCodeAt(0) + offset);
    return String.fromCodePoint(...codePoints);
}

// Add notification icons based on aircraft notification status
function addNotificationIcons(container, aircraft) {
    // Always hide the container - this removes notification icons from the aircraft cards
    container.style.display = 'none';
    return;
    
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
    const countdownElement = document.getElementById('nextUpdateCountdown');
    
    if (!lastUpdateTime) {
        countdownElement.textContent = '-';
        return;
    }
    
    // Calculate time passed since last update
    const now = new Date();
    const elapsedSeconds = Math.floor((now - lastUpdateTime) / 1000);
    const secondsLeft = REFRESH_PERIOD - elapsedSeconds;
    
    // Format the countdown
    if (secondsLeft <= 0) {
        countdownElement.textContent = 'Refreshing...';
        countdownElement.classList.add('refreshing');
        countdownElement.parentElement.title = 'Refreshing data...';
    } else {
        const minutes = Math.floor(secondsLeft / 60);
        const seconds = secondsLeft % 60;
        
        countdownElement.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;
        countdownElement.classList.remove('refreshing');
        
        // Update tooltip with more detailed information
        if (minutes > 0) {
            countdownElement.parentElement.title = `Next refresh in ${minutes} minute${minutes !== 1 ? 's' : ''} and ${seconds} second${seconds !== 1 ? 's' : ''}`;
        } else {
            countdownElement.parentElement.title = `Next refresh in ${seconds} second${seconds !== 1 ? 's' : ''}`;
        }
    }
}

// Fetch version information from the API
async function fetchVersionInfo() {
    try {
        const response = await fetch('/api/version');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const versionData = await response.json();
        appVersion = versionData.version || 'dev';
        
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

// Initialize map link with coordinates
async function initializeMapLink() {
    try {
        // Fetch coordinates from the API
        const response = await fetch('/api/config');
        
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        
        const config = await response.json();
        
        // Store coordinates for future use
        if (config.Latitude && config.Longitude) {
            configCoordinates.latitude = config.Latitude;
            configCoordinates.longitude = config.Longitude;
            
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

// Parse and clean airport name, handling placeholder values
function cleanAirportName(name) {
    // Check for placeholder values and other invalid patterns
    if (!name || name === 'null' || name === '@@@@@' || name.includes('@@@')) {
        return 'Unknown Airport';
    }
    return name;
}

// Generate consistent flight progress percentage based on ICAO or callsign
function getConsistentFlightProgress(aircraft) {
    // Use ICAO or callsign as a seed for consistent position
    const seed = aircraft.ICAO || aircraft.Callsign || '';
    
    // Create a simple hash from the seed string
    let hash = 0;
    for (let i = 0; i < seed.length; i++) {
        hash = ((hash << 5) - hash) + seed.charCodeAt(i);
        hash = hash & hash; // Convert to 32bit integer
    }
    
    // Generate a value between 10 and 90 (avoid extremes)
    return Math.abs(hash % 81) + 10; 
}

// Flight details popup handler
function showFlightDetails(aircraft) {
    // Create overlay
    let overlay = document.querySelector('.flight-details-overlay');
    if (!overlay) {
        overlay = document.createElement('div');
        overlay.className = 'flight-details-overlay';
        document.body.appendChild(overlay);
    }

    // Create popup if it doesn't exist
    let popup = document.querySelector('.flight-details-popup');
    if (!popup) {
        popup = document.createElement('div');
        popup.className = 'flight-details-popup';
        document.body.appendChild(popup);
    }

    // Parse origin data
    let originDisplay = 'Unknown';
    let originName = '';
    let originCity = '';
    let originDetails = '';
    let originCountryCode = '';
    let originIataCode = 'N/A';
    let originIcaoCode = '';
    let originCountryName = '';
    let originCoordinates = null;
    
    if (aircraft.Origin) {
        try {
            const origin = typeof aircraft.Origin === 'string' 
                ? JSON.parse(aircraft.Origin) 
                : aircraft.Origin;
            
            if (origin) {
                originName = cleanAirportName(origin.name);
                originCity = origin.municipality || '';
                originCountryName = origin.country_name || '';
                
                originDisplay = originName;
                originDetails = originCity ? (originCountryName ? `${originCity}, ${originCountryName}` : originCity) : originCountryName;
                originCountryCode = origin.country_iso_name || '';
                originIataCode = origin.iata_code || (origin.icao_code ? origin.icao_code.substring(1) : '?');
                originIcaoCode = origin.icao_code || '';
                
                // Get coordinates if available
                if (origin.longitude_deg !== undefined && origin.latitude_deg !== undefined) {
                    originCoordinates = {
                        lat: origin.latitude_deg,
                        lng: origin.longitude_deg
                    };
                }
            }
        } catch (e) {
            console.error('Error parsing origin data:', e);
            originDisplay = String(aircraft.Origin);
        }
    }
    
    // Parse destination data
    let destinationDisplay = 'Unknown';
    let destinationName = '';
    let destinationCity = '';
    let destinationDetails = '';
    let destinationCountryCode = '';
    let destinationIataCode = 'N/A';
    let destinationIcaoCode = '';
    let destinationCountryName = '';
    let destinationCoordinates = null;
    
    if (aircraft.Destination) {
        try {
            const destination = typeof aircraft.Destination === 'string' 
                ? JSON.parse(aircraft.Destination) 
                : aircraft.Destination;
            
            if (destination) {
                destinationName = cleanAirportName(destination.name);
                destinationCity = destination.municipality || '';
                destinationCountryName = destination.country_name || '';
                
                destinationDisplay = destinationName;
                destinationDetails = destinationCity ? (destinationCountryName ? `${destinationCity}, ${destinationCountryName}` : destinationCity) : destinationCountryName;
                destinationCountryCode = destination.country_iso_name || '';
                destinationIataCode = destination.iata_code || (destination.icao_code ? destination.icao_code.substring(1) : '?');
                destinationIcaoCode = destination.icao_code || '';
                
                // Get coordinates if available
                if (destination.longitude_deg !== undefined && destination.latitude_deg !== undefined) {
                    destinationCoordinates = {
                        lat: destination.latitude_deg,
                        lng: destination.longitude_deg
                    };
                }
            }
        } catch (e) {
            console.error('Error parsing destination data:', e);
            destinationDisplay = String(aircraft.Destination);
        }
    }
    
    // Parse airline data
    let airlineDisplay = 'Unknown';
    let airlineDetails = '';
    let airlineCountryCode = '';
    let airlineCodeDetails = [];
    
    if (aircraft.Airline) {
        try {
            const airline = typeof aircraft.Airline === 'string' 
                ? JSON.parse(aircraft.Airline) 
                : aircraft.Airline;
            
            if (airline && airline.name) {
                airlineDisplay = airline.name;
                if (airline.iata) airlineCodeDetails.push(`IATA: ${airline.iata}`);
                if (airline.icao) airlineCodeDetails.push(`ICAO: ${airline.icao}`);
                airlineDetails = airlineCodeDetails.join(' | ');
                airlineCountryCode = airline.country_iso || '';
            }
        } catch (e) {
            console.error('Error parsing airline data:', e);
            airlineDisplay = String(aircraft.Airline);
        }
    }
    
    // Create Google Maps URLs
    const originMapsUrl = originCoordinates 
        ? `https://www.google.com/maps?q=${originCoordinates.lat},${originCoordinates.lng}` 
        : null;
        
    const destinationMapsUrl = destinationCoordinates 
        ? `https://www.google.com/maps?q=${destinationCoordinates.lat},${destinationCoordinates.lng}` 
        : null;

    // Get country flags
    const originFlag = createFlagEmoji(originCountryCode);
    const destinationFlag = createFlagEmoji(destinationCountryCode);
    const airlineFlag = createFlagEmoji(airlineCountryCode);

    // Always position the flight progress animation in the middle (50%)
    const flightProgressPercent = 50;
    
    // Estimate distance between airports if we have coordinates
    let flightDistanceText = '';
    if (originCoordinates && destinationCoordinates) {
        // Using Haversine formula for a rough distance calculation
        const R = 6371; // Radius of the Earth in km
        const dLat = (destinationCoordinates.lat - originCoordinates.lat) * Math.PI / 180;
        const dLon = (destinationCoordinates.lng - originCoordinates.lng) * Math.PI / 180;
        const a = 
            Math.sin(dLat/2) * Math.sin(dLat/2) +
            Math.cos(originCoordinates.lat * Math.PI / 180) * Math.cos(destinationCoordinates.lat * Math.PI / 180) * 
            Math.sin(dLon/2) * Math.sin(dLon/2);
        const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
        const distance = Math.round(R * c);
        
        flightDistanceText = `${distance.toLocaleString()} km`;
    }
    
    // Determine flight time estimate (very rough approximation)
    const flightTimeText = flightDistanceText ? `~${Math.max(1, Math.round(parseInt(flightDistanceText) / 800))}h ${Math.round(Math.random() * 59)}m` : '';
    
    // Flight number/callsign
    const flightNumberDisplay = aircraft.Callsign || 'Unknown';
    
    // Registration with formatting
    let registrationDisplay = aircraft.Registration || 'Unknown';
    if (registrationDisplay && !registrationDisplay.includes('-') && registrationDisplay.length > 2) {
        // For registrations like N12345 that don't have hyphens, keep as is
        if (!registrationDisplay.startsWith('N')) {
            // For other registration formats, try to add hyphen for better readability
            // This is a simple heuristic and might not work for all registration formats
            registrationDisplay = `${registrationDisplay.substring(0, 2)}-${registrationDisplay.substring(2)}`;
        }
    }
    
    // SVG icons
    const planeIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
        <path d="M21,16v-2l-8-5V3.5C13,2.67,12.33,2,11.5,2S10,2.67,10,3.5V9l-8,5v2l8-2.5V19l-2,1.5V22l3.5-1l3.5,1v-1.5L13,19v-5.5L21,16z"/>
    </svg>`;
    
    const locationIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8 16s6-5.686 6-10A6 6 0 0 0 2 6c0 4.314 6 10 6 10m0-7a3 3 0 1 1 0-6 3 3 0 0 1 0 6"/>
    </svg>`;
    
    const distanceIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path fill-rule="evenodd" d="M1 11.5a.5.5 0 0 0 .5.5h11.793l-3.147 3.146a.5.5 0 0 0 .708.708l4-4a.5.5 0 0 0 0-.708l-4-4a.5.5 0 1 1-.708.708L13.293 11H1.5a.5.5 0 0 0-.5.5m14-7a.5.5 0 0 1-.5.5H2.707l3.147 3.146a.5.5 0 1 1-.708.708l-4-4a.5.5 0 0 1 0-.708l4-4a.5.5 0 1 1 .708.708L2.707 4H14.5a.5.5 0 0 1 .5.5"/>
    </svg>`;
    
    const timeIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
        <path d="M8 3.5a.5.5 0 0 0-1 0V9a.5.5 0 0 0 .252.434l3.5 2a.5.5 0 0 0 .496-.868L8 8.71V3.5z"/>
        <path d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16m7-8A7 7 0 1 1 1 8a7 7 0 0 1 14 0"/>
    </svg>`;

    // Populate popup with flight details
    popup.innerHTML = `
        <div class="flight-details-header">
            <div class="flight-details-title">
                ${planeIcon}
                <span class="callsign">${aircraft.Callsign}</span>
                <span class="aircraft-type">${aircraft.Type || 'Unknown Aircraft'}</span>
            </div>
            <button class="flight-details-close" aria-label="Close">&times;</button>
        </div>
        <div class="flight-details-route">
            <div class="flight-route-line">
                <div class="route-airport origin">
                    <div class="airport-flag">
                        ${originFlag || ''}
                    </div>
                    <div class="airport-code">${originIataCode}</div>
                    <div class="airport-name">${originName}</div>
                    <div class="airport-city">${originCity}</div>
                    <div class="airport-location-icon">
                        ${originCoordinates ? `<a href="${originMapsUrl}" target="_blank" class="airport-map-link" title="View on Google Maps">${locationIcon}</a>` : ''}
                    </div>
                </div>
                <div class="route-line">
                    <div class="route-progress" style="width: ${flightProgressPercent}%"></div>
                    <div class="route-aircraft" style="left: ${flightProgressPercent}%">
                        ${planeIcon}
                    </div>
                </div>
                <div class="route-airport destination">
                    <div class="airport-flag">
                        ${destinationFlag || ''}
                    </div>
                    <div class="airport-code">${destinationIataCode}</div>
                    <div class="airport-name">${destinationName}</div>
                    <div class="airport-city">${destinationCity}</div>
                    <div class="airport-location-icon">
                        ${destinationCoordinates ? `<a href="${destinationMapsUrl}" target="_blank" class="airport-map-link" title="View on Google Maps">${locationIcon}</a>` : ''}
                    </div>
                </div>
            </div>
            ${(flightDistanceText || flightTimeText) ? `
                <div class="flight-info-badge">
                    ${flightDistanceText ? `
                        <span title="Flight distance">
                            ${distanceIcon} 
                            ${flightDistanceText}
                        </span>
                    ` : ''}
                    ${flightTimeText && flightDistanceText ? ' • ' : ''}
                    ${flightTimeText ? `
                        <span title="Estimated flight time">
                            ${timeIcon} 
                            ${flightTimeText}
                        </span>
                    ` : ''}
                </div>
            ` : ''}
        </div>
        <div class="flight-details-content">
            <div class="flight-details-section flight-details-airline">
                <div class="flight-details-label">Airline</div>
                <div class="flight-details-value">
                    ${airlineFlag ? `<span class="country-flag">${airlineFlag}</span>` : ''}
                    ${airlineDisplay}
                </div>
                ${airlineDetails ? `<div class="flight-details-subvalue">${airlineDetails}</div>` : ''}
            </div>
            <div class="flight-details-section flight-details-flight">
                <div class="flight-details-label">Flight</div>
                <div class="flight-details-value">${flightNumberDisplay}</div>
                <div class="flight-details-subvalue">Callsign</div>
            </div>
            <div class="flight-details-section flight-details-aircraft">
                <div class="flight-details-label">Aircraft</div>
                <div class="flight-details-value">${aircraft.Type || 'Unknown'}</div>
                <div class="flight-details-subvalue">${aircraft.Description || ''}</div>
            </div>
            <div class="flight-details-section flight-details-registration">
                <div class="flight-details-label">Registration</div>
                <div class="flight-details-value">${registrationDisplay}</div>
                <div class="flight-details-subvalue">${aircraft.Country || 'Unknown'}</div>
            </div>
            <div class="flight-details-section flight-details-distance">
                <div class="flight-details-label">Current Position</div>
                <div class="flight-details-value">${aircraft.Distance || '?'} km from you</div>
                <div class="flight-details-subvalue">Altitude: ${Math.round(aircraft.Altitude).toLocaleString() || '?'} ft • Speed: ${aircraft.Speed || '?'} knots</div>
            </div>
        </div>
        <div class="flight-details-actions">
            <a href="${aircraft.TrackerURL || 'javascript:void(0)'}" target="_blank" class="action-button track-aircraft-button ${!aircraft.TrackerURL ? 'disabled' : ''}">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" stroke="currentColor" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0"></path>
                    <path d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"></path>
                    <path d="M12 9l0 -2"></path>
                    <path d="M12 15l0 2"></path>
                    <path d="M9 12l-2 0"></path>
                    <path d="M15 12l2 0"></path>
                </svg>
                Track Aircraft
            </a>
            ${aircraft.ImageURL ? `
                <a href="${aircraft.ImageURL}" target="_blank" class="action-button view-image-button">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" stroke="currentColor" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M15 8h.01"></path>
                        <path d="M3 6a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v12a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3v-12z"></path>
                        <path d="M3 16l5 -5c.928 -.893 2.072 -.893 3 0l5 5"></path>
                        <path d="M14 14l1 -1c.928 -.893 2.072 -.893 3 0l3 3"></path>
                    </svg>
                    View Image
                </a>
            ` : ''}
        </div>
    `;

    // Add event listeners
    const closeButton = popup.querySelector('.flight-details-close');
    closeButton.addEventListener('click', closeFlightDetailsPopup);
    
    // Close when clicking outside the popup
    overlay.addEventListener('click', closeFlightDetailsPopup);

    // Show the popup and overlay with a slight delay for better animation
    setTimeout(() => {
        popup.classList.add('active');
        overlay.classList.add('active');
    }, 10);
}

// Create flag emoji from country code
function createFlagEmoji(countryCode) {
    if (!countryCode || countryCode.length !== 2) {
        return '';
    }
    
    // Convert ISO country code to country flag emoji
    // Regional Indicator Symbols are Unicode characters U+1F1E6 to U+1F1FF
    const offset = 127397; // Offset to convert ASCII letter to Regional Indicator Symbol
    const codePoints = [...countryCode.toUpperCase()].map(char => char.charCodeAt(0) + offset);
    
    try {
        return String.fromCodePoint(...codePoints);
    } catch (e) {
        console.error('Error creating flag emoji:', e);
        return '';
    }
}

// Close flight details popup
function closeFlightDetailsPopup() {
    const popup = document.querySelector('.flight-details-popup');
    const overlay = document.querySelector('.flight-details-overlay');
    
    if (popup) popup.classList.remove('active');
    if (overlay) overlay.classList.remove('active');
}

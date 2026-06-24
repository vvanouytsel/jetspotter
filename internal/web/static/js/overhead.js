// Global state for the overhead prediction page
let allCandidates = [];
let countdownInterval = null;
let lastUpdateTime = null;
let overheadMap = null;
let mapMarkers = {};
let mapLines = {};
let mapCpaMarkers = {};
let radiusCircle = null;
let siteMarker = null;
let appVersion = 'dev';

document.addEventListener('DOMContentLoaded', () => {
    initializeTheme();
    document.getElementById('themeToggle').addEventListener('click', toggleTheme);

    const userDropdownButton = document.getElementById('userDropdownButton');
    if (userDropdownButton) {
        userDropdownButton.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            const dropdown = this.closest('.user-dropdown');
            dropdown.classList.toggle('active');
            document.addEventListener('click', function closeDropdown(e) {
                if (!dropdown.contains(e.target)) {
                    dropdown.classList.remove('active');
                    document.removeEventListener('click', closeDropdown);
                }
            });
        });
    }

    updateMapLink();
    fetchVersionInfo();

    initOverheadMap();

    if (OVERHEAD_ENABLED) {
        fetchOverheadCandidates();
        setInterval(fetchOverheadCandidates, REFRESH_PERIOD * 1000);
    } else {
        document.getElementById('overhead-disabled-banner').style.display = 'block';
    }

    startCountdownTimer();
});

function initializeTheme() {
    const storedTheme = localStorage.getItem('theme');
    if (storedTheme) {
        document.documentElement.setAttribute('data-theme', storedTheme);
    } else {
        const prefersDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
        const initialTheme = prefersDarkMode ? 'dark' : 'light';
        document.documentElement.setAttribute('data-theme', initialTheme);
        localStorage.setItem('theme', initialTheme);
    }
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'light';
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    applyMapTheme();
}

function applyMapTheme() {
    if (overheadMap) {
        overheadMap.eachLayer(layer => {
            if (layer instanceof L.TileLayer) {
                overheadMap.removeLayer(layer);
            }
        });
        const theme = document.documentElement.getAttribute('data-theme') || 'light';
        if (theme === 'dark') {
            L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png', {
                attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
                maxZoom: 19
            }).addTo(overheadMap);
        } else {
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
                maxZoom: 19
            }).addTo(overheadMap);
        }
    }
}

function initOverheadMap() {
    overheadMap = L.map('overhead-map').setView([SITE_LATITUDE, SITE_LONGITUDE], 11);
    applyMapTheme();

    const siteIcon = L.divIcon({
        className: 'site-marker',
        html: '<div class="site-marker-inner"></div>',
        iconSize: [14, 14],
        iconAnchor: [7, 7]
    });
    siteMarker = L.marker([SITE_LATITUDE, SITE_LONGITUDE], { icon: siteIcon, title: 'Your location' })
        .bindPopup('Your location')
        .addTo(overheadMap);

    radiusCircle = L.circle([SITE_LATITUDE, SITE_LONGITUDE], {
        radius: OVERHEAD_RADIUS_KM * 1000,
        color: '#888',
        weight: 1,
        fillColor: '#888',
        fillOpacity: 0.05,
        dashArray: '4 4'
    }).addTo(overheadMap);
}

async function fetchOverheadCandidates() {
    try {
        const response = await fetch('/api/overhead');
        if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
        const newCandidates = await response.json();
        allCandidates = newCandidates || [];
        lastUpdateTime = new Date();
        renderOverheadMap();
        renderOverheadTable();
        updateSummary();
    } catch (error) {
        console.error('Error fetching overhead candidates:', error);
        lastUpdateTime = new Date();
    }
}

function stateColor(state) {
    switch (state) {
        case 'confirmed': return '#22c55e';
        case 'pending':   return '#f59e0b';
        case 'passed':    return '#6b7280';
        case 'abandoned': return '#ef4444';
        default:          return '#6b7280';
    }
}

function aircraftIcon(state, heading) {
    const color = stateColor(state);
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24" fill="${color}" stroke="#fff" stroke-width="1"><path d="M12 2L8 12 2 14l6 2-2 6 6-4 6 4-2-6 6-2-6-2z" transform="rotate(0 12 12)"/></svg>`;
    return L.divIcon({
        className: 'aircraft-map-marker',
        html: `<div style="transform: rotate(${heading}deg);">${svg}</div>`,
        iconSize: [24, 24],
        iconAnchor: [12, 12]
    });
}

function cpaIcon() {
    return L.divIcon({
        className: 'cpa-marker',
        html: '<div class="cpa-marker-inner"></div>',
        iconSize: [12, 12],
        iconAnchor: [6, 6]
    });
}

function advancePosition(lat, lon, headingDeg, km) {
    const kmPerDeg = 111.32;
    const brg = headingDeg * Math.PI / 180;
    const dLat = (km * Math.cos(brg)) / kmPerDeg;
    const cosLat = Math.cos(lat * Math.PI / 180) || 1e-6;
    const dLon = (km * Math.sin(brg)) / (kmPerDeg * cosLat);
    return [lat + dLat, lon + dLon];
}

function renderOverheadMap() {
    if (!overheadMap) return;

    const seenIcao = new Set();

    allCandidates.forEach(c => {
        seenIcao.add(c.icao);
        const ac = c.aircraft || {};
        const lat = ac.Latitude || ac.latitude;
        const lon = ac.Longitude || ac.longitude;
        if (lat == null || lon == null) return;

        const heading = ac.Heading || ac.heading || 0;
        const state = c.confirmationState || 'pending';

        // Aircraft marker
        if (mapMarkers[c.icao]) {
            mapMarkers[c.icao].setLatLng([lat, lon]);
            mapMarkers[c.icao].setIcon(aircraftIcon(state, heading));
        } else {
            mapMarkers[c.icao] = L.marker([lat, lon], { icon: aircraftIcon(state, heading) })
                .bindPopup(buildPopup(c))
                .addTo(overheadMap);
        }
        mapMarkers[c.icao].setPopupContent(buildPopup(c));

        // Trajectory line up to the predicted CPA point
        const linePoints = computeTrajectoryToCPA(c, lat, lon, heading);
        if (linePoints && linePoints.length > 1) {
            if (mapLines[c.icao]) {
                mapLines[c.icao].setLatLngs(linePoints);
            } else {
                mapLines[c.icao] = L.polyline(linePoints, {
                    color: stateColor(state),
                    weight: 2,
                    opacity: 0.7,
                    dashArray: '5 5'
                }).addTo(overheadMap);
            }
        }

        // CPA marker
        if (linePoints && linePoints.length > 1) {
            const cpaPoint = linePoints[linePoints.length - 1];
            if (mapCpaMarkers[c.icao]) {
                mapCpaMarkers[c.icao].setLatLng(cpaPoint);
            } else {
                mapCpaMarkers[c.icao] = L.marker(cpaPoint, { icon: cpaIcon() })
                    .bindTooltip('Predicted overhead point')
                    .addTo(overheadMap);
            }
        }
    });

    // Remove markers for candidates that disappeared
    Object.keys(mapMarkers).forEach(icao => {
        if (!seenIcao.has(icao)) {
            overheadMap.removeLayer(mapMarkers[icao]);
            delete mapMarkers[icao];
            if (mapLines[icao]) { overheadMap.removeLayer(mapLines[icao]); delete mapLines[icao]; }
            if (mapCpaMarkers[icao]) { overheadMap.removeLayer(mapCpaMarkers[icao]); delete mapCpaMarkers[icao]; }
        }
    });
}

function computeTrajectoryToCPA(c, lat, lon, heading) {
    const ac = c.aircraft || {};
    const speedKnots = ac.Speed || ac.speed || 0;
    if (speedKnots <= 0) return null;
    const kmph = speedKnots * 1.852;
    const kmPerSec = kmph / 3600;

    const minutesLeft = c.minutesUntilOverhead || 0;
    const totalSec = minutesLeft * 60;
    if (totalSec <= 0) return [[lat, lon]];

    const stepSec = 5;
    const steps = Math.min(Math.ceil(totalSec / stepSec), 240);
    const points = [[lat, lon]];
    let curLat = lat, curLon = lon;
    for (let i = 1; i <= steps; i++) {
        const km = kmPerSec * stepSec;
        const [nLat, nLon] = advancePosition(curLat, curLon, heading, km);
        points.push([nLat, nLon]);
        curLat = nLat; curLon = nLon;
    }
    return points;
}

function buildPopup(c) {
    const ac = c.aircraft || {};
    const callsign = ac.Callsign || ac.callsign || 'Unknown';
    const type = ac.Type || ac.type || '';
    const heading = Math.round(ac.Heading || ac.heading || 0);
    const minutes = c.minutesUntilOverhead;
    const headline = minutes <= 0
        ? `Look up now! A ${type || 'aircraft'} is passing overhead heading ${c.compassWord}.`
        : `In ${minutes} min, look up — you'll see a ${type || 'aircraft'} heading ${c.compassWord}.`;
    return `<div style="min-width:180px;">
        <strong>${callsign}</strong><br>
        ${type}<br>
        Heading: ${c.compassWord} (${heading}°)<br>
        CPA distance: ${c.cpaDistanceKm} km<br>
        State: ${c.confirmationState}
    </div>`;
}

function renderOverheadTable() {
    const tbody = document.getElementById('overhead-table-body');
    const noCandidates = document.getElementById('noCandidatesMessage');
    tbody.innerHTML = '';

    if (!allCandidates || allCandidates.length === 0) {
        noCandidates.style.display = 'flex';
        return;
    }
    noCandidates.style.display = 'none';

    const sorted = [...allCandidates].sort((a, b) => {
        const aT = a.predictedCPATime || '';
        const bT = b.predictedCPATime || '';
        if (aT < bT) return -1;
        if (aT > bT) return 1;
        return 0;
    });

    sorted.forEach(c => {
        const ac = c.aircraft || {};
        const row = document.createElement('tr');
        row.className = `conf-${c.confirmationState}`;
        row.innerHTML = `
            <td>${ac.Callsign || ac.callsign || 'Unknown'}</td>
            <td>${ac.Type || ac.type || ''}</td>
            <td>${c.compassWord || ''} (${Math.round(ac.Heading || ac.heading || 0)}°)</td>
            <td>${formatAltitude(ac)}</td>
            <td>${ac.Speed || ac.speed || ''} kn</td>
            <td>${ac.Distance != null ? ac.Distance : ''} km</td>
            <td>${c.cpaDistanceKm} km</td>
            <td>${c.minutesUntilOverhead <= 0 ? 'now' : c.minutesUntilOverhead + ' min'}</td>
            <td><span class="state-badge conf-${c.confirmationState}">${c.confirmationState}</span></td>
            <td>${ac.TrackerURL ? `<a href="${ac.TrackerURL}" target="_blank">track</a>` : ''}</td>
        `;
        tbody.appendChild(row);
    });
}

function formatAltitude(ac) {
    const alt = ac.Altitude != null ? ac.Altitude : ac.altitude;
    if (ac.OnGround || ac.onGround) return 'On ground';
    if (alt == null) return '';
    return `${Math.round(alt).toLocaleString()} ft`;
}

function updateSummary() {
    let confirmed = 0, pending = 0, passed = 0;
    (allCandidates || []).forEach(c => {
        if (c.confirmationState === 'confirmed') confirmed++;
        else if (c.confirmationState === 'pending') pending++;
        else if (c.confirmationState === 'passed') passed++;
    });
    document.getElementById('confirmedCount').textContent = confirmed;
    document.getElementById('pendingCount').textContent = pending;
    document.getElementById('passedCount').textContent = passed;
}

function startCountdownTimer() {
    if (countdownInterval) clearInterval(countdownInterval);
    lastUpdateTime = new Date();
    updateCountdown();
    countdownInterval = setInterval(updateCountdown, 1000);
}

function updateCountdown() {
    const el = document.getElementById('nextUpdateCountdown');
    if (!el) return;
    if (!lastUpdateTime) { el.textContent = '-'; return; }
    const now = new Date();
    const elapsed = Math.floor((now - lastUpdateTime) / 1000);
    const left = REFRESH_PERIOD - elapsed;
    if (left <= 0) {
        el.textContent = 'Refreshing...';
        el.classList.add('refreshing');
    } else {
        const m = Math.floor(left / 60);
        const s = left % 60;
        el.textContent = `${m}:${s.toString().padStart(2, '0')}`;
        el.classList.remove('refreshing');
    }
}

async function fetchVersionInfo() {
    try {
        const response = await fetch('/api/version');
        if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
        const v = await response.json();
        appVersion = v.version || 'dev';
        const el = document.getElementById('appVersion');
        if (el) el.textContent = appVersion;
    } catch (e) {
        const el = document.getElementById('appVersion');
        if (el) el.textContent = 'dev';
    }
}

function updateMapLink() {
    const link = document.getElementById('mapLink');
    if (!link) return;
    if (typeof SITE_LATITUDE !== 'undefined' && typeof SITE_LONGITUDE !== 'undefined') {
        link.href = `https://globe.airplanes.live/?lat=${SITE_LATITUDE}&lon=${SITE_LONGITUDE}&SiteLat=${SITE_LATITUDE}&SiteLon=${SITE_LONGITUDE}&zoom=11&enableLabels&extendedLabels=1&hideSidebar`;
    } else {
        link.href = 'https://globe.airplanes.live/';
    }
}

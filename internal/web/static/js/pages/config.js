// The configuration page: shows what Jetspotter is currently set to.
import { currentTheme, initShell } from '../shell.js';
import { html, toString } from '../lib/dom.js';

const root = document.getElementById('config-content');

const NOTIFICATION_COLOR = '#e5484d';
const SCAN_COLOR = '#3b82f6';

const TILES = {
    light: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
    dark: 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png',
};
const TILE_ATTRIBUTION = '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors';

// Text ----------------------------------------------------------------------

function describeAircraftTypes(types = []) {
    if (types.length === 0) return 'No specific types configured';
    if (types.includes('ALL')) return 'All aircraft';
    if (types.includes('MILITARY')) {
        const others = types.filter((t) => t !== 'MILITARY');
        return others.length === 0 ? 'Military aircraft only' : `Military aircraft and ${others.join(', ')}`;
    }
    return types.join(', ');
}

// Never show webhook URLs or tokens here, only whether a service is set up.
function notificationServices(config) {
    return [
        { name: 'Terminal', active: true },
        {
            name: 'Discord',
            active: Boolean(config.DiscordWebHookURL),
            detail: config.DiscordWebHookURL ? `Color by altitude: ${config.DiscordColorAltitude === 'true' ? 'yes' : 'no'}` : '',
        },
        { name: 'Slack', active: Boolean(config.SlackWebHookURL) },
        { name: 'Gotify', active: Boolean(config.GotifyURL && config.GotifyToken) },
        {
            name: 'ntfy',
            active: Boolean(config.NtfyTopic),
            detail: config.NtfyTopic ? `Topic ${config.NtfyTopic} on ${config.NtfyServer}` : '',
        },
    ];
}

// Drawing -------------------------------------------------------------------

function setting(name, value, hint) {
    return html`<div class="setting">
        <dt>${name}</dt>
        <dd>${value}<span class="setting__hint">${hint}</span></dd>
    </div>`;
}

function panel(title, body) {
    return html`<section class="panel">
        <h2 class="panel__title">${title}</h2>
        <div class="panel__body">${body}</div>
    </section>`;
}

function renderConfig(config, scanRange) {
    const noLimit = config.MaxAltitudeFeet === 0;

    return html`
        ${panel('Location', html`
            <div class="location">
                <dl class="settings">
                    ${setting('Latitude', config.Location.Lat)}
                    ${setting('Longitude', config.Location.Lon)}
                </dl>
                <div>
                    <div class="map" id="map"></div>
                    <div class="legend">
                        <span class="legend__item"><span class="legend__swatch" style="--swatch: ${NOTIFICATION_COLOR}"></span>Notification range</span>
                        <span class="legend__item"><span class="legend__swatch" style="--swatch: ${SCAN_COLOR}"></span>Scan range</span>
                    </div>
                </div>
            </div>`)}

        ${panel('Tracking', html`<dl class="settings">
            ${setting('Scan range', `${scanRange} km`, 'How far from your location to look for aircraft.')}
            ${setting('Maximum altitude', noLimit ? 'No limit' : `${config.MaxAltitudeFeet} ft`, 'The highest altitude aircraft are tracked at.')}
            ${setting('Update interval', `${config.FetchInterval} seconds`, 'How often new aircraft data is fetched.')}
        </dl>`)}

        ${panel('Notifications', html`<dl class="settings">
            ${setting('Notification range', `${config.MaxRangeKilometers} km`, 'Only aircraft inside this distance trigger a notification.')}
            ${setting('Aircraft types', describeAircraftTypes(config.AircraftTypes), 'The kinds of aircraft you get notified about.')}
        </dl>`)}

        ${panel('Notification services', html`<div class="services">
            ${notificationServices(config).map(
                (service) => html`<div class="service" ${service.active ? 'data-active' : ''}>
                    <span class="service__name">${service.name}</span>
                    <span class="service__status">${service.active ? 'Active' : 'Not configured'}</span>
                    ${service.detail ? html`<span class="service__detail">${service.detail}</span>` : ''}
                </div>`,
            )}
        </div>`)}
    `;
}

// Map -----------------------------------------------------------------------

function drawMap(config, scanRange) {
    if (typeof L === 'undefined') {
        document.getElementById('map').textContent = 'The map could not be loaded.';
        return;
    }

    const center = [config.Location.Lat, config.Location.Lon];
    const map = L.map('map').setView(center, 10);

    let tiles = null;
    const showTiles = () => {
        tiles?.remove();
        tiles = L.tileLayer(TILES[currentTheme()], { attribution: TILE_ATTRIBUTION }).addTo(map);
    };
    showTiles();
    document.addEventListener('themechange', showTiles);

    L.marker(center).addTo(map).bindPopup(`<b>Tracking location</b><br>${center[0]}, ${center[1]}`);

    const notificationCircle = L.circle(center, {
        color: NOTIFICATION_COLOR,
        fillColor: NOTIFICATION_COLOR,
        fillOpacity: 0.2,
        radius: config.MaxRangeKilometers * 1000,
    }).addTo(map);

    let outer = notificationCircle;
    if (scanRange > config.MaxRangeKilometers) {
        outer = L.circle(center, {
            color: SCAN_COLOR,
            fillColor: SCAN_COLOR,
            fillOpacity: 0.1,
            radius: scanRange * 1000,
        }).addTo(map);
    }
    map.fitBounds(outer.getBounds());
}

// Start ---------------------------------------------------------------------

async function start() {
    try {
        const response = await fetch('/api/config');
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const config = await response.json();

        const scanRange = config.MaxScanRangeKilometers || config.MaxRangeKilometers;
        root.innerHTML = toString(renderConfig(config, scanRange));
        drawMap(config, scanRange);
    } catch (error) {
        console.error('Could not load the configuration:', error);
        root.innerHTML = toString(html`<div class="notice">Could not load the configuration. Refresh the page to try again.</div>`);
    }
}

initShell();
start();

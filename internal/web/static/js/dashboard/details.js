// The "flight details" dialog: route, airline and aircraft information.
import { html, icon, toString } from '../lib/dom.js';
import { cleanName, distanceKm, flagFromIso, formatNumber, EMPTY } from '../lib/format.js';

/** @typedef {import('../lib/types.js').Aircraft} Aircraft */
/** @typedef {import('../lib/types.js').Airport} Airport */

const dialog = () => /** @type {HTMLDialogElement} */ (document.getElementById('details-dialog'));

function hasCoordinates(airport) {
    return (
        Number.isFinite(airport.latitude) &&
        Number.isFinite(airport.longitude) &&
        (airport.latitude !== 0 || airport.longitude !== 0)
    );
}

/** @param {Airport} airport */
function renderAirport(airport) {
    const code = airport.iata_code || airport.icao_code || EMPTY;
    const city = [airport.municipality, airport.country_name].filter(Boolean).join(', ');
    return html`<div class="route__airport">
        <span class="route__code">${code}</span>
        <span class="route__name">${flagFromIso(airport.country_iso_name)} ${cleanName(airport.name) || 'Unknown airport'}</span>
        <span class="route__city">${city}</span>
    </div>`;
}

/**
 * @param {Airport} origin
 * @param {Airport} destination
 */
function renderRoute(origin, destination) {
    const known = (airport) => airport.name || airport.icao_code;
    if (!known(origin) && !known(destination)) return '';

    const km =
        hasCoordinates(origin) && hasCoordinates(destination)
            ? distanceKm(
                  { lat: origin.latitude, lon: origin.longitude },
                  { lat: destination.latitude, lon: destination.longitude },
              )
            : null;

    return html`<div class="route">
        ${renderAirport(origin)}
        <div class="route__line">
            ${icon('plane')}
            <span>${km !== null ? `${formatNumber(km)} km` : ''}</span>
        </div>
        ${renderAirport(destination)}
    </div>`;
}

/** @param {Aircraft} a */
function renderDetails(a) {
    const airline = a.Airline || {};
    const origin = a.Origin || {};
    const destination = a.Destination || {};
    const airlineCodes = [airline.iata && `IATA ${airline.iata}`, airline.icao && `ICAO ${airline.icao}`]
        .filter(Boolean)
        .join(' · ');
    const photo = a.ImageThumbnailURL || a.ImageURL;

    return html`<div class="details">
        <header class="details__hero" ${photo ? html`data-has-photo` : ''}>
            ${photo ? html`<img class="details__photo" src="${photo}" alt="">` : ''}
            <div>
                <div class="details__callsign">${a.Callsign || 'Unknown'}</div>
                <div class="details__type">${a.Type || 'Unknown aircraft'}</div>
            </div>
            <button class="icon-button details__close" type="button" data-close aria-label="Close">${icon('x')}</button>
        </header>

        ${renderRoute(origin, destination)}

        <dl class="details__grid">
            <div>
                <dt>Airline</dt>
                <dd>${flagFromIso(airline.country_iso)} ${airline.name || EMPTY}<small>${airlineCodes}</small></dd>
            </div>
            <div>
                <dt>Aircraft</dt>
                <dd>${a.Type || EMPTY}<small>${a.Description || ''}</small></dd>
            </div>
            <div>
                <dt>Registration</dt>
                <dd class="mono">${a.Registration || EMPTY}<small>${a.Country || ''}</small></dd>
            </div>
            <div>
                <dt>Position</dt>
                <dd>${formatNumber(a.Distance)} km from you
                    <small>${a.OnGround ? 'On the ground' : `${formatNumber(a.Altitude)} ft · ${formatNumber(a.Speed)} kt`}</small>
                </dd>
            </div>
        </dl>

        <footer class="details__actions">
            ${a.TrackerURL
                ? html`<a class="button button--primary" href="${a.TrackerURL}" target="_blank" rel="noopener">${icon('radar')} Track live</a>`
                : ''}
            ${a.ImageURL
                ? html`<a class="button" href="${a.ImageURL}" target="_blank" rel="noopener">${icon('image')} Photos</a>`
                : ''}
        </footer>
    </div>`;
}

/** @param {Aircraft} aircraft */
export function openDetails(aircraft) {
    const el = dialog();
    el.innerHTML = toString(renderDetails(aircraft));
    el.showModal();
}

/** Close the dialog with its close button, or by clicking outside of it. Call once. */
export function initDetails() {
    dialog().addEventListener('click', (event) => {
        const target = /** @type {HTMLElement} */ (event.target);
        if (target === dialog() || target.closest('[data-close]')) dialog().close();
    });
}

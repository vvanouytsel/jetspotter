// One aircraft card.
import { html, icon } from '../lib/dom.js';
import {
    altitudeBand,
    altitudePercent,
    cardinal,
    flagFromIso,
    formatNumber,
    isoFromCountryName,
    EMPTY,
} from '../lib/format.js';

/** @typedef {import('../lib/types.js').Aircraft} Aircraft */

function photo(a) {
    const src = a.ImageThumbnailURL || a.ImageURL;
    const image = src
        ? html`<img src="${src}" alt="${a.Type || 'Aircraft'} ${a.Registration || ''}" loading="lazy"
              ${a.Photographer ? html`title="Photo by ${a.Photographer}"` : ''}>`
        : html`<img src="/static/images/aircraft_not_found.png" alt="No photo available" data-fallback>`;

    return a.ImageURL
        ? html`<a href="${a.ImageURL}" target="_blank" rel="noopener" title="More photos of this aircraft">${image}</a>`
        : image;
}

/** The badges are drawn on top of the photo. */
function badges(a) {
    const list = [
        a.Military ? html`<span class="badge badge--military">Military</span>` : '',
        a.Inbound ? html`<span class="badge badge--inbound">Inbound</span>` : '',
        a.OnGround ? html`<span class="badge">On ground</span>` : '',
    ];
    return list.some(Boolean) ? html`<div class="card__badges">${list}</div>` : '';
}

function metric(label, value, unit, className = '') {
    return html`<div class="metric ${className}">
        <dt>${label}</dt>
        <dd>${value}${unit ? html`<small> ${unit}</small>` : ''}</dd>
    </div>`;
}

/** A small round compass whose arrow points at `degrees`. */
function dial(label, degrees, text) {
    if (!Number.isFinite(degrees)) return '';
    return html`<div class="dial" title="${label}: ${Math.round(degrees)}°">
        <span class="dial__face" style="--angle: ${degrees}deg">${icon('arrow-up')}</span>
        <span>
            <span class="dial__label">${label}</span>
            <span class="dial__value">${text}</span>
        </span>
    </div>`;
}

/**
 * Returns the HTML of a card for one aircraft.
 * @param {Aircraft} a
 */
export function renderCard(a) {
    const onGround = Boolean(a.OnGround);
    const flag = flagFromIso(isoFromCountryName(a.Country));

    return html`<article class="card" data-band="${onGround ? 0 : altitudeBand(a.Altitude)}">
        <div class="card__media">
            ${photo(a)}
            ${badges(a)}
            <div class="card__overlay">
                <div class="card__title">
                    <span class="card__callsign">${a.Callsign || 'Unknown'}</span>
                    ${flag ? html`<span class="card__flag" title="${a.Country}">${flag}</span>` : ''}
                </div>
                <div class="card__type">${a.Description || a.Type || 'Unknown aircraft'}</div>
            </div>
        </div>

        <div class="card__body">
            <dl class="metrics">
                ${metric('Altitude', onGround ? 'Ground' : formatNumber(a.Altitude), onGround ? '' : 'ft', 'metric--altitude')}
                ${metric('Speed', onGround ? EMPTY : formatNumber(a.Speed), onGround ? '' : 'kt')}
                ${metric('Distance', formatNumber(a.Distance), 'km')}
            </dl>
            <div class="altbar" title="Altitude" style="--fill: ${altitudePercent(a.Altitude)}%"><span></span></div>

            <ul class="pills">
                <li><span>Reg</span>${a.Registration || EMPTY}</li>
                <li><span>ICAO</span>${a.ICAO || EMPTY}</li>
                ${a.Country ? html`<li><span>From</span>${a.Country}</li>` : ''}
            </ul>
        </div>

        <footer class="card__footer">
            <div class="card__compass">
                ${onGround ? '' : dial('Heading', a.Heading, `${Math.round(a.Heading)}°`)}
                ${dial('From you', a.BearingFromLocation, `${cardinal(a.BearingFromLocation)} ${Math.round(a.BearingFromLocation)}°`)}
            </div>
            <div class="card__actions">
                ${a.TrackerURL
                    ? html`<a class="button button--small" href="${a.TrackerURL}" target="_blank" rel="noopener" title="Track this aircraft live">${icon('radar')} Track</a>`
                    : ''}
                <button class="button button--small button--primary" type="button" data-details="${a.ICAO}">Details</button>
            </div>
        </footer>
    </article>`;
}

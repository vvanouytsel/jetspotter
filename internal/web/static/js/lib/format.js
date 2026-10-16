// Formatting helpers: numbers, altitude, directions, countries and distances.

export const EMPTY = '–';

/** The highest altitude (in feet) the altitude bar on a card can show. */
export const MAX_BAR_ALTITUDE = 40000;

/**
 * A number with thousand separators, or a dash when there is no value.
 * @param {number | undefined} value
 */
export function formatNumber(value) {
    return Number.isFinite(value) ? Math.round(value).toLocaleString() : EMPTY;
}

/**
 * Which color band an altitude belongs to (0 to 6).
 * 0 means unknown or on the ground. The bands match the --alt-N colors in base.css.
 * @param {number | undefined} feet
 */
export function altitudeBand(feet) {
    if (!Number.isFinite(feet) || feet <= 0) return 0;
    if (feet < 2000) return 1;
    if (feet < 5000) return 2;
    if (feet < 10000) return 3;
    if (feet < 20000) return 4;
    if (feet < 30000) return 5;
    return 6;
}

/**
 * How full the altitude bar is, as a percentage.
 * @param {number | undefined} feet
 */
export function altitudePercent(feet) {
    if (!Number.isFinite(feet) || feet <= 0) return 0;
    return Math.min(100, Math.round((feet / MAX_BAR_ALTITUDE) * 100));
}

const CARDINALS = ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW'];

/**
 * The compass direction for a bearing in degrees: 0 is N, 90 is E, ...
 * @param {number} degrees
 */
export function cardinal(degrees) {
    return CARDINALS[Math.round((((degrees % 360) + 360) % 360) / 45) % 8];
}

// Countries ------------------------------------------------------------------

// The backend gives a country as an English name ("Belgium"). Names that the
// browser spells differently are listed here.
const COUNTRY_ALIASES = {
    'czech republic': 'CZ',
    'ivory coast': 'CI',
    'turkey': 'TR',
    'myanmar': 'MM',
    'russia': 'RU',
    'south korea': 'KR',
    'north korea': 'KP',
    'vietnam': 'VN',
    'laos': 'LA',
    'syria': 'SY',
    'iran': 'IR',
    'moldova': 'MD',
    'brunei': 'BN',
    'cape verde': 'CV',
    'swaziland': 'SZ',
};

/** @type {Map<string, string> | null} */
let isoByName = null;

function buildCountryIndex() {
    isoByName = new Map(Object.entries(COUNTRY_ALIASES));
    const names = new Intl.DisplayNames(['en'], { type: 'region' });
    const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    for (const a of letters) {
        for (const b of letters) {
            const code = a + b;
            try {
                const name = names.of(code);
                if (name && name !== code && !isoByName.has(name.toLowerCase())) {
                    isoByName.set(name.toLowerCase(), code);
                }
            } catch (e) {
                // Not a real region code, skip it
            }
        }
    }
}

/**
 * The two letter code (BE, NL, ...) for an English country name, or '' when the country is not known.
 * @param {string | undefined} name
 */
export function isoFromCountryName(name) {
    if (!name) return '';
    if (!isoByName) buildCountryIndex();
    return isoByName.get(name.toLowerCase()) || '';
}

/**
 * Turns a two letter country code into a flag emoji. Returns '' if the code is not valid.
 * @param {string | undefined} code
 */
export function flagFromIso(code) {
    if (typeof code !== 'string' || !/^[a-z]{2}$/i.test(code)) return '';
    const REGIONAL_INDICATOR_OFFSET = 127397;
    return String.fromCodePoint(...[...code.toUpperCase()].map((ch) => ch.charCodeAt(0) + REGIONAL_INDICATOR_OFFSET));
}

// Airports -------------------------------------------------------------------

/**
 * The route data uses placeholder text such as "@@@@@" for unknown names.
 * @param {string | undefined} name
 */
export function cleanName(name) {
    if (!name || name === 'null' || name.includes('@@@')) return '';
    return name;
}

/**
 * Straight line distance in kilometers between two coordinates.
 * @param {{lat: number, lon: number}} from
 * @param {{lat: number, lon: number}} to
 */
export function distanceKm(from, to) {
    const EARTH_RADIUS_KM = 6371;
    const toRad = (deg) => (deg * Math.PI) / 180;
    const dLat = toRad(to.lat - from.lat);
    const dLon = toRad(to.lon - from.lon);
    const a =
        Math.sin(dLat / 2) ** 2 +
        Math.cos(toRad(from.lat)) * Math.cos(toRad(to.lat)) * Math.sin(dLon / 2) ** 2;
    return Math.round(EARTH_RADIUS_KM * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a)));
}

// Which aircraft are shown, and in what order.
// This file has no knowledge of the page: it only works on data.

/** @typedef {import('../lib/types.js').Aircraft} Aircraft */

/** What the user has picked in the toolbar. */
export const defaultView = () => ({
    query: '',
    type: '',
    inbound: false,
    military: false,
    hideGround: false,
    sortBy: 'distance',
    sortDirection: 'asc',
});

/**
 * What kind of aircraft this is, used for the type filter and for sorting.
 * @param {Aircraft} aircraft
 */
export function aircraftType(aircraft) {
    return aircraft.Description || aircraft.Type || '';
}

const SORTERS = {
    distance: (a, b) => (a.Distance || 0) - (b.Distance || 0),
    altitude: (a, b) => (a.Altitude || 0) - (b.Altitude || 0),
    speed: (a, b) => (a.Speed || 0) - (b.Speed || 0),
    type: (a, b) => aircraftType(a).localeCompare(aircraftType(b)),
};

/** Does the text the user typed appear in one of the aircraft's names? */
function matchesQuery(aircraft, query) {
    if (!query) return true;
    const haystack = [aircraft.Callsign, aircraft.Registration, aircraft.ICAO, aircraft.Type, aircraft.Description, aircraft.Country]
        .join(' ')
        .toLowerCase();
    return haystack.includes(query.trim().toLowerCase());
}

/**
 * Returns a new list with the filters and the sorting of `view` applied.
 * @param {Aircraft[]} aircraft
 * @param {ReturnType<typeof defaultView>} view
 */
export function applyView(aircraft, view) {
    const shown = aircraft.filter(
        (a) =>
            matchesQuery(a, view.query) &&
            (!view.type || aircraftType(a) === view.type) &&
            (!view.inbound || a.Inbound) &&
            (!view.military || a.Military) &&
            (!view.hideGround || !a.OnGround),
    );

    const direction = view.sortDirection === 'desc' ? -1 : 1;
    return shown.sort((a, b) => direction * SORTERS[view.sortBy](a, b));
}

/**
 * The sorted list of aircraft types that are currently around.
 * @param {Aircraft[]} aircraft
 */
export function availableTypes(aircraft) {
    return [...new Set(aircraft.map(aircraftType).filter(Boolean))].sort();
}

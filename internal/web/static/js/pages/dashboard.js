// The dashboard page: fetches aircraft, keeps the filters, draws the cards.
import { initShell } from '../shell.js';
import { html, icon, toString } from '../lib/dom.js';
import { renderCard } from '../dashboard/card.js';
import { initDetails, openDetails } from '../dashboard/details.js';
import { applyView, availableTypes, defaultView } from '../dashboard/filters.js';

const root = document.getElementById('dashboard');
const refreshSeconds = Number(root.dataset.refreshSeconds) || 60;

const el = {
    grid: document.getElementById('aircraft-grid'),
    search: document.getElementById('search'),
    typeSelect: document.getElementById('filter-type'),
    sortBy: document.getElementById('sort-by'),
    sortDirection: document.getElementById('sort-direction'),
    countdown: document.getElementById('countdown'),
    resultCount: document.getElementById('result-count'),
    counts: {
        total: document.getElementById('count-total'),
        inbound: document.getElementById('count-inbound'),
        military: document.getElementById('count-military'),
    },
};

const state = {
    aircraft: [],
    view: defaultView(),
    loading: true,
    failed: false,
    nextRefreshAt: null,
};

// Drawing -------------------------------------------------------------------

function renderEmptyState() {
    if (state.loading) {
        return html`<div class="empty-state"><div class="radar" role="status"></div>
            <div><strong>Scanning the skies…</strong>Looking for aircraft in your vicinity.</div></div>`;
    }
    if (state.failed && state.aircraft.length === 0) {
        return html`<div class="empty-state"><div>${icon('info')}
            <strong>Could not load aircraft</strong>Trying again at the next refresh.</div></div>`;
    }
    const filtered = state.aircraft.length > 0;
    return html`<div class="empty-state"><div>
        <strong>${filtered ? 'No aircraft match your filters' : 'No aircraft right now'}</strong>
        ${filtered ? 'Change or clear the filters to see more.' : 'Nothing is flying within range. This page updates automatically.'}
    </div></div>`;
}

function renderGrid() {
    const shown = applyView(state.aircraft, state.view);
    el.grid.innerHTML = shown.length > 0 ? toString(shown.map(renderCard)) : toString(renderEmptyState());

    const total = state.aircraft.length;
    el.resultCount.textContent = state.loading || total === 0 ? '' : `Showing ${shown.length} of ${total} aircraft`;
}

function renderCounts() {
    el.counts.total.textContent = state.aircraft.length;
    el.counts.inbound.textContent = state.aircraft.filter((a) => a.Inbound).length;
    el.counts.military.textContent = state.aircraft.filter((a) => a.Military).length;
}

function renderTypeOptions() {
    const types = availableTypes(state.aircraft);
    // A type that disappeared from the sky should not stay selected
    if (state.view.type && !types.includes(state.view.type)) state.view.type = '';

    const options = [html`<option value="">All types</option>`, ...types.map((t) => html`<option value="${t}">${t}</option>`)];
    el.typeSelect.innerHTML = toString(options);
    el.typeSelect.value = state.view.type;
}

function renderControls() {
    document.querySelectorAll('[data-filter]').forEach((button) => {
        button.setAttribute('aria-pressed', String(Boolean(state.view[button.dataset.filter])));
    });

    const descending = state.view.sortDirection === 'desc';
    el.sortDirection.dataset.direction = state.view.sortDirection;
    el.sortDirection.setAttribute('aria-label', `Sort direction: ${descending ? 'descending' : 'ascending'}`);
}

function render() {
    renderGrid();
    renderControls();
}

// Data ----------------------------------------------------------------------

async function loadAircraft() {
    if (new URLSearchParams(location.search).get('test') === 'true') {
        return (await import('../dashboard/demo-data.js')).demoAircraft;
    }
    const response = await fetch('/api/aircraft');
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return (await response.json()) ?? [];
}

async function refresh() {
    let next = state.aircraft;
    try {
        next = await loadAircraft();
        state.failed = false;
    } catch (error) {
        console.error('Could not load aircraft:', error);
        state.failed = true;
    }

    // Only redraw when something changed, so photos do not flicker on every refresh
    const changed = state.loading || state.failed || JSON.stringify(next) !== JSON.stringify(state.aircraft);
    state.aircraft = next;
    state.loading = false;
    state.nextRefreshAt = Date.now() + refreshSeconds * 1000;

    if (changed) {
        renderCounts();
        renderTypeOptions();
        render();
    }
    updateCountdown();
    setTimeout(refresh, refreshSeconds * 1000);
}

function updateCountdown() {
    if (!state.nextRefreshAt) return;
    const secondsLeft = Math.round((state.nextRefreshAt - Date.now()) / 1000);
    if (secondsLeft <= 0) {
        el.countdown.textContent = 'Refreshing…';
        el.countdown.dataset.state = 'refreshing';
        return;
    }
    const minutes = Math.floor(secondsLeft / 60);
    el.countdown.textContent = `${minutes}:${String(secondsLeft % 60).padStart(2, '0')}`;
    el.countdown.title = `Next refresh in ${secondsLeft} seconds`;
    delete el.countdown.dataset.state;
}

// Events --------------------------------------------------------------------

function initEvents() {
    // Chips and tiles toggle a yes/no filter
    document.querySelectorAll('[data-filter]').forEach((button) => {
        button.addEventListener('click', () => {
            const name = button.dataset.filter;
            state.view[name] = !state.view[name];
            render();
        });
    });

    // The "total" tile clears all filters
    document.getElementById('tile-total').addEventListener('click', () => {
        const { sortBy, sortDirection } = state.view;
        state.view = { ...defaultView(), sortBy, sortDirection };
        el.typeSelect.value = '';
        el.search.value = '';
        render();
    });

    el.search.addEventListener('input', () => {
        state.view.query = el.search.value;
        render();
    });

    el.typeSelect.addEventListener('change', () => {
        state.view.type = el.typeSelect.value;
        render();
    });

    el.sortBy.addEventListener('change', () => {
        state.view.sortBy = el.sortBy.value;
        render();
    });

    el.sortDirection.addEventListener('click', () => {
        state.view.sortDirection = state.view.sortDirection === 'asc' ? 'desc' : 'asc';
        render();
    });

    el.grid.addEventListener('click', (event) => {
        const button = event.target.closest('[data-details]');
        const aircraft = state.aircraft.find((a) => a.ICAO === button?.dataset.details);
        if (aircraft) openDetails(aircraft);
    });

    // An image that fails to load is replaced by the "no photo" picture
    el.grid.addEventListener(
        'error',
        (event) => {
            const img = event.target;
            if (img.tagName !== 'IMG' || img.hasAttribute('data-fallback')) return;
            img.src = '/static/images/aircraft_not_found.png';
            img.setAttribute('data-fallback', '');
        },
        true,
    );
}

initShell();
initDetails();
initEvents();
render();
refresh();
setInterval(updateCountdown, 1000);

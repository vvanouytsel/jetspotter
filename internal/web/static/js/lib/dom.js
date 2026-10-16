// Small helpers for building HTML safely.
//
// Aircraft data comes from the internet (callsigns, airport names, ...), so it
// must never be pasted into HTML as-is. The html`...` tag escapes every value
// automatically. Use raw() only for HTML that this code made itself.

const ESCAPES = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' };

export function esc(value) {
    return String(value).replace(/[&<>"']/g, (ch) => ESCAPES[ch]);
}

class Raw {
    constructor(value) {
        this.value = value;
    }
}

/** Marks a string as trusted HTML so html`` will not escape it. */
export function raw(value) {
    return new Raw(value);
}

function toHtml(value) {
    if (value === null || value === undefined || value === false) return '';
    if (value instanceof Raw) return value.value;
    if (Array.isArray(value)) return value.map(toHtml).join('');
    return esc(value);
}

/** Template tag: html`<b>${userText}</b>` returns trusted HTML (as a Raw). */
export function html(strings, ...values) {
    let out = strings[0];
    values.forEach((value, i) => {
        out += toHtml(value) + strings[i + 1];
    });
    return new Raw(out);
}

/** Turns the result of html`` into a string that can be assigned to innerHTML. */
export function toString(value) {
    return toHtml(value);
}

/** An icon from the sprite in templates/partials/icons.html. */
export function icon(name) {
    return raw(`<svg class="icon" aria-hidden="true"><use href="#i-${name}"></use></svg>`);
}

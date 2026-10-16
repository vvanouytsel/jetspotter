// Behavior shared by every page: the theme toggle in the header.

export function currentTheme() {
    return document.documentElement.getAttribute('data-theme') === 'dark' ? 'dark' : 'light';
}

function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    try {
        localStorage.setItem('theme', theme);
    } catch (e) {
        // Storage can be blocked; the theme just will not be remembered
    }
    document.dispatchEvent(new CustomEvent('themechange', { detail: theme }));
}

export function initShell() {
    document.getElementById('theme-toggle')?.addEventListener('click', () => {
        setTheme(currentTheme() === 'dark' ? 'light' : 'dark');
    });
}

// Runs in the <head> before the page is drawn: picks the saved theme, or
// the system theme when nothing was saved yet. Kept tiny and free of
// modules on purpose, so there is no flash of the wrong theme.
(function () {
    var theme = null;
    try {
        theme = localStorage.getItem('theme');
    } catch (e) {
        // Storage can be blocked; fall back to the system theme
    }
    if (theme !== 'light' && theme !== 'dark') {
        theme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }
    document.documentElement.setAttribute('data-theme', theme);
})();

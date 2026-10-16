# Web UI

The web UI is plain HTML, CSS and JavaScript. There is no build step and no
framework: everything in `templates/` and `static/` is embedded into the Go
binary, so after changing a file you restart the app to see the change.

## Layout

```text
server.go               Routes, login, and the page renderer
templates/
  layouts/base.html     The frame around every page (<head>, header, footer)
  partials/             header, footer, and the icon sprite
  pages/                One file per page: dashboard, config, login, loading
static/
  css/
    base.css            Colors, spacing, fonts (the "design tokens") and the reset
    components.css      Buttons, chips, badges, header, footer, dialog, loader
    dashboard.css       Summary tiles, toolbar, aircraft cards, details dialog
    config.css          Configuration page
    auth.css            Login and "starting up" pages
  js/
    theme-init.js       Picks light/dark before the page is drawn
    shell.js            Behavior every page shares (theme toggle)
    pages/              One entry point per page
    dashboard/          Dashboard pieces: card, details dialog, filters, demo data
    lib/                dom.js (safe HTML), format.js (numbers, altitude, countries)
  images/
```

## How to change things

- **Colors, spacing, fonts**: edit the variables at the top of `static/css/base.css`.
  The dark theme is the `:root[data-theme="dark"]` block right below the light one.
- **Header or footer**: edit `templates/partials/header.html` or `footer.html`.
  The change shows up on every page.
- **A new icon**: add a `<symbol>` to `templates/partials/icons.html`, then use it with
  `<svg class="icon"><use href="#i-name"></use></svg>`, or `icon('name')` in JavaScript.
- **What an aircraft card shows**: `static/js/dashboard/card.js`.
- **A new page**: add `templates/pages/<name>.html`, add the name to `allPages` in
  `server.go`, and add a route.

## Rules the code follows

- Aircraft data comes from the internet, so never put it into HTML by hand.
  Build HTML with the `html` tag from `lib/dom.js`, which escapes every value.
- CSS lives in the file named for its page or, if two pages use it, in `components.css`.
  No `<style>` blocks inside templates.
- Scripts are ES modules. Each page has one entry file in `static/js/pages/`.

## Trying the dashboard without live data

Open the dashboard with `?test=true` (for example `http://localhost:8080/?test=true`)
to show the fake aircraft from `static/js/dashboard/demo-data.js`.

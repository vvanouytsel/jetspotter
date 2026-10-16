// Shapes of the data the backend sends. These are only documentation for people
// (and editors): nothing here runs. It mirrors `Aircraft` in internal/jetspotter/types.go.

/**
 * @typedef {object} Airport
 * @property {string} [name]
 * @property {string} [municipality]
 * @property {string} [country_name]
 * @property {string} [country_iso_name]  Two letter country code
 * @property {string} [iata_code]
 * @property {string} [icao_code]
 * @property {number} [latitude]
 * @property {number} [longitude]
 */

/**
 * @typedef {object} Airline
 * @property {string} [name]
 * @property {string} [iata]
 * @property {string} [icao]
 * @property {string} [country_iso]
 */

/**
 * @typedef {object} Aircraft
 * @property {string} ICAO               Unique id of the aircraft, also used as the key in the UI
 * @property {string} Callsign
 * @property {string} Description
 * @property {string} Type
 * @property {string} Registration
 * @property {string} Country            English country name
 * @property {number} Altitude           Feet
 * @property {number} Speed              Knots
 * @property {number} Distance           Kilometers from your location
 * @property {number} Heading            Degrees
 * @property {number} BearingFromLocation Degrees, direction from you to the aircraft
 * @property {boolean} Military
 * @property {boolean} Inbound           Flying towards your location
 * @property {boolean} OnGround
 * @property {string} TrackerURL
 * @property {string} ImageURL           Page with more photos
 * @property {string} ImageThumbnailURL
 * @property {string} Photographer
 * @property {Airline} Airline
 * @property {Airport} Origin
 * @property {Airport} Destination
 */

export {};

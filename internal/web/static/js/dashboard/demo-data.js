// Fake aircraft for trying out the dashboard without a connection to the
// ADS-B API. Open the dashboard with ?test=true to use them.

const photo = (id) => ({
    ImageURL: `https://picsum.photos/seed/${id}/600/400`,
    ImageThumbnailURL: `https://picsum.photos/seed/${id}/300/200`,
});

const brussels = { name: 'Brussels Airport', municipality: 'Brussels', country_name: 'Belgium', country_iso_name: 'BE', iata_code: 'BRU', icao_code: 'EBBR', latitude: 50.9, longitude: 4.48 };
const newark = { name: 'Newark Liberty International Airport', municipality: 'Newark', country_name: 'United States', country_iso_name: 'US', iata_code: 'EWR', icao_code: 'KEWR', latitude: 40.69, longitude: -74.17 };
const none = {};

export const demoAircraft = [
    { Callsign: 'ARMY01', Country: 'United States', Military: true, ICAO: 'AE01FF', Altitude: 25000, Speed: 450, Distance: 42, Heading: 225, BearingFromLocation: 170, Type: 'C-17A Globemaster III', Description: 'Military Transport Aircraft', Registration: '02-1111', TrackerURL: '#', Airline: none, Origin: none, Destination: none },
    { Callsign: 'NAVY42', Country: 'United States', Military: true, Inbound: true, ICAO: 'AE1234', Altitude: 18000, Speed: 380, Distance: 25, Heading: 90, BearingFromLocation: 45, Type: 'F/A-18E Super Hornet', Description: 'Military Fighter Aircraft', Registration: '165667', TrackerURL: '#', Airline: none, Origin: none, Destination: none },
    { Callsign: 'UAL123', Country: 'United States', ICAO: 'A12345', Altitude: 35000, Speed: 520, Distance: 15, Heading: 45, BearingFromLocation: 280, Type: 'Boeing 777-300ER', Description: 'Twin-Engine Passenger Jet', Registration: 'N12345', TrackerURL: '#', ...photo('777'), Airline: { name: 'United Airlines', iata: 'UA', icao: 'UAL', country_iso: 'US' }, Origin: brussels, Destination: newark },
    { Callsign: 'DAL456', Country: 'United States', Inbound: true, ICAO: 'A99999', Altitude: 1500, Speed: 160, Distance: 8, Heading: 180, BearingFromLocation: 90, Type: 'Airbus A320', Description: 'Twin-Engine Passenger Jet', Registration: 'N987AA', TrackerURL: '#', ...photo('a320'), Airline: { name: 'Delta Air Lines', iata: 'DL', icao: 'DAL', country_iso: 'US' }, Origin: newark, Destination: brussels },
    { Callsign: 'ABC123', Country: 'Does not exist', ICAO: 'A12346', Altitude: 7500, Speed: 300, Distance: 31, Heading: 10, BearingFromLocation: 300, Type: 'Cessna 172', Description: 'Light Aircraft', Registration: 'ABCD', Airline: none, Origin: none, Destination: none },
    { Callsign: 'AF101', Country: 'United States', Military: true, OnGround: true, ICAO: 'AE2468', Altitude: 0, Speed: 0, Distance: 5, Heading: 90, BearingFromLocation: 120, Type: 'F-16 Fighting Falcon', Description: 'Military Fighter Aircraft', Registration: '86-0241', TrackerURL: '#', ...photo('f16'), Airline: none, Origin: none, Destination: none },
    { Callsign: 'UAL789', Country: 'Belgium', OnGround: true, ICAO: 'A54321', Altitude: 0, Speed: 0, Distance: 3, Heading: 270, BearingFromLocation: 200, Type: 'Boeing 737-800', Description: 'Twin-Engine Passenger Jet', Registration: 'OO-ABC', TrackerURL: '#', ...photo('737'), Airline: none, Origin: none, Destination: none },
];

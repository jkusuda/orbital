// Example asteroids. id / name / hazardous / diameter_km match the backend's
// /positions response; `orbit` is made up for the visualization and goes away
// once positions come from the API.
//
// orbit: a = semi-major axis (Earth radii), e = eccentricity,
//        inclination / node / phase in radians.
export const asteroidData = [
  { id: '2099942', name: 'Apophis', hazardous: true, diameter_km: 0.37, orbit: { a: 2.2, e: 0.19, inclination: 0.06, node: 3.6, phase: 0.4 } },
  { id: '2101955', name: 'Bennu', hazardous: true, diameter_km: 0.49, orbit: { a: 2.8, e: 0.2, inclination: 0.3, node: 0.1, phase: 2.1 } },
  { id: '2065803', name: 'Didymos', hazardous: true, diameter_km: 0.78, orbit: { a: 3.4, e: 0.38, inclination: -0.2, node: 1.3, phase: 4.4 } },
  { id: '2000433', name: 'Eros', hazardous: false, diameter_km: 16.8, orbit: { a: 5.6, e: 0.22, inclination: 0.19, node: 5.3, phase: 1.0 } },
  { id: '2162173', name: 'Ryugu', hazardous: true, diameter_km: 0.9, orbit: { a: 3.0, e: 0.19, inclination: -0.5, node: 4.4, phase: 5.5 } },
  { id: '2025143', name: 'Itokawa', hazardous: false, diameter_km: 0.33, orbit: { a: 3.8, e: 0.28, inclination: 0.7, node: 1.2, phase: 3.0 } },
  { id: '2004179', name: 'Toutatis', hazardous: true, diameter_km: 2.5, orbit: { a: 4.6, e: 0.4, inclination: 0.1, node: 2.2, phase: 0.2 } },
  { id: '2003200', name: 'Phaethon', hazardous: true, diameter_km: 5.8, orbit: { a: 4.2, e: 0.5, inclination: -0.9, node: 4.6, phase: 3.9 } },
  { id: '2001036', name: 'Ganymed', hazardous: false, diameter_km: 37.7, orbit: { a: 6.4, e: 0.3, inclination: 0.46, node: 3.8, phase: 2.6 } },
  { id: '2001620', name: 'Geographos', hazardous: false, diameter_km: 2.5, orbit: { a: 3.6, e: 0.1, inclination: 1.2, node: 5.9, phase: 4.9 } },
  { id: '2003122', name: 'Florence', hazardous: false, diameter_km: 4.9, orbit: { a: 5.0, e: 0.15, inclination: -0.4, node: 0.6, phase: 5.8 } },
  { id: '54339127', name: '2023 DW', hazardous: false, diameter_km: 0.05, orbit: { a: 2.5, e: 0.05, inclination: -0.1, node: 2.8, phase: 1.6 } },
];

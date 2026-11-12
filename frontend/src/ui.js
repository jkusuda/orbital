const $ = (id) => document.getElementById(id);

export function showInfo(asteroid) {
  $('info-name').textContent = asteroid.name;
  $('info-diameter').textContent = formatDiameter(asteroid.diameter_km);
  $('info-status').textContent = asteroid.hazardous ? 'Hazardous' : 'Not hazardous';
  $('info-status').classList.toggle('hazardous', asteroid.hazardous);
  $('info-more').href = asteroid.jpl_url || `https://ssd.jpl.nasa.gov/tools/sbdb_lookup.html#/?sstr=${asteroid.id}`;
  $('info-more').ariaLabel = `More about ${asteroid.name} on NASA JPL's Small-Body Database (opens in a new tab)`;
  $('info').hidden = false;
}

export function onInfoClose(handler) {
  $('info-close').addEventListener('click', handler);
}

// NASA's estimates are rough, so two or three significant figures is honest.
function formatDiameter(km) {
  return km < 1 ? `~${Math.round(km * 1000)} m` : `~${km.toFixed(1)} km`;
}

export function updateDistance(km) {
  $('info-distance').textContent = `${Math.round(km).toLocaleString()} km`;
}

export function hideInfo() {
  $('info').hidden = true;
}

// Loading / error message in the middle of the screen; null hides it.
export function showStatus(message) {
  $('status').textContent = message ?? '';
  $('status').hidden = message == null;
}

// About popup; the dialog's own form closes it.
$('about-open').addEventListener('click', () => $('about').showModal());

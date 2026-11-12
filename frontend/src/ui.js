const $ = (id) => document.getElementById(id);

// Simulation clock state, driven by the pause button and speed slider.
export const clock = { paused: false, speed: 1 };

$('pause').addEventListener('click', () => {
  clock.paused = !clock.paused;
  $('pause').textContent = clock.paused ? '▶' : '❚❚';
  $('pause').setAttribute('aria-label', clock.paused ? 'Play' : 'Pause');
});

$('speed').addEventListener('input', (event) => {
  clock.speed = Number(event.target.value);
  $('speed-label').textContent = `${clock.speed.toFixed(1)}×`;
});

export function showInfo(asteroid) {
  $('info-name').textContent = asteroid.name;
  $('info-diameter').textContent = `${asteroid.diameter_km} km`;
  $('info-status').textContent = asteroid.hazardous ? 'Hazardous' : 'Not hazardous';
  $('info-status').classList.toggle('hazardous', asteroid.hazardous);
  $('info').hidden = false;
}

export function updateDistance(km) {
  $('info-distance').textContent = `${Math.round(km).toLocaleString()} km`;
}

export function hideInfo() {
  $('info').hidden = true;
}

// About popup; the dialog's own form closes it.
$('about-open').addEventListener('click', () => $('about').showModal());

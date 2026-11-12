// The simulated clock behind the time panel.
//   Live:      locked to real time (the default).
//   Otherwise: the viewer scrubbed or is playing a time-lapse, so simulated
//              time sits somewhere within RANGE_MS of now.

const HOUR_MS = 3_600_000;
export const RANGE_MS = 7 * 24 * HOUR_MS; // slider reach either side of now

const SPEEDS = [
  { label: '1×', rate: 1 }, // simulated ms per real ms
  { label: '1 h/s', rate: 3_600 },
  { label: '1 d/s', rate: 86_400 },
];

const $ = (id) => document.getElementById(id);
const scrub = $('scrub');
scrub.min = String(-RANGE_MS / HOUR_MS);
scrub.max = String(RANGE_MS / HOUR_MS);

const state = { live: true, playing: true, speed: 0, date: Date.now() };
let dragging = false;

$('live').addEventListener('click', () => {
  Object.assign(state, { live: true, playing: true, speed: 0 });
});

$('play').addEventListener('click', () => {
  if (state.live) {
    state.live = false; // pausing live freezes the view at this moment
    state.playing = false;
  } else {
    state.playing = !state.playing;
  }
});

$('speed').addEventListener('click', () => {
  state.speed = (state.speed + 1) % SPEEDS.length;
  if (state.live && state.speed !== 0) {
    state.live = false; // speeding up means leaving real time
    state.playing = true;
  }
});

scrub.addEventListener('pointerdown', () => { dragging = true; });
window.addEventListener('pointerup', () => { dragging = false; });
scrub.addEventListener('input', () => {
  state.live = false;
  state.date = Date.now() + Number(scrub.value) * HOUR_MS;
});

// Advances the clock by one frame (deltaMs of real time) and returns the
// simulated time in ms since the Unix epoch.
export function tick(deltaMs) {
  const now = Date.now();
  if (state.live) state.date = now;
  else if (state.playing) state.date += deltaMs * SPEEDS[state.speed].rate;

  // Stop at the slider's ends rather than running past the loaded data.
  const clamped = Math.min(Math.max(state.date, now - RANGE_MS), now + RANGE_MS);
  if (clamped !== state.date) {
    state.date = clamped;
    state.playing = false;
  }

  render(now);
  return state.date;
}

// --- Rendering the panel ---

const dateFormat = new Intl.DateTimeFormat(undefined, {
  weekday: 'short',
  month: 'short',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  second: '2-digit',
});

// "+2d 4h", "−3h 12m", "+0m"
function formatOffset(ms) {
  const sign = ms < 0 ? '−' : '+';
  let minutes = Math.round(Math.abs(ms) / 60_000);
  const days = Math.floor(minutes / 1440);
  minutes -= days * 1440;
  const hours = Math.floor(minutes / 60);
  minutes -= hours * 60;
  const parts = days ? [`${days}d`, `${hours}h`] : hours ? [`${hours}h`, `${minutes}m`] : [`${minutes}m`];
  return sign + parts.join(' ');
}

// Called every frame, so only touch the DOM when something changed.
const shown = {};
function show(id, text) {
  if (shown[id] === text) return;
  shown[id] = text;
  $(id).textContent = text;
}

function render(now) {
  show('date', dateFormat.format(state.date));
  show('offset', state.live ? 'Live' : `${formatOffset(state.date - now)} from now`);
  show('speed', SPEEDS[state.speed].label);
  show('play', state.playing ? '❚❚' : '▶');
  $('play').ariaLabel = state.playing ? 'Pause' : 'Play';
  $('offset').classList.toggle('live', state.live);
  $('live').classList.toggle('active', state.live);
  if (!dragging) scrub.value = String((state.date - now) / HOUR_MS);
}

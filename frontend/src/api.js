async function get(path, params) {
  const res = await fetch(`/api/${path}?${new URLSearchParams(params)}`);
  if (!res.ok) throw new Error(`${path}: ${res.status}`);
  return res.json();
}

export function fetchPositions(time = new Date()) {
  return get('positions', { time: time.toISOString() });
}

export function fetchTrajectories({ start, end, step }) {
  return get('trajectories', { start: start.toISOString(), end: end.toISOString(), step });
}

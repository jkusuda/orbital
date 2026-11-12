import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import { createEarth } from './earth.js';
import { createStarfield } from './background.js';
import { createSun } from './sun.js';
import { createAsteroids } from './asteroids.js';
import { fetchTrajectories } from './api.js';
import { sampleAt, toSceneDirection } from './scale.js';
import { showInfo, updateDistance, hideInfo, showStatus, onInfoClose } from './ui.js';
import { RANGE_MS, tick } from './timeline.js';
import { createCameraRig } from './cameraRig.js';

const HOUR_MS = 3_600_000;
const LOAD_MARGIN_MS = 24 * HOUR_MS; // load a day past the slider's reach on each side
const STEP = '2h'; // backend sample spacing; motion in between is interpolated
const REFRESH_MS = 6 * HOUR_MS; // re-centre the loaded window as real time moves on

// --- Scene setup ---

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 0.1, 30_000); // far enough to reach the sun
camera.position.set(0, 16, 36);

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

const controls = new OrbitControls(camera, renderer.domElement);
controls.enableDamping = true;
controls.minDistance = 2;
controls.maxDistance = 70;
const rig = createCameraRig(camera, controls);

const sun = createSun(new THREE.Vector3(1, 0, 0)); // aimed for real on the first frame
scene.add(createStarfield(), sun, new THREE.AmbientLight(0xffffff, 0.08));

const earth = await createEarth();
scene.add(earth);

// --- Hover and selection ---

const raycaster = new THREE.Raycaster();
const pointer = new THREE.Vector2();
let asteroids = [];
let rocks = [];
let hovered = null;
let selected = null;

function asteroidAt(event) {
  pointer.x = (event.clientX / window.innerWidth) * 2 - 1;
  pointer.y = -(event.clientY / window.innerHeight) * 2 + 1;
  raycaster.setFromCamera(pointer, camera);
  const hit = raycaster.intersectObjects(rocks)[0];
  return hit ? asteroids[rocks.indexOf(hit.object)] : null;
}

function hover(asteroid) {
  hovered?.setHovered(false);
  hovered = asteroid;
  hovered?.setHovered(true);
  renderer.domElement.style.cursor = hovered ? 'pointer' : '';
}

// Selecting flies the camera to the asteroid; deselecting flies it back.
function select(asteroid) {
  if (asteroid === selected) return;
  selected?.setSelected(false);
  selected = asteroid;
  selected?.setSelected(true);
  if (selected) {
    showInfo(selected.data);
    rig.focusOn(selected);
  } else {
    hideInfo();
    rig.release();
  }
}

onInfoClose(() => select(null));

renderer.domElement.addEventListener('pointermove', (event) => {
  const asteroid = asteroidAt(event);
  if (asteroid !== hovered) hover(asteroid);
});

// Only treat it as a click if the pointer didn't move (i.e. not a camera drag).
let pointerDown = null;
renderer.domElement.addEventListener('pointerdown', (event) => {
  pointerDown = { x: event.clientX, y: event.clientY };
});
renderer.domElement.addEventListener('pointerup', (event) => {
  if (!pointerDown) return;
  const moved = Math.hypot(event.clientX - pointerDown.x, event.clientY - pointerDown.y);
  if (moved < 5) select(asteroidAt(event));
  pointerDown = null;
});

// --- Data ---

let data = null; // { firstTime, stepMs, count, sun } for the loaded window

// Fetches trajectories around the real now and swaps them into the scene,
// keeping the current selection if that asteroid is still there.
async function loadData() {
  const now = Date.now();
  const reach = RANGE_MS + LOAD_MARGIN_MS;
  const res = await fetchTrajectories({ start: new Date(now - reach), end: new Date(now + reach), step: STEP });
  const next = await createAsteroids(res.asteroids);

  hover(null);
  asteroids.forEach((asteroid) => {
    scene.remove(asteroid.group);
    asteroid.dispose();
  });

  asteroids = next;
  rocks = asteroids.map((asteroid) => asteroid.rock);
  asteroids.forEach((asteroid) => scene.add(asteroid.group));

  // Keep the selection (and the camera on it) if that asteroid is still here.
  const same = selected && asteroids.find((asteroid) => asteroid.data.id === selected.data.id);
  if (same) {
    selected = same;
    selected.setSelected(true);
    rig.swap(same);
  } else {
    select(null);
  }

  const times = res.times.map((t) => Date.parse(t));
  data = {
    firstTime: times[0],
    stepMs: times.length > 1 ? times[1] - times[0] : 1,
    count: times.length,
    sun: res.sun,
  };
}

try {
  await loadData();
  showStatus(asteroids.length ? null : 'No asteroids loaded.');
} catch (err) {
  showStatus(`Couldn't load asteroids (${err.message}). Is the backend running on :8080?`);
  throw err;
}

setInterval(() => {
  loadData().catch((err) => console.error('refresh failed, keeping previous data:', err));
}, REFRESH_MS);

// --- Animation loop ---

let lastFrame = performance.now();
const sunAU = [0, 0, 0];
const sunDir = new THREE.Vector3();

renderer.setAnimationLoop((frame) => {
  const date = tick(frame - lastFrame);
  lastFrame = frame;

  // Fractional index into the trajectory samples for the simulated time.
  const sample = THREE.MathUtils.clamp((date - data.firstTime) / data.stepMs, 0, data.count - 1);

  earth.userData.update(date);
  sun.userData.setDirection(toSceneDirection(sampleAt(data.sun, sample, sunAU), sunDir));
  asteroids.forEach((asteroid) => asteroid.update(sample, frame / 1000));
  if (selected) updateDistance(selected.distanceKm);

  rig.update();
  controls.update();
  renderer.render(scene, camera);
});

window.addEventListener('resize', () => {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
});

import * as THREE from 'three';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';
import { createEarth } from './earth.js';
import { createStarfield } from './background.js';
import { createSun } from './sun.js';
import { createAsteroids } from './asteroids.js';
import { asteroidData } from './data.js';
import { clock, showInfo, updateDistance, hideInfo } from './ui.js';

const EARTH_RADIUS_KM = 6371; // one scene unit

// --- Scene setup ---

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(50, window.innerWidth / window.innerHeight, 0.1, 30_000); // far enough to reach the sun
camera.position.set(0, 5, 12);

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

const controls = new OrbitControls(camera, renderer.domElement);
controls.enableDamping = true;
controls.minDistance = 2;
controls.maxDistance = 40;

scene.add(
  createStarfield(),
  createSun(new THREE.Vector3(10, 4, 6)),
  new THREE.AmbientLight(0xffffff, 0.08),
);

const earth = await createEarth();
scene.add(earth);

const asteroids = await createAsteroids(asteroidData);
asteroids.forEach((asteroid) => scene.add(asteroid.group));

// --- Hover and selection ---

const raycaster = new THREE.Raycaster();
const pointer = new THREE.Vector2();
const rocks = asteroids.map((asteroid) => asteroid.rock);
let hovered = null;
let selected = null;

function asteroidAt(event) {
  pointer.x = (event.clientX / window.innerWidth) * 2 - 1;
  pointer.y = -(event.clientY / window.innerHeight) * 2 + 1;
  raycaster.setFromCamera(pointer, camera);
  const hit = raycaster.intersectObjects(rocks)[0];
  return hit ? asteroids[rocks.indexOf(hit.object)] : null;
}

function select(asteroid) {
  selected?.setSelected(false);
  selected = asteroid;
  selected?.setSelected(true);
  if (selected) showInfo(selected.data);
  else hideInfo();
}

renderer.domElement.addEventListener('pointermove', (event) => {
  const asteroid = asteroidAt(event);
  if (asteroid === hovered) return;
  hovered?.setHovered(false);
  hovered = asteroid;
  hovered?.setHovered(true);
  renderer.domElement.style.cursor = hovered ? 'pointer' : '';
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

// --- Animation loop ---

let simTime = 0;
let lastFrame = performance.now();

renderer.setAnimationLoop((now) => {
  const delta = (now - lastFrame) / 1000;
  lastFrame = now;
  if (!clock.paused) simTime += delta * clock.speed;

  earth.userData.update(simTime);
  asteroids.forEach((asteroid) => asteroid.update(simTime));
  if (selected) {
    const position = selected.rock.getWorldPosition(new THREE.Vector3());
    updateDistance(position.length() * EARTH_RADIUS_KM);
  }

  controls.update();
  renderer.render(scene, camera);
});

window.addEventListener('resize', () => {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
});

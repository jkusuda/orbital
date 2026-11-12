import * as THREE from 'three';
import { OBJLoader } from 'three/addons/loaders/OBJLoader.js';

const HAZARDOUS_COLOR = 0xff5a47;
const SAFE_COLOR = 0x6ea8ff;
const ORBIT_SEGMENTS = 128;
const ORBIT_OPACITY = 0.5;

// Rock models from https://github.com/bobbyroe/solar-system (MIT).
async function loadRockGeometries() {
  const loader = new OBJLoader();
  const rocks = await Promise.all(
    ['Rock1', 'Rock2', 'Rock3'].map((name) => loader.loadAsync(`/rocks/${name}.obj`)),
  );
  return rocks.map((rock) => rock.getObjectByProperty('isMesh', true).geometry);
}

// Point on an ellipse lying in the XZ plane with Earth at one focus.
function orbitPoint({ a, e }, angle) {
  const r = (a * (1 - e * e)) / (1 + e * Math.cos(angle));
  return new THREE.Vector3(Math.cos(angle) * r, 0, Math.sin(angle) * r);
}

function createOrbitLine(orbit, color) {
  const points = [];
  for (let i = 0; i < ORBIT_SEGMENTS; i += 1) {
    points.push(orbitPoint(orbit, (i / ORBIT_SEGMENTS) * Math.PI * 2));
  }
  const geometry = new THREE.BufferGeometry().setFromPoints(points);
  const material = new THREE.LineBasicMaterial({ color, transparent: true, opacity: ORBIT_OPACITY });
  return new THREE.LineLoop(geometry, material);
}

function createAsteroid(data, geometry) {
  const { orbit } = data;
  const color = data.hazardous ? HAZARDOUS_COLOR : SAFE_COLOR;

  const rock = new THREE.Mesh(
    geometry,
    new THREE.MeshStandardMaterial({ color: 0x9a938a, flatShading: true }),
  );
  // Real sizes would be invisible at this scale, so exaggerate and clamp.
  rock.scale.setScalar(THREE.MathUtils.clamp(0.06 + data.diameter_km * 0.01, 0.06, 0.22));
  rock.rotation.setFromVector3(new THREE.Vector3().randomDirection().multiplyScalar(Math.PI));
  rock.userData.asteroid = data;

  const line = createOrbitLine(orbit, color);

  // The group tilts the flat orbit into place; rock and line share that frame.
  const group = new THREE.Group();
  group.rotation.set(orbit.inclination, orbit.node, 0, 'YXZ');
  group.add(rock, line);

  // Not real Kepler motion (the angle advances uniformly), but wider orbits are slower.
  const angularSpeed = 0.5 / orbit.a ** 1.5;
  const tumbleSpeed = 0.2 + Math.random() * 0.4;

  return {
    data,
    group,
    rock,
    update(t) {
      rock.position.copy(orbitPoint(orbit, orbit.phase + t * angularSpeed));
      rock.rotation.y = t * tumbleSpeed;
    },
    setHovered(hovered) {
      rock.material.emissive.set(hovered ? color : 0x000000);
    },
    setSelected(selected) {
      line.material.opacity = selected ? 1 : ORBIT_OPACITY;
    },
  };
}

export async function createAsteroids(asteroidData) {
  const geometries = await loadRockGeometries();
  return asteroidData.map((data, i) => createAsteroid(data, geometries[i % geometries.length]));
}

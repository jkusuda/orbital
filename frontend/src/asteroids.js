import * as THREE from 'three';
import { OBJLoader } from 'three/addons/loaders/OBJLoader.js';
import { KM_PER_AU, sampleAt, toScene } from './scale.js';

const HAZARDOUS_COLOR = 0xff5a47;
const SAFE_COLOR = 0x6ea8ff;
const PATH_OPACITY = 0.35;

// Rock models from https://github.com/bobbyroe/solar-system (MIT).
// Loaded once and shared by every asteroid, across data refreshes too.
let rockGeometries = null;
function loadRockGeometries() {
  rockGeometries ??= (async () => {
    const loader = new OBJLoader();
    const rocks = await Promise.all(
      ['Rock1', 'Rock2', 'Rock3'].map((name) => loader.loadAsync(`/rocks/${name}.obj`)),
    );
    return rocks.map((rock) => rock.getObjectByProperty('isMesh', true).geometry);
  })();
  return rockGeometries;
}

// The real path relative to Earth over the loaded time window, scaled like the rock.
function createPathLine(points, color) {
  const geometry = new THREE.BufferGeometry().setFromPoints(points.map((p) => toScene(p)));
  const material = new THREE.LineBasicMaterial({ color, transparent: true, opacity: PATH_OPACITY });
  return new THREE.Line(geometry, material);
}

// `data` is one entry of the /api/trajectories response.
function createAsteroid(data, geometry) {
  const color = data.hazardous ? HAZARDOUS_COLOR : SAFE_COLOR;

  const rock = new THREE.Mesh(
    geometry,
    new THREE.MeshStandardMaterial({ color: 0x9a938a, flatShading: true }),
  );
  // Real sizes would be invisible at this scale, so exaggerate and clamp.
  rock.scale.setScalar(THREE.MathUtils.clamp(0.2 + data.diameter_km * 0.03, 0.2, 0.6));
  rock.rotation.setFromVector3(new THREE.Vector3().randomDirection().multiplyScalar(Math.PI));
  rock.userData.asteroid = data;

  const line = createPathLine(data.points, color);

  const group = new THREE.Group();
  group.add(rock, line);

  const tumbleSpeed = 0.2 + Math.random() * 0.4;
  const au = [0, 0, 0];

  return {
    data,
    group,
    rock,
    distanceKm: 0, // real distance from Earth, updated each frame
    // sample: fractional index into data.points; t: seconds, for the tumble.
    update(sample, t) {
      const p = sampleAt(data.points, sample, au);
      toScene(p, rock.position);
      this.distanceKm = Math.hypot(p[0], p[1], p[2]) * KM_PER_AU;
      rock.rotation.y = t * tumbleSpeed;
    },
    setHovered(hovered) {
      rock.material.emissive.set(hovered ? color : 0x000000);
    },
    setSelected(selected) {
      line.material.opacity = selected ? 1 : PATH_OPACITY;
    },
    // Frees GPU memory when the asteroid is replaced. The rock geometry is shared, so it stays.
    dispose() {
      line.geometry.dispose();
      line.material.dispose();
      rock.material.dispose();
    },
  };
}

export async function createAsteroids(trajectories) {
  const geometries = await loadRockGeometries();
  return trajectories.map((data, i) => createAsteroid(data, geometries[i % geometries.length]));
}

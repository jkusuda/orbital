// Starfield adapted from https://github.com/bobbyroe/solar-system (MIT).
import * as THREE from 'three';

const textureLoader = new THREE.TextureLoader();

export function createStarfield({ numStars = 2000, size = 0.6 } = {}) {
  const positions = [];
  const colors = [];
  const color = new THREE.Color();

  for (let i = 0; i < numStars; i += 1) {
    const star = new THREE.Vector3().randomDirection().multiplyScalar(80 + Math.random() * 80);
    positions.push(star.x, star.y, star.z);
    color.setHSL(0.6, 0.2, Math.random());
    colors.push(color.r, color.g, color.b);
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position', new THREE.Float32BufferAttribute(positions, 3));
  geometry.setAttribute('color', new THREE.Float32BufferAttribute(colors, 3));

  const material = new THREE.PointsMaterial({
    size,
    vertexColors: true,
    map: textureLoader.load('/textures/circle.png'),
    transparent: true,
    depthWrite: false,
  });
  return new THREE.Points(geometry, material);
}

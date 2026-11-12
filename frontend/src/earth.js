import * as THREE from 'three';
import { createFresnelMat } from './fresnelMat.js';

const AXIAL_TILT = THREE.MathUtils.degToRad(23.4);
const SPIN_RATE = 0.05; // radians per simulated second
const DETAIL = 6; // icosahedron subdivisions; lower = chunkier
const EDGE_INSET = 6; // px to skip at the map's borders, which are white

// Reads the equirectangular map's pixels so faces can be coloured from it.
async function loadMapPixels(url) {
  const image = await new THREE.ImageLoader().loadAsync(url);
  const canvas = document.createElement('canvas');
  canvas.width = image.width;
  canvas.height = image.height;
  const context = canvas.getContext('2d');
  context.drawImage(image, 0, 0);
  return context.getImageData(0, 0, image.width, image.height);
}

// Colour of the map at the point where `direction` leaves the globe.
function sampleMap(pixels, direction) {
  const u = Math.atan2(direction.z, -direction.x) / (Math.PI * 2) + 0.5;
  const v = Math.acos(direction.y) / Math.PI;
  const maxX = pixels.width - 1 - EDGE_INSET;
  const maxY = pixels.height - 1 - EDGE_INSET;
  const x = Math.round(THREE.MathUtils.clamp(u * pixels.width, EDGE_INSET, maxX));
  const y = Math.round(THREE.MathUtils.clamp(v * pixels.height, EDGE_INSET, maxY));
  const i = (y * pixels.width + x) * 4;
  return new THREE.Color().setRGB(
    pixels.data[i] / 255,
    pixels.data[i + 1] / 255,
    pixels.data[i + 2] / 255,
    THREE.SRGBColorSpace,
  );
}

// Gives each triangle one flat colour, sampled from the map at its centre.
function colorFaces(geometry, pixels) {
  const position = geometry.attributes.position;
  const colors = new Float32Array(position.count * 3);
  const a = new THREE.Vector3();
  const b = new THREE.Vector3();
  const c = new THREE.Vector3();

  for (let i = 0; i < position.count; i += 3) {
    a.fromBufferAttribute(position, i);
    b.fromBufferAttribute(position, i + 1);
    c.fromBufferAttribute(position, i + 2);
    const color = sampleMap(pixels, a.add(b).add(c).normalize());
    for (let j = 0; j < 3; j += 1) color.toArray(colors, (i + j) * 3);
  }
  geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));
}

// Earth has radius 1; everything else in the scene is sized relative to it.
export async function createEarth() {
  const geometry = new THREE.IcosahedronGeometry(1, DETAIL);
  colorFaces(geometry, await loadMapPixels('/textures/earth.png'));

  const surface = new THREE.Mesh(
    geometry,
    new THREE.MeshStandardMaterial({ vertexColors: true, flatShading: true }),
  );

  const atmosphere = new THREE.Mesh(new THREE.IcosahedronGeometry(1.02, 12), createFresnelMat());

  const earth = new THREE.Group();
  earth.rotation.z = -AXIAL_TILT;
  earth.add(surface, atmosphere);

  earth.userData.update = (t) => {
    surface.rotation.y = t * SPIN_RATE;
  };
  return earth;
}

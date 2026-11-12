import * as THREE from 'three';

// True proportions, in Earth radii (the scene unit).
const SUN_RADIUS = 109.2;
const SUN_DISTANCE = 23_455; // 1 AU
const SUN_COLOR = 0xffe2a8;
const GLOW_SCALE = 6; // glow sprite width, in sun diameters

// The sun at real size and distance, plus the light it casts on the scene.
export function createSun(direction) {
  const disc = new THREE.Mesh(
    new THREE.IcosahedronGeometry(SUN_RADIUS, 2),
    new THREE.MeshBasicMaterial({ color: SUN_COLOR }),
  );

  const glow = new THREE.Sprite(
    new THREE.SpriteMaterial({
      color: SUN_COLOR,
      map: new THREE.TextureLoader().load('/textures/rad-grad.png'),
      blending: THREE.AdditiveBlending,
      depthWrite: false,
    }),
  );
  glow.scale.setScalar(SUN_RADIUS * 2 * GLOW_SCALE);

  const light = new THREE.DirectionalLight(0xfff4e0, 3);

  const sun = new THREE.Group();
  sun.position.copy(direction).normalize().multiplyScalar(SUN_DISTANCE);
  sun.add(disc, glow, light);
  return sun;
}

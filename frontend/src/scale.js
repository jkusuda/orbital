import * as THREE from 'three';

export const KM_PER_AU = 149_597_870.7;

const R0_AU = 0.001;
const K = 2.5;

export function toScene([x, y, z], target = new THREE.Vector3()) {
  target.set(x, z, -y);
  const r = target.length();
  if (r === 0) return target;
  return target.multiplyScalar((1 + K * Math.log1p(r / R0_AU)) / r);
}

export function toSceneDirection([x, y, z], target = new THREE.Vector3()) {
  return target.set(x, z, -y).normalize();
}

export function sampleAt(points, s, out = [0, 0, 0]) {
  if (points.length === 1) return points[0];
  const i = Math.min(Math.max(Math.floor(s), 0), points.length - 2);
  const f = s - i;
  const a = points[i];
  const b = points[i + 1];
  out[0] = a[0] + (b[0] - a[0]) * f;
  out[1] = a[1] + (b[1] - a[1]) * f;
  out[2] = a[2] + (b[2] - a[2]) * f;
  return out;
}

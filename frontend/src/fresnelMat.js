// Rim-glow shader from https://github.com/bobbyroe/solar-system (MIT).
import * as THREE from 'three';

const vertexShader = /* glsl */ `
  uniform float fresnelBias;
  uniform float fresnelScale;
  uniform float fresnelPower;
  varying float vReflectionFactor;

  void main() {
    vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
    vec4 worldPosition = modelMatrix * vec4(position, 1.0);
    vec3 worldNormal = normalize(mat3(modelMatrix[0].xyz, modelMatrix[1].xyz, modelMatrix[2].xyz) * normal);
    vec3 I = worldPosition.xyz - cameraPosition;
    vReflectionFactor = fresnelBias + fresnelScale * pow(1.0 + dot(normalize(I), worldNormal), fresnelPower);
    gl_Position = projectionMatrix * mvPosition;
  }
`;

const fragmentShader = /* glsl */ `
  uniform vec3 rimColor;
  uniform vec3 facingColor;
  varying float vReflectionFactor;

  void main() {
    float f = clamp(vReflectionFactor, 0.0, 1.0);
    gl_FragColor = vec4(mix(facingColor, rimColor, vec3(f)), f);
  }
`;

export function createFresnelMat({ rimHex = 0x0088ff, facingHex = 0x000000 } = {}) {
  return new THREE.ShaderMaterial({
    uniforms: {
      rimColor: { value: new THREE.Color(rimHex) },
      facingColor: { value: new THREE.Color(facingHex) },
      fresnelBias: { value: 0.1 },
      fresnelScale: { value: 1.0 },
      fresnelPower: { value: 4.0 },
    },
    vertexShader,
    fragmentShader,
    transparent: true,
    blending: THREE.AdditiveBlending,
  });
}

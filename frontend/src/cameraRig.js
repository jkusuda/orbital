import * as THREE from 'three';

const FLY_MS = 1200;
const BEHIND = 3.5; // how far past the asteroid (away from Earth) the camera sits
const LIFT = 1.0; // and above it, so the rock doesn't hide Earth
const SIDE = 1.0; // and to one side, so the rock sits off-centre
const UP = new THREE.Vector3(0, 1, 0);
const EARTH = new THREE.Vector3(0, 0, 0);

const ease = (t) => (t < 0.5 ? 4 * t * t * t : 1 - (-2 * t + 2) ** 3 / 2);

// Moves the camera between a free "home" view and a view from just behind a
// selected asteroid, looking back at Earth. While focused it follows the
// asteroid as simulated time moves it; OrbitControls still work throughout.
export function createCameraRig(camera, controls) {
  const home = { position: new THREE.Vector3(), target: new THREE.Vector3() };
  let focus = null; // asteroid being followed
  let flight = null; // { fromPosition, fromTarget, start } while animating
  let resync = false; // focus was swapped in place: re-read its position, don't jump
  const lastRock = new THREE.Vector3();
  const goal = new THREE.Vector3();
  const offset = new THREE.Vector3();
  const fromDir = new THREE.Vector3();
  const toDir = new THREE.Vector3();
  const turn = new THREE.Quaternion();

  function viewpoint(rockPosition, out) {
    const outward = rockPosition.clone().normalize();
    const side = new THREE.Vector3().crossVectors(outward, UP);
    if (side.lengthSq() < 1e-6) side.set(1, 0, 0); // asteroid straight above/below Earth
    side.normalize();
    return out
      .copy(rockPosition)
      .addScaledVector(outward, BEHIND)
      .addScaledVector(UP, LIFT)
      .addScaledVector(side, SIDE);
  }

  function startFlight() {
    flight = { fromPosition: camera.position.clone(), fromTarget: controls.target.clone(), start: performance.now() };
    controls.enabled = false;
  }

  return {
    // Fly to `asteroid` and keep following it.
    focusOn(asteroid) {
      if (!focus && !flight) {
        home.position.copy(camera.position); // remember the free view to return to
        home.target.copy(controls.target);
      }
      focus = asteroid;
      startFlight();
    },

    // Swap the followed asteroid for an equivalent one (after a data refresh)
    // without flying anywhere.
    swap(asteroid) {
      focus = asteroid;
      resync = true;
    },

    // Fly back to the free view.
    release() {
      if (!focus) return;
      focus = null;
      startFlight();
    },

    // Call once per frame, after asteroid positions update and before controls.update().
    update() {
      if (flight) {
        const t = ease(Math.min((performance.now() - flight.start) / FLY_MS, 1));
        if (focus) {
          viewpoint(focus.rock.position, goal); // re-aimed each frame: the asteroid may be moving
          camera.position.lerpVectors(flight.fromPosition, goal, t);
          controls.target.lerpVectors(flight.fromTarget, EARTH, t);
        } else {
          camera.position.lerpVectors(flight.fromPosition, home.position, t);
          controls.target.lerpVectors(flight.fromTarget, home.target, t);
        }
        if (t === 1) {
          flight = null;
          controls.enabled = true;
        }
        if (focus) lastRock.copy(focus.rock.position);
        return;
      }

      if (!focus) return;
      const rock = focus.rock.position;
      if (resync) {
        resync = false;
      } else if (!rock.equals(lastRock)) {
        // Follow: carry the camera's offset from the asteroid along with it, turned by
        // however much the asteroid's direction from Earth turned. That keeps the
        // asteroid in the same spot on screen while the camera keeps looking at Earth.
        turn.setFromUnitVectors(fromDir.copy(lastRock).normalize(), toDir.copy(rock).normalize());
        offset.subVectors(camera.position, lastRock).applyQuaternion(turn);
        camera.position.copy(rock).add(offset);
      }
      lastRock.copy(rock);
    },
  };
}

#pragma once

struct Elements {
    double epoch;          // JD
    double eccentricity;
    double semiMajorAxis;  // AU
    double inclination;    // deg
    double ascendingNode;  // deg
    double perihelionArg;  // deg
    double meanAnomaly;    // deg at epoch
};

struct Vec3 {
    double x, y, z;
};

Vec3 position(const Elements& el, double jd);

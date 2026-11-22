#include "propagate.h"

#include <cmath>
#include <utility>

namespace {

constexpr double kPi = 3.14159265358979323846;
constexpr double kGauss = 0.01720209895;

double degToRad(double d) { return d * kPi / 180.0; }

double meanAnomaly(const Elements& el, double jd) {
    double n = kGauss / std::pow(el.semiMajorAxis, 1.5);
    double M = degToRad(el.meanAnomaly) + n * (jd - el.epoch);
    M = std::fmod(M, 2 * kPi);
    if (M < 0)
        M += 2 * kPi;
    return M;
}

double solveKepler(double M, double e) {
    double E = (e > 0.8) ? kPi : M;
    for (int i = 0; i < 50; ++i) {
        double next = E - (E - e * std::sin(E) - M) / (1 - e * std::cos(E));
        if (std::abs(next - E) < 1e-12) return next;
        E = next;
    }
    return E;
}

std::pair<double, double> orbitalPlane(double a, double e, double E) {
    return {a * (std::cos(E) - e), a * std::sqrt(1 - e * e) * std::sin(E)};
}

Vec3 rotate(double x, double y, const Elements& el) {
    double i = degToRad(el.inclination);
    double O = degToRad(el.ascendingNode);
    double w = degToRad(el.perihelionArg);

    double cosO = std::cos(O), sinO = std::sin(O);
    double cosw = std::cos(w), sinw = std::sin(w);
    double cosi = std::cos(i), sini = std::sin(i);

    return {
        (cosw * cosO - sinw * sinO * cosi) * x + (-sinw * cosO - cosw * sinO * cosi) * y,
        (cosw * sinO + sinw * cosO * cosi) * x + (-sinw * sinO + cosw * cosO * cosi) * y,
        (sinw * sini) * x + (cosw * sini) * y,
    };
}

}  // namespace

Vec3 position(const Elements& el, double jd) {
    double M = meanAnomaly(el, jd);
    double E = solveKepler(M, el.eccentricity);
    auto [x, y] = orbitalPlane(el.semiMajorAxis, el.eccentricity, E);
    return rotate(x, y, el);
}

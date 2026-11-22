#include <cstdio>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>
#include <grpcpp/ext/proto_server_reflection_plugin.h>

#include "orbital/v1/engine.grpc.pb.h"
#include "propagate.h"

namespace ov1 = orbital::v1;

class PropagatorImpl final : public ov1::PropagatorService::Service {
    grpc::Status Propagate(grpc::ServerContext*, const ov1::PropagateRequest* req,
                           ov1::PropagateResponse* resp) override {
        for (const auto& a : req->asteroids()) {
            Elements el{a.epoch_jd(),        a.eccentricity(),       a.semi_major_axis_au(),
                        a.inclination_deg(), a.ascending_node_deg(), a.perihelion_arg_deg(),
                        a.mean_anomaly_deg()};

            ov1::Trajectory* traj = resp->add_trajectories();
            for (double jd : req->times_jd()) {
                Vec3 p = position(el, jd);
                ov1::Vec3* out = traj->add_points();
                out->set_x(p.x);
                out->set_y(p.y);
                out->set_z(p.z);
            }
        }
        return grpc::Status::OK;
    }
};

int main() {
    const std::string addr = "0.0.0.0:50051";
    PropagatorImpl service;

    grpc::reflection::InitProtoReflectionServerBuilderPlugin();
    grpc::ServerBuilder builder;
    builder.AddListeningPort(addr, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);

    std::unique_ptr<grpc::Server> server = builder.BuildAndStart();
    if (!server) {
        std::fprintf(stderr, "failed to listen on %s\n", addr.c_str());
        return 1;
    }
    std::printf("engine listening on %s\n", addr.c_str());
    std::fflush(stdout);
    server->Wait();
}

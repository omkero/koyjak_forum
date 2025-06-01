import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
  thresholds: {
    // Assert that 99% of requests finish within 3000ms.
    http_req_duration: ["p(99) < 3000"],
  },
  // Ramp the number of virtual users up and down
  stages: [
    { duration: '1s', target: 1 },
  ],
};

// Simulated user behavior
export default function () {
  let res = http.get("http://127.0.0.1:8080/");
  // Validate response status
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}
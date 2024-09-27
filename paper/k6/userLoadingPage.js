import http from "k6/http";
import { sleep } from "k6";

export const options = {
    vus: 50, // vu stands for virtual user
    duration: '30s',
};

export default function() {
    http.get("http://localhost:8080");
    sleep(0.5);
}

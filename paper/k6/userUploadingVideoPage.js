import http from "k6/http";
import exec from "k6/execution";
import { sleep } from "k6";

export const options = {
    vus: 10, // vu stands for virtual user
    duration: '10s',
};

const video = open("/home/ayaka/Downloads/furina-upside-down-underwater-genshin-impact-moewalls-com.mp4", "b");

export default function() {
    const data = {
        title: `video-${exec.vu.idInTest}`,
        creator: `${exec.vu.idInTest}`,
        media: http.file(video, 'video.mp4'),
    };

    http.post("http://localhost:8080/api/upload", data);
    sleep(0.5);
}

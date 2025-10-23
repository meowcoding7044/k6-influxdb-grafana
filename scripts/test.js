import http from 'k6/http';

export let options = {
  vus: 5,          // virtual users
  duration: '5s', // test duration
};

export default function () {
  http.get('http://host.docker.internal:8000/hello');
}
//docker compose run --rm k6 run /scripts/test.js
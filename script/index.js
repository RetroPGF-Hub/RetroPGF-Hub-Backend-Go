import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 50,
  duration: '30s',
  ext: {
    loadimpact: {
      name: 'Test Echo (20/01/2024)'
    }
  }
};

export default function() {
  http.get('http://localhost:5001/project_v1/projects?limit=80&skip=0');
  sleep(1);
}
import http from 'k6/http'
import { check } from 'k6'

export default function load() {
    let res = http.get('http://localhost:22313')

    check(res, { 'is status 200': (r) => r.status === 200 })
}

load();
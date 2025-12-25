import { updateWorkerList, updateStatus } from './ui.js';
import { initChart, addDataToChart } from './charts.js';

const socket = new WebSocket(`ws://${window.location.host}/ws`);

socket.onopen = () => updateStatus(true);
socket.onclose = () => updateStatus(false);

socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    
    switch(msg.type) {
        case "WORKER_JOIN":
            updateWorkerList(msg.data);
            break;
        case "RESULT":
            addDataToChart(msg.data);
            updateWorkerMetrics(msg.data.worker_id, msg.data.metadata); // Mise à jour IHM temps réel
            break;
    }
};

document.getElementById('btn-fibo').onclick = async () => {
    await fetch('/run', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            handler: "fibonacci",
            params: { series: 5, limit: 400000 }
        })
    });
};

initChart();
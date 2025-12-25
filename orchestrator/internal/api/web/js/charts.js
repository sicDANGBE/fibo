let chart;
const workerColors = {
    'rust': '#f97316',
    'go': '#00add8',
    'node': '#84cc16',
    'python': '#3b82f6'
};

export function initChart() {
    const ctx = document.getElementById('mainChart').getContext('2d');
    chart = new Chart(ctx, {
        type: 'line',
        data: { datasets: [] },
        options: {
            animation: false,
            scales: {
                x: { type: 'linear', grid: { color: '#1e293b' } },
                y: { grid: { color: '#1e293b' } }
            },
            plugins: { legend: { position: 'bottom' } }
        }
    });
}

export function addDataToChart(res) {
    let dataset = chart.data.datasets.find(d => d.id === res.worker_id);
    
    if (!dataset) {
        dataset = {
            id: res.worker_id,
            label: `${res.handler} - ${res.worker_id.substring(0,6)}`,
            data: [],
            borderColor: workerColors[res.handler] || '#fff',
            borderWidth: 2,
            pointRadius: 0
        };
        chart.data.datasets.push(dataset);
    }

    dataset.data.push({ x: res.index, y: res.timestamp % 10000 }); // Exemple de métrique
    if (dataset.data.length > 100) dataset.data.shift();
    chart.update('none');
}
export function updateStatus(connected) {
    const indicator = document.getElementById('ws-status');
    if (connected) {
        indicator.classList.replace('bg-red-500', 'bg-emerald-500');
        indicator.classList.replace('shadow-[0_0_10px_red]', 'shadow-[0_0_10px_#10b981]');
    } else {
        indicator.classList.replace('bg-emerald-500', 'bg-red-500');
        indicator.classList.replace('shadow-[0_0_10px_#10b981]', 'shadow-[0_0_10px_red]');
    }
}

export function updateWorkerMetrics(workerId, metadata) {
    const el = document.getElementById(`worker-${workerId}`);
    if (!el) return;

    const statsContainer = el.querySelector('.worker-stats') || document.createElement('div');
    if (!el.querySelector('.worker-stats')) {
        statsContainer.className = "worker-stats mt-3 pt-2 border-t border-slate-600/30 grid grid-cols-2 gap-2";
        el.appendChild(statsContainer);
    }

    statsContainer.innerHTML = `
        <div class="flex flex-col">
            <span class="text-[9px] text-slate-500 uppercase">Load</span>
            <span class="text-xs font-mono text-blue-400">${metadata.cpu} gor</span>
        </div>
        <div class="flex flex-col text-right">
            <span class="text-[9px] text-slate-500 uppercase">Mem</span>
            <span class="text-xs font-mono text-emerald-400">${metadata.ram}MB</span>
        </div>
        <div class="col-span-2 text-center text-[10px] text-slate-400 font-mono mt-1">
            Data Stream: ${(metadata.net / 1024).toFixed(2)} KB
        </div>
    `;
}

export function updateWorkerList(worker) {
    const list = document.getElementById('worker-list');
    const id = `worker-${worker.id}`;
    if (document.getElementById(id)) return;

    const el = document.createElement('div');
    el.id = id;
    // Ajout des classes pour l'expansion et l'interactivité 
    el.className = "p-4 bg-slate-800/80 rounded-2xl border border-slate-700 cursor-pointer transition-all duration-300 hover:border-blue-500 overflow-hidden mb-4";
    
    el.onclick = () => {
        // Toggle pour agrandir la tuile sur tout le conteneur si besoin
        el.classList.toggle('ring-2');
        el.classList.toggle('ring-blue-500');
        el.querySelector('.metrics-grid').classList.toggle('hidden');
    };

    el.innerHTML = `
        <div class="flex justify-between items-center">
            <div class="flex flex-col">
                <span class="text-xs font-black text-blue-400 italic tracking-tighter">${worker.language.toUpperCase()}</span>
                <span class="text-[9px] font-mono text-slate-500">${worker.id.substring(0,16)}</span>
            </div>
            <div class="status-dot w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_emerald]"></div>
        </div>
        
        <div class="metrics-grid hidden mt-4 pt-4 border-t border-slate-700 grid grid-cols-2 gap-4">
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">RAM Usage</p>
                <p class="ram-val text-sm font-black text-white">-- MB</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">CPU Threads</p>
                <p class="cpu-val text-sm font-black text-white">-- gor</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">Net I/O</p>
                <p class="net-val text-xs font-mono text-slate-300">--</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">Disk Activity</p>
                <p class="disk-val text-xs font-mono text-slate-300">--</p>
            </div>
        </div>
    `;
    list.appendChild(el);
}

export function updateHealthData(data) {
    const el = document.getElementById(`worker-${data.worker_id}`);
    if (!el) return;
    console.log('Updating health data for', data.worker_id, data);
    // Mise à jour des valeurs en temps réel 
    el.querySelector('.ram-val').innerText = `${data.ram} MB`;
    el.querySelector('.cpu-val').innerText = `${data.cpu} thr`;
    el.querySelector('.net-val').innerText = data.net_io;
    el.querySelector('.disk-val').innerText = data.disk_io;
}
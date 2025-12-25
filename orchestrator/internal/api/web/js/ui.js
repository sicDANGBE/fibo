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
    el.className = "p-4 bg-slate-700/50 rounded-xl border border-slate-600 mb-3";
    el.innerHTML = `
        <div class="flex items-center justify-between mb-2">
            <div class="flex flex-col">
                <span class="font-bold text-blue-400 text-xs">${worker.language.toUpperCase()}</span>
                <span class="text-[9px] font-mono text-slate-500">${worker.id.substring(0,8)}</span>
            </div>
            <div class="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_5px_emerald]"></div>
        </div>
        <div class="worker-stats border-t border-slate-600/50 pt-2">
            <span class="text-[10px] text-slate-500 italic">Waiting for metrics...</span>
        </div>
    `;
    list.appendChild(el);
}
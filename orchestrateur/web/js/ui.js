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

export function updateWorkerList(worker) {
    const list = document.getElementById('worker-list');
    const id = `worker-${worker.id}`;
    if (document.getElementById(id)) return;

    const el = document.createElement('div');
    el.id = id;
    el.className = "flex items-center justify-between p-4 bg-slate-700/50 rounded-xl border border-slate-600 animate-pulse";
    el.innerHTML = `
        <div class="flex flex-col">
            <span class="font-bold text-blue-400">${worker.language.toUpperCase()}</span>
            <span class="text-[10px] font-mono text-slate-400">${worker.id.substring(0,16)}</span>
        </div>
        <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
    `;
    list.appendChild(el);
    setTimeout(() => el.classList.remove('animate-pulse'), 2000);
}
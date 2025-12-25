const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync('sync.proto');
const syncProto = grpc.loadPackageDefinition(packageDefinition).sync;

function runFibo() {
    for (let run = 1; run <= 10; run++) {
        let a = 0n, b = 1n;
        const start = Date.now();
        for (let i = 0; i <= 400000; i++) {
            [a, b] = [b, a + b];
            if (i % 10000 === 0 && i > 0) {
                console.log(`[NODE] Run ${run} - ${i} iters - Temps: ${(Date.now() - start)/1000}s`);
            }
        }
    }
}

const client = new syncProto.Barrier('fibo-go:50051', grpc.credentials.createInsecure());
console.log("Node prêt, en attente du signal...");
client.waitToStart({}, (err) => {
    if (err) console.error(err);
    else runFibo();
});
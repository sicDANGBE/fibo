import grpc
import time
from sync_pb2 import Empty
from sync_pb2_grpc import BarrierStub

def run_fibo():
    for run in range(1, 11):
        a, b = 0, 1
        start_time = time.time()
        for i in range(400001):
            a, b = b, a + b
            if i % 10000 == 0 and i > 0:
                print(f"[PYTHON] Run {run} - {i} iters - Temps: {time.time() - start_time:.4f}s")

if __name__ == "__main__":
    with grpc.insecure_channel('fibo-go:50051') as channel:
        stub = BarrierStub(channel)
        print("Python prêt, en attente du signal...")
        stub.WaitToStart(Empty())
        run_fibo()
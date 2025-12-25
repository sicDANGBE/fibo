use futures_util::StreamExt;
use lapin::{options::*, types::FieldTable, BasicProperties, Connection, ConnectionProperties};
use num_bigint::BigInt;
use num_traits::{One, Zero};
use serde::{Deserialize, Serialize};
use std::time::Instant;

#[derive(Serialize, Deserialize, Debug)]
struct TaskMessage {
    language: String,
    serie: String,
    limit: u32,
}

#[derive(Serialize, Debug)]
struct ResultMessage {
    id: String,
    Série: String,
    num: u32,
    value: String,
}

async fn run_fibo_and_publish(chan: &lapin::Channel, task: TaskMessage) {
    let mut a: BigInt = BigInt::zero();
    let mut b: BigInt = BigInt::one();
    let worker_id = "rust-worker-01".to_string();

    println!("[RUST] Démarrage de la série : {}", task.serie);

    for i in 0..=task.limit {
        let temp = a.clone() + &b;
        a = b;
        b = temp;

        let res = ResultMessage {
            id: worker_id.clone(),
            Série: task.serie.clone(),
            num: i,
            value: a.to_string(), // Sérialisation BigInt en string
        };

        let payload = serde_json::to_vec(&res).unwrap();
        
        // Envoi immédiat à RabbitMQ
        chan.basic_publish(
            "",
            "fibo_results",
            BasicPublishOptions::default(),
            &payload,
            BasicProperties::default(),
        )
        .await
        .expect("Erreur lors de la publication");

        if i % 10000 == 0 {
            println!("[RUST] {} itérations envoyées...", i);
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("AMQP_ADDR").unwrap_or_else(|_| "amqp://guest:guest@rabbitmq:5672/%2f".into());
    let conn = Connection::connect(&addr, ConnectionProperties::default()).await?;
    let channel = conn.create_channel().await?;

    // Déclaration des queues
    channel.queue_declare("fibo_tasks", QueueDeclareOptions::default(), FieldTable::default()).await?;
    channel.queue_declare("fibo_results", QueueDeclareOptions::default(), FieldTable::default()).await?;

    println!("Rust est connecté à RabbitMQ. En attente de tâches...");

    let mut consumer = channel
        .basic_consume("fibo_tasks", "rust_consumer", BasicConsumeOptions::default(), FieldTable::default())
        .await?;

    while let Some(delivery) = consumer.next().await {
        let (_, delivery) = delivery.expect("error in consumer");
        let task: TaskMessage = serde_json::from_slice(&delivery.data)?;

        if task.language == "rust" || task.language == "all" {
            run_fibo_and_publish(&channel, task).await;
            channel.basic_ack(delivery.delivery_tag, BasicAckOptions::default()).await?;
        }
    }

    Ok(())
}
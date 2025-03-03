use actix_cors::Cors;
use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::sync::{Arc, Mutex};
use std::time::Duration;
use tokio::time::sleep;

#[derive(Debug, Serialize, Deserialize, Clone)]
struct StatEntry {
    id: u32,
    event_type: String,
    user_id: u32,
    timestamp: DateTime<Utc>,
    value: f64,
    metadata: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
struct StatSummary {
    total_users: u32,
    active_users: u32,
    new_today: u32,
    revenue: f64,
}

#[derive(Debug, Serialize, Deserialize)]
struct StatRequest {
    event_type: String,
    user_id: u32,
    value: f64,
    metadata: Option<String>,
}

struct AppState {
    stats: Mutex<Vec<StatEntry>>,
    next_id: Mutex<u32>,
}

async fn delay() {
    // Simüle edilmiş işleme gecikmesi
    sleep(Duration::from_millis(50)).await;
}

#[get("/stats")]
async fn get_all_stats(data: web::Data<Arc<AppState>>) -> impl Responder {
    delay().await;
    
    let stats = data.stats.lock().unwrap();
    HttpResponse::Ok().json(&*stats)
}

#[get("/stats/summary")]
async fn get_stats_summary(data: web::Data<Arc<AppState>>) -> impl Responder {
    delay().await;
    
    let stats = data.stats.lock().unwrap();
    
    // Gerçek uygulamada burada veri analizi yapılır
    // Burada sadece örnek değerler döndürüyoruz
    let summary = StatSummary {
        total_users: 8942,
        active_users: 2531,
        new_today: 147,
        revenue: 134952.89,
    };
    
    HttpResponse::Ok().json(summary)
}

#[post("/stats")]
async fn add_stat(
    data: web::Data<Arc<AppState>>, 
    stat_req: web::Json<StatRequest>
) -> impl Responder {
    delay().await;
    
    let mut id = data.next_id.lock().unwrap();
    let new_id = *id;
    *id += 1;
    
    let new_stat = StatEntry {
        id: new_id,
        event_type: stat_req.event_type.clone(),
        user_id: stat_req.user_id,
        timestamp: Utc::now(),
        value: stat_req.value,
        metadata: stat_req.metadata.clone(),
    };
    
    let mut stats = data.stats.lock().unwrap();
    stats.push(new_stat.clone());
    
    HttpResponse::Created().json(new_stat)
}

#[get("/stats/revenue")]
async fn get_revenue_stats(data: web::Data<Arc<AppState>>) -> impl Responder {
    delay().await;
    
    let stats = data.stats.lock().unwrap();
    
    // Burada gelir verilerini filtreliyoruz
    let revenue_stats: Vec<&StatEntry> = stats
        .iter()
        .filter(|s| s.event_type == "revenue" || s.event_type == "purchase")
        .collect();
    
    HttpResponse::Ok().json(revenue_stats)
}

#[get("/stats/user/{user_id}")]
async fn get_user_stats(
    data: web::Data<Arc<AppState>>,
    path: web::Path<u32>
) -> impl Responder {
    delay().await;
    
    let user_id = path.into_inner();
    let stats = data.stats.lock().unwrap();
    
    let user_stats: Vec<&StatEntry> = stats
        .iter()
        .filter(|s| s.user_id == user_id)
        .collect();
    
    if user_stats.is_empty() {
        return HttpResponse::NotFound().body(format!("Kullanıcı için istatistik bulunamadı (ID: {})", user_id));
    }
    
    HttpResponse::Ok().json(user_stats)
}

#[get("/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().body("Sağlık durumu: Çalışıyor")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Mock veri oluşturma
    let mut initial_stats = Vec::new();
    let mock_events = vec!["login", "pageview", "purchase", "revenue", "signup"];
    
    for i in 1..20 {
        let event_type = mock_events[i % mock_events.len()];
        let user_id = (i % 10) + 1;
        let value = match event_type {
            "revenue" | "purchase" => (i as f64) * 23.45,
            _ => i as f64,
        };
        
        initial_stats.push(StatEntry {
            id: i,
            event_type: event_type.to_string(),
            user_id,
            timestamp: Utc::now(),
            value,
            metadata: None,
        });
    }
    
    // Uygulama durumu oluşturma
    let app_state = Arc::new(AppState {
        stats: Mutex::new(initial_stats),
        next_id: Mutex::new(20),
    });
    
    println!("Rust İstatistik Servisi 8081 portunda başlatılıyor");
    
    // HTTP sunucusunu başlatma
    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header()
            .max_age(3600);
            
        App::new()
            .wrap(cors)
            .app_data(web::Data::new(app_state.clone()))
            .service(get_all_stats)
            .service(get_stats_summary)
            .service(add_stat)
            .service(get_revenue_stats)
            .service(get_user_stats)
            .service(health_check)
    })
    .bind("127.0.0.1:8081")?
    .run()
    .await
}
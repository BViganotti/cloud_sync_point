use std::{
    collections::HashMap,
    sync::{Arc, Mutex},
    time::Duration,
};
use sync::oneshot;
use tokio::*;
use uuid::Uuid;
use warp::Filter;

// fwd declaring our thread safe state map
type SharedState = Arc<Mutex<HashMap<Uuid, oneshot::Sender<()>>>>;

#[tokio::main]
async fn main() {
    // setup our path endpoint requiring a uuid
    let endpoint_path = warp::path("wait-for-second-party").and(warp::path::param::<Uuid>());

    // this will keep track of the requests
    let keep_track_of_state = Arc::new(Mutex::new(HashMap::new()));

    // this filter will match any route, we need this because the unique id will be unique
    let state_filter = warp::any().map(move || Arc::clone(&keep_track_of_state));

    let post_endpoint = endpoint_path
        .and(warp::post())
        .and(state_filter.clone())
        .and_then(request_handler);

    let get_endpoint = endpoint_path
        .and(warp::get())
        .and(state_filter.clone())
        .and_then(request_handler);

    // i want this endpoint to be able to handle both post and get requests
    let endpoints = post_endpoint.or(get_endpoint);

    println!("cloud sync-point running at: 127.0.0.1:3030");
    warp::serve(endpoints).run(([127, 0, 0, 1], 3030)).await;
}

async fn request_handler(
    unique_id: Uuid,
    req_state: SharedState,
) -> Result<impl warp::Reply, warp::Rejection> {
    let (tx, rx) = oneshot::channel();

    println!("request id: {}", unique_id);

    {
        // make state map data safe to access
        let mut req_state_guard = req_state.lock().unwrap();
        // we try to remove the unique_id key we got from the map, if we succeed (meaning sender won't be None),
        // that mean the second party requested our uri
        if let Some(sender) = req_state_guard.remove(&unique_id) {
            // we have our second party, let's unblock and reply OK
            sender.send(()).ok();
            //println!("The second party requested the URI");
            return Ok(warp::reply::with_status(
                "the second party requested the URI",
                warp::http::StatusCode::OK,
            ));
        } else {
            // the unique_id wasn't in the map, that means it is just the first party requesting the uri
            // this makes us able to send a signal to the second party when it will request the uri
            println!("the first party requested the URI");
            req_state_guard.insert(unique_id, tx);
        }
    }

    tokio::select! {
        _ = rx => {
            println!("the second party requested the URI\n");
            Ok(warp::reply::with_status("the second party requested the URI", warp::http::StatusCode::OK))},
        // if this complete, it means that the second party didn't request the uri in time
        _ = tokio::time::sleep(Duration::from_secs(10)) => {
            // locking the map before removing the state
            println!("the second party did not request the URI, reponding with TIMEOUT\n");
            let mut req_state_guard = req_state.lock().unwrap();
            req_state_guard.remove(&unique_id);
            Ok(warp::reply::with_status("TIMEOUT", warp::http::StatusCode::REQUEST_TIMEOUT))
        },
    }
}

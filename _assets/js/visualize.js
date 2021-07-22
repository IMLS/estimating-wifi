
function setup () {
    console.log("setup");
    // Boxes
    createCanvas(screen.width, 600);
    background("#ffffff");
    colorMode(HSB, 100);
    noStroke();
}

class Block {
    constructor (mfg, w, yndx) {
        this.width = w;
        this.height = height / gql.unique.manufacturers.size;
        this.x = width;
        this.y = yndx * this.height;
        this.mfg = mfg;
    }

    move () {
        this.x = this.x - 1;
    }

    draw () {
        fill(Math.floor(360 / this.mfg), 80, 80);
        rect(this.x, this.y, this.width, this.height)
    }
}

class Pixel {
    constructor (ndx) {
        this.x = ndx;
        this.y = 0;
        this.color = color((ndx * 3) % 100, 80, 80);
    }

    move () {
        this.y = this.y + 1;
    }

    draw () {
        this.move()
        fill(this.color);
        ellipse(this.x, this.y, 3, 3);
    }
}

let ticker = 0;
let slow_tick = 0;
let prev_slow_tick = 0;

let blocks = [];
let block_ndx = 0;
let block_width = 10;

let pixel_setup = false;
let active_points = new Map();
let events = [];

function draw () {
    if (gql && gql.STATUS) {
        if (!pixel_setup) {
            pixel_setup = true;
            // Iterate through the event ids.
            for (let e of gql.events) {
                events.push(parseInt(e));
            }
            events.sort((a, b) => { return a - b; });
            console.log("events", events);
            // createCanvas(gql.unique.devices.size, 360);
            // resizeCanvas(, windowHeight);
            clear();
        } else {
            if ((ticker % 5) == 0) {slow_tick += 1;}
            if ((slow_tick < events.length) && (prev_slow_tick != slow_tick)) {
                clear();
                // console.log("slow_tick", slow_tick, "events_length", events.length);
                prev_slow_tick = slow_tick;
                current_event = events[slow_tick];
                // console.log("ce", current_event);
        
                // This is slow.
                for (let o of gql.data) {
                    if (o.event_id == current_event) {
                        // console.log("created", o.event_id, o.id)
                        active_points.set(o.id, new Pixel(o.patron_index));
                    }
                }
                
                // console.log("active_points count", active_points.size)
                for (let k of active_points.keys()) {
                    active_points.get(k).draw();
                }
        
                for (let k of active_points.keys()) {
                    if (active_points.get(k).y > height) {
                        active_points.delete(k);
                    }
                }
            } // end if over slow ticks
            ticker += 1;    
        }
    }
}


function draw2 () {
    clear()
    if (gql && gql.STATUS) {

        if ((ticker % block_width) == 0) {
            for (o of gql.data) {
                if (o.event_id == block_ndx) {
                    console.log("event", o.event_id);
                    blocks.push(new Block(o.manufacturer_index, block_width, o.manufacturer_index));
                }
            }
            block_ndx += 1;
        }

        for (let b of blocks) {
            b.move();
            b.draw();
        }
    
        ticker += 1;
    }
}
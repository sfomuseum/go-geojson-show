window.addEventListener("load", function load(event){

    // Null Island
    var map = L.map('map').setView([0.0, 0.0], 12);

    var select = function(show_id){

	unselect();
	
	var el = document.getElementById(show_id);
	
	if (el){
	    el.setAttribute("class", "selected");
	    el.scrollIntoView();
	}
	
    };
    
    var unselect = function(){
	
	var current = document.querySelector(".selected");
	
	if (current){
	    current.classList.remove("selected");
	}
    };

    map.on("click", function(e){
	unselect();
    });
    
    var init = function(cfg) {
	
	fetch("/features.geojson")
	    .then((rsp) => rsp.json())
	    .then((f) => {

		var features = f.features;
		var count = features.length;
		
		for (var i=0; i < count; i++){
		    var show_id = "show-" + (i+1);
		    f.features[i]["properties"]["show:id"] = show_id;
		}
		
		var raw_el = document.querySelector("#raw");
		
		var format = function(show_id, str){
		    
		    // Remember: wof_format is defined by the /wasm/wof_format.wasm binary.
			// Details below.
			
			wof_format(str).then((rsp) => {
			    append(show_id, rsp);
			}).catch((err) => {
			    console.warn("Unable to format feature", err, str);
			    append(show_id, str);
			});
		};
		
		var append = function(show_id, str) {
		    var pre = document.createElement("pre");
		    pre.setAttribute("id", show_id);
		    pre.appendChild(document.createTextNode(str));		    
		    raw_el.appendChild(pre);
		};
		
		if (raw_el){
		    
		    // Remember: Both sfomuseum.wasm.fetch and the WASM binary are imported and registered
		    // in show.go. For details see: https://github.com/whosonfirst/go-whosonfirst-format-wasm
		    
		    sfomuseum.wasm.fetch("/wasm/wof_format.wasm").then(rsp => {
			
			var features = f.features;
			var count = features.length;
			
			for (var i=0; i < count; i++){
			    
			    var show_id = features[i]["properties"]["show:id"];
			    var this_f = structuredClone(features[i]);
			    
			    delete(this_f["properties"]["show:id"]);
			    var str_f = JSON.stringify(this_f);
			    
			    format(show_id, str_f);
			}
			
		    }).catch((err) => {
			console.warn("Unable to load wof_format.wasm", err);
			var str_f = JSON.stringify(f, "", " ");		    
			append(0, str_f);
		    });
		    
		}

		var geojson_args = {
		    onEachFeature: function (feature, layer) {

			layer.on("click", function(e){			    
			    var show_id = feature["properties"]["show:id"];
			    select(show_id);
			});

			var label_props = cfg.label_properties;

			if (label_props){
			    var count_props = label_props.length;
			    
			    if (count_props > 0) {
				
				var label_text = [];
				
				for (var i=0; i < count_props; i++){
				    
				    var prop = label_props[i];
				    var value = feature.properties[ prop ];
				    
				    label_text.push("<strong>" + prop + "</strong> " + value);
				}
				
				if (label_text.length > 0){ 
				    layer.bindPopup(label_text.join("<br />"));
				}
			    }
			    
			}
		    }
		};

		if (cfg.style){
		    geojson_args.style = cfg.style;
		}

		if (cfg.point_style) {

		    geojson_args.pointToLayer = function (feature, latlng) {
			return L.circleMarker(latlng, cfg.point_style);
		    }
		    
		}

		var geojson_layer = L.geoJSON(f, geojson_args);
		geojson_layer.addTo(map);
		
		var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);
		
		var sw = bounds[0];
		var ne = bounds[1];
		
		if ((sw[0] == ne[0]) && (sw[1] == ne[1])){
		    map.setView(sw, 12);
		} else {
		    map.fitBounds(bounds);
		}
		
	    }).catch((err) => {
		console.error("Failed to render features", err);
	    });
    };

    fetch("/map.json")
	.then((rsp) => rsp.json())
	.then((cfg) => {

	    switch (cfg.provider) {
		case "leaflet":

		    var tile_url = cfg.tile_url;

		    var tile_layer = L.tileLayer(tile_url, {
			maxZoom: 19,
		    });
		    
		    tile_layer.addTo(map);
		    break;
		    
		case "protomaps":		    

		    var tile_url = cfg.tile_url;

		    var tile_layer = protomapsL.leafletLayer({
			url: tile_url,
			theme: cfg.protomaps.theme,
		    })

		    tile_layer.addTo(map);
		    break;
		    
		default:
		    console.error("Uknown or unsupported map provider");
		    return;
	    }
	    
	    init(cfg);
	    
	}).catch((err) => {
	    console.error("Failed to retrieve features", err);
	});
    
});

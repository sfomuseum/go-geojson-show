window.addEventListener("load", function load(event){

    // Null Island
    const map = L.map('map').setView([0.0, 0.0], 12);

    const applyCustomStyles = function(feature, style){

	if (! "custom" in style){
	    return style;
	}
	
	if ("color_map" in style.custom){
	    
	    const color_map = style.custom.color_map;
	    const prop = color_map.property;
	    
	    if (prop in feature.properties){
		
		const v = feature.properties[prop];
		const str_v = String(v);
		
		if (str_v in color_map.key){
		    
		    if ("color" in color_map.key[str_v]){
			style.color = color_map.key[str_v]["color"];
		    }
		    
		    if ("opacity" in color_map.key[v]){
			style.opacity = color_map.key[v]["opacity"];
		    }				    
		}
	    }
	}
	
	if ("fill_map" in style.custom){
	    
	    const fill_map = style.custom.fill_map;
	    const prop = fill_map.property;
	    
	    if (feature.properties[prop]){
		
		const v = feature.properties[prop];
		const str_v = String(v);
		
		if (str_v in fill_map.key){

		    if ("color" in fill_map.key[v]){ 
			style.fillColor = fill_map.key[v]["color"];
		    }
		    
		    if ("opacity" in fill_map.key[v]){ 
			style.fillOpacity = fill_map.key[v]["opacity"];
		    }
		    
		}
	    }
	    
	}
	
	if ("pane_map" in style.custom){

	    const pane_map = style.custom.pane_map;
	    const prop = pane_map.property;

	    if (prop in feature.properties){

		const v = feature.properties[prop];
		const str_v = String(v);

		if (str_v in pane_map.key){
		    
		    const label = pane_map.key[str_v];
		    
		    if (map.getPane(label)){
			// console.log("YES Assign feature to pane", label);			
			style.pane = label;
		    }
		    
		} else if ("*" in pane_map.key){

		    const label = pane_map.key["*"];
		    
		    if (map.getPane(label)){
			// console.log("YES Assign feature to pane", label);			
			style.pane = label;
		    }
		} else {
		    // pass
		}
	    }
	}
	    
	return style;
    };
    
    const select = function(show_id){

	unselect();
	
	var el = document.getElementById(show_id);
	
	if (el){
	    el.setAttribute("class", "selected");
	    el.scrollIntoView();
	}
	
    };
    
    const unselect = function(){
	
	var current = document.querySelector(".selected");
	
	if (current){
	    current.classList.remove("selected");
	}
    };

    map.on("click", function(e){
	unselect();
    });
    
    const init = function(local_cfg, map_cfg) {

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
		    
		    sfomuseum.golang.wasm.fetch("/wasm/wof_format.wasm").then(rsp => {
			
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

			if (map_cfg.leaflet) {
			    
			    var label_props = map_cfg.leaflet.label_properties;
			    
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
		    }
		};

		if ((map_cfg.leaflet) && (map_cfg.leaflet.style)){
		    // This doesn't work because we don't know what feature is...
		    // const style = applyCustomStyles(feature, map_cfg.leaflet.style);
		    const style = map_cfg.leaflet.style;
		    geojson_args.style = style;
		}

		if ((map_cfg.leaflet) && (map_cfg.leaflet.point_style)){

		    geojson_args.pointToLayer = function (feature, latlng) {
			const style = applyCustomStyles(feature, map_cfg.leaflet.point_style);			
			return L.circleMarker(latlng, style);
		    }
		    
		}

		var geojson_layer = L.geoJSON(f, geojson_args);

		if (local_cfg.cluster_markers){
		    const markers = L.markerClusterGroup();
		    markers.addLayer(geojson_layer);
		    markers.addTo(map);
		} else {
		    geojson_layer.addTo(map);
		}
		
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

    fetch("/config.json").then(rsp =>
	rsp.json()
    ).then((local_cfg) => {
	
	fetch("/map.json")
	    .then((rsp) => rsp.json())
	    .then((map_cfg) => {
		
		switch (map_cfg.provider) {
		    case "leaflet":
			
			var tile_url = map_cfg.tile_url;
			
			var tile_layer = L.tileLayer(tile_url, {
			    maxZoom: 19,
			});
			
			tile_layer.addTo(map);
			break;

		    case "esri":

			if ("esri_feature_layers" in local_cfg){
			    
			    const layers = local_cfg.esri_feature_layers;
			    const count_layers = layers.length;
			    
			    for (var i=0; i < count_layers; i++){
				
				var tile_uri = layers[i];
				
				var tile_style = {
				    color: '#000',
				    weight: 1,
				    opacity: 1,
				    fillColor: '#fff',
				    fillOpacity: 0,
				};
				
				const tile_u = new URL(tile_uri);
				const tile_q = tile_u.searchParams;
				
				for (k in tile_style){
				    
				    var q_key = "_" + k;
				    
				    if (tile_q.has(q_key)){
					tile_style[k] = tile_q.get(q_key);
					tile_q.delete(q_key);
				    }
				}
				
				tile_u.searchParams = tile_q;
				tile_uri = tile_u.toString();
				
				var tile_args = {
				    url: tile_uri,
				    style: tile_style,
				    /*
				       pointToLayer: function(feature, latlng) {
				       return L.circleMarker(latlng, {
				       radius: 8,
				       fillColor: "#ff0000",
				       color: "#fff",
				       weight: 1,
				       opacity: 1,
				       fillOpacity: 0.1
				       });
				       }
				     */
				};
				
				var tile_layer = L.esri.featureLayer(tile_args);
				tile_layer.addTo(map);
			    }
			}
			
			break;
			
		    case "protomaps":		    
			
			var tile_url = map_cfg.tile_url;		

			var pm_args = {
			    url: tile_url,
			    theme: map_cfg.protomaps.theme,
			    flavor: map_cfg.protomaps.theme,
			};
			
			if ("max_data_zoom" in map_cfg.protomaps){
			    pm_args.maxDataZoom = map_cfg.protomaps.max_data_zoom;
			}

			var tile_layer = protomapsL.leafletLayer(pm_args);
			
			tile_layer.addTo(map);
			break;
			
		    default:
			console.error("Uknown or unsupported map provider");
			return;
		}
		
		if (("leaflet" in map_cfg) && ("panes" in map_cfg.leaflet)){
		    
		    for (label in map_cfg.leaflet.panes){
			const p = map.createPane(label);
			p.style.zIndex = map_cfg.leaflet.panes.label;
			console.debug("Created pane", label, map_cfg.leaflet.panes.label);
		    }
		}

		init(local_cfg, map_cfg);
		
	    }).catch((err) => {
		console.error("Failed to retrieve map config", err);
	    });
	
    }).catch((err) => {
	console.error("Failed to retrieve local cfg", err);
    });
    
});

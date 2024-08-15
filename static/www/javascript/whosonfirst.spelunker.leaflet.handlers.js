var whosonfirst = whosonfirst || {};
whosonfirst.spelunker = whosonfirst.spelunker || {};
whosonfirst.spelunker.leaflet = whosonfirst.spelunker.leaflet || {};

whosonfirst.spelunker.leaflet.handlers = (function(){

	var self = {

	    'point': function(layer_args){
		
		return function(feature, latlon){
		    
		    var m = L.circleMarker(latlon, layer_args);
		    
		    try {
			var props = feature['properties'];
			var label = props['lflt:label_text'];
			var href = props['lflt:label_href'];
			
			if ((! label) && (props['lflt:label_names'])){
			    var str_coords = JSON.stringify([ latlon.lng, latlon.lat ]);
			    label = props['lflt:label_names'][str_coords];
			}

			if ((! href) && (props['lflt:label_links'])){
			    var str_coords = JSON.stringify([ latlon.lng, latlon.lat ]);
			    href = props['lflt:label_links'][str_coords];
			}
			
			if (label){
			    
			    var label_args = {
				noHide: false,
				interactive: true,
			    }

			    if (layer_args.tooltips_pane){
				label_args.pane = layer_args.tooltips_pane;
			    }
			    
			    var t = m.bindTooltip(label, label_args);

			    if (href){
				
				t.on("click", function(){
				    location.href = href;
				    return false;
				});
			    }
			}
		    }
		    
		    catch (e){
			console.log("failed to bind label because " + e);
		    }
		    
		    return m;
		};
	    },
	};
    
	return self;
})();

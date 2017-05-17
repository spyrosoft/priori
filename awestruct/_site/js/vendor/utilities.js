// Source: https://github.com/spyrosoft/javascript-utilities

var Utilities = {
	keyCodeLookupTable : {
		13 : 'enter',
		27 : 'escape',
		16 : 'shift',
		17 : 'control',
		18 : 'alt',
		91 : 'meta',
		20 : 'capslock',
		8 : 'backspace',
		9 : 'tab',
		189 : '-',
		187 : '=',
		219 : '[',
		221 : ']',
		220 : '\\',
		186 : ';',
		222 : '\'',
		188 : ',',
		190 : '.',
		191 : '/'
	},
	
	keyCodeShiftLookupTable : {
		// The ASCII is out of order for the number row symbols so we can't do a translation like we did for numbers and letters
		49 : '!',
		50 : '@',
		51 : '#',
		52 : '$',
		53 : '%',
		54 : '^',
		55 : '&',
		56 : '*',
		57 : '(',
		58 : ')',
		189 : '_',
		187 : '+',
		219 : '{',
		221 : '}',
		220 : '|',
		186 : ':',
		222 : '"',
		188 : '<',
		190 : '>',
		191 : '?'
	},
	
	keyCodeAlphaMin : 65,
	keyCodeAlphaMax : 90,
	
	keyCodeLookup : function( keyEvent ) {
		if ( keyEvent.shiftKey ) { return this.keyCodeShiftLookup( keyEvent ); }
		
		var keyCode = keyEvent.keyCode;
		
		var numericMin = 48;
		var numericMax = 57;
		if ( keyCode <= numericMax && keyCode >= numericMin ) {
			return String.fromCharCode( keyCode );
		}
		
		var alphaOffset = 32;
		if ( keyCode <= this.keyCodeAlphaMax && keyCode >= this.keyCodeAlphaMin ) {
			return String.fromCharCode( keyCode + alphaOffset );
		}
		
		if ( typeof this.keyCodeLookupTable[ keyCode ] === 'undefined' ) {
			return null;
		} else {
			return this.keyCodeLookupTable[ keyCode ];
		}
		
		return null;
	},
	
	keyCodeShiftLookup( keyEvent ) {
		var keyCode = keyEvent.keyCode;
		
		if ( keyCode <= this.keyCodeAlphaMax && keyCode >= this.keyCodeAlphaMin ) {
			return String.fromCharCode( keyCode );
		}
		
		if ( typeof this.keyCodeShiftLookupTable[ keyCode ] === 'undefined' ) {
			if ( typeof this.keyCodeLookupTable[ keyCode ] === 'undefined' ) {
				return null;
			} else {
				return this.keyCodeLookupTable[ keyCode ];
			}
		} else {
			return this.keyCodeShiftLookupTable[ keyCode ];
		}
		
		return null;
	},
	
	isEnter : function( keyEvent ) {
		var keyCode = keyEvent.keyCode;
		if ( typeof this.keyCodeLookupTable[ keyCode ] === 'undefined' ) { return false; }
		return this.keyCodeLookupTable[ keyCode ] === 'enter';
	},
	
	//TODO: Is this useful for debugging, or does it obfuscate where the error originated?
	customError : function( errorName, errorMessage ) {
		throw {
			'name' : errorName,
			'message' : errorMessage
		};
	},
	
	sixHexToRgb : function( sixHex ) {
		var rgb = new Array();
		rgb[ 0 ] = this.hexToDecimal( sixHex.substr( 0, 2 ) );
		rgb[ 1 ] = this.hexToDecimal( sixHex.substr( 2, 2 ) );
		rgb[ 2 ] = this.hexToDecimal( sixHex.substr( 4, 2 ) );
		return rgb;
	},
	
	rgbToSixHex : function( rgb ) {
		var rHex = this.padStringFront( this.decimalToHex( rgb[ 0 ] ), '0', 2 );
		var gHex = this.padStringFront( this.decimalToHex( rgb[ 1 ] ), '0', 2 );
		var bHex = this.padStringFront( this.decimalToHex( rgb[ 2 ] ), '0', 2 );
		var sixHex = rHex + gHex + bHex;
		return sixHex;
	},
	
	decimalToHex : function( decimal ) {
		if ( typeof decimal !== 'number' ) { return undefined; }
		return decimal.toString( 16 );
	},
	
	hexToDecimal : function( hex ) {
		return parseInt( hex, 16 );
	},
	
	padStringFront : function( string, padWith, padLength ) {
		if ( typeof string !== 'string' ) { return undefined; }
		var paddedString = string;
		while ( paddedString.length < padLength ) {
			paddedString = padWith + paddedString;
		}
		return paddedString;
	},
	
	strictInteger : function( input ) {
		if ( typeof input === 'number' ) {
			if ( parseFloat( input ) === parseInt( input ) ) { return input; }
			return undefined;
		}
		if ( typeof input === 'string' && input === parseInt( input ).toString() ) {
			return parseInt( input );
		}
		return undefined;
	},
	
	strictFloat : function( input ) {
		if ( isFinite( input ) ) { return parseFloat( input ); }
		return undefined;
	},
	
	objectEqual : function( object1, object2 ) {
		if ( typeof object1 !== 'object' || typeof object2 !== 'object' ) { return false; }
		
		var objectsAreEqual = true;
		
		for ( var key in object1 ) {
			var value = object1[ key ];
			if ( object2[ key ] === undefined ) {
				objectsAreEqual = false;
				break;
			}
			if ( typeof object2[ key ] === 'object' ) {
				objectsAreEqual = this.objectEqual( object1[ key ], object2[ key ] );
				if ( ! objectsAreEqual ) { break; }
			} else if ( object1[ key ] !== object2[ key ] ) {
				objectsAreEqual = false;
				break;
			}
		}
		
		return objectsAreEqual;
	},
	
	arrayEqual : function( array1, array2 ) {
		return this.objectEqual( array1, array2 );
	},
	
	formatHumanReadableDollars : function( number ) {
		number = this.strictFloat( number );
		if ( number === undefined ) { throw 'The argument provided was not a number: ' + number; }
		number = number.toFixed( 2 ).split( '.' );
		return '$'
			+ number[ 0 ].replace( /(\d)(?=(\d{3})+$)/g, '$1,' )
			+ '.'
			+ number[ 1 ];
	},
	
	convertDollarsToCents : function( dollars ) {
		dollars = this.strictFloat( dollars );
		if ( dollars !== undefined ) { throw 'The argument provided was not a dollar amount: ' + dollars; }
		return Math.ceil( parseFloat( dollars ) * 100 );
	},
	
	parseGetParametersFromURL : function() {
		var full_url = window.location.href;
		if ( full_url.indexOf( '?' ) === -1 ) { return {}; }
		var get_parameter_string = full_url.substring( full_url.indexOf( '?' ) + 1, full_url.length );
		return this.parseGetParameters( get_parameter_string );
	},
	
	parseGetParameters : function( get_parameter_string ) {
		var get_parameters = {};
		var get_parameters_and_values = get_parameter_string.split( '&' );
		for ( var i in get_parameters_and_values ) {
			var get_parameter_and_value = get_parameters_and_values[ i ].split( '=' );
			var get_parameter = decodeURIComponent( get_parameter_and_value[ 0 ] );
			var get_value = decodeURIComponent( get_parameter_and_value[ 1 ] );
			get_parameters[ get_parameter ] = get_value;
		}
		return get_parameters;
	},
	
	parseJSON : function( json_to_parse ) {
		try {
			var parsed = JSON.parse( json_to_parse );
			return parsed;
		} catch(e) {}
		return undefined;
	}
};

if ( typeof console === 'undefined' ) { var console = { log : function() {} }; }

if ( typeof Belt === 'undefined' ) { var Belt = Utilities; }

var live_or_dev = 'live';
if (window.location.hostname === 'localhost') {
	live_or_dev = 'dev';
}

var error_response_types_and_callbacks = {
	'errors' : display_error,
	'messages' : display_message,
	'debug' : debug,
	'fields' : highlight_field
};

var success_response_callbacks = {};

$('form.ajax').submit(ajax_form_submission);


function ajax_form_submission(submit_event) {
	submit_event.preventDefault();
	var post_url = $(this).attr('action');
	var post_data = $(this).serialize();
	var callback_name = $(this).attr('data-callback');
	ajax_submission(post_url, post_data, callback_name);
}

function api_call(post_data, callback_name) {
	ajax_submission(
		'/api/',
		post_data,
		callback_name
	);
}

function ajax_submission(post_url, post_data, callback_name) {
	var ajax_request_response = function(response_data) {
		ajax_response_success_or_errors(
			post_data,
			response_data,
			callback_name
		);
	};
	$.post(post_url, post_data)
		.done(ajax_request_response)
		.fail(ajax_request_failure);
}

function ajax_response_success_or_errors(post_data, response_data, callback_name) {
	try {
		var response = JSON.parse(response_data);
		if (response.success === false) {
			ajax_request_response_errors(response);
		} else {
			if (callback_name !== undefined
				&& success_response_callbacks[callback_name] !== undefined) {
				success_response_callbacks[callback_name](post_data, response);
			} else {
				display_message('Success');
			}
		}
		//TODO: Handle unexpected JSON responses
	} catch (error) {
		notify_admin('ajax_request_response: ' + error + '. Data: ' + response_data.toString());
	}
}

function ajax_request_response_errors(response) {
	for (var response_type in error_response_types_and_callbacks) {
		if (typeof response[response_type] !== 'undefined') {
			ajax_response_actions(response[response_type], error_response_types_and_callbacks[response_type]);
		}
	}
}

function ajax_response_actions(actions, callback) {
	if (actions === null || actions.length === undefined || actions.length === 0) {return;}
	for (var i = 0; i < actions.length; i++) {
		callback(actions[i]);
	}
	if (callback === highlight_field) {
		$('[name=' + actions[0] + ']').focus();
	}
}

function highlight_field(field) {
	$('[name=' + field + ']')
		.addClass('issue')
		.on('focus', function() {
			$(this).removeClass('issue');
		});
}

function ajax_request_failure() {
	display_error('Connection with the server failed. Please check your internet connection. Otherwise, something is wrong on our end - please try again later.');
}

function debug(message) {
	if (live_or_dev === 'dev') {
		display_message('DEBUG: ' + message);
	} else if (typeof console !== 'undefined' && typeof console.log === 'function') {
		console.log(message);
	}
}

function display_error(message) {
	append_to_messages('error', message);
}

function display_message(message) {
	append_to_messages('message', message);
}

function append_to_messages(template_type, message) {
	var error_template = $('.' + template_type + '.template')[0];
	if (error_template === undefined) {
		alert(message);
		return;
	}
	var add_message = function(clone) {
		clone.html(message);
		var remove_myself = function() {
			$(this).remove();
		};
		clone.on('click', remove_myself);
	};
	use_template(error_template, '.messages', add_message);
}

function use_template(template, parent_selector, functionality) {
	var clone = $(template).clone().removeClass('template');
	if (typeof functionality === 'function') {
		functionality(clone);
	}
	$(parent_selector).append(clone);
}

function notify_admin(message) {
	if (live_or_dev === 'dev') {
		debug(message);
	} else {
		display_error('Something went wrong. An admin has been notified and will resolve this problem as soon as possible.');
	}
	var notify_admin_fail = function() {
		display_error('The attempt to contact an admin has failed. It is possible your internet connection has been interrupted. Otherwise if the issue persists, please contact us directly to resolve the issue.');
	};
	var server_message = {
		'function' : 'notify-admin',
		'message' : message,
		'token' : 'PpPub4GjM4'
	};
	$.post('/api/', server_message)
		.fail(notify_admin_fail);
}

$(document).foundation();

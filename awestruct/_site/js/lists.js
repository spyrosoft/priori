success_response_callbacks['reset-lists'] = populate_lists;

function reset_lists() {
	$('.list-items').html('');
	ajax_submission(
		'/api/',
		{'action': 'lists'},
		populate_lists
	);
}

function populate_lists(list_data) {
	console.log(list_data)
	//populate_lists();
}

function add_list(list_name) {
	var list_template = $('.list-item.template')[0];
	$(list_template).find('a').html(list_name);
	use_template(list_template, '.list-items');
}

function new_list_success(form_json, response_json) {
	add_list(response_json['name']);
	$('.new-list-name').first().val('');
}

success_response_callbacks['new-list-success'] = new_list_success;

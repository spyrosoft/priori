function reset_lists() {
	$('.list-items').html('');
	api_call(
		{'action': 'lists'},
		'populate-lists'
	);
}

success_response_callbacks['populate-lists'] = populate_lists_callback;

function populate_lists_callback(post_data, lists_data) {
	for (var i = 0; i < lists_data.length; i++) {
		add_list(lists_data[i]['name'], lists_data[i]['id']);
	}
}

function add_list(name, id) {
	var list_template = $('.list-item.template')[0];
	$(list_template).find('a')
		.html(name)
		.attr('href', '/list/#' + id);
	use_template(list_template, '.list-items');
}



success_response_callbacks['new-list-success'] = new_list_success;

function new_list_success(post_data, new_list_data) {
	add_list(new_list_data['name'], new_list_data['id']);
	$('.new-list-name').first().val('');
}



reset_lists();